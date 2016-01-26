package server

import (
	"sync"
	"sync/atomic"
)

// TODO: Add an 'alive' property so that a 'close' call stops the listeners.

type BytePipe struct {
	buf  []byte     // Internal buffer
	size int32      // capacity of the buffer
	r    int32      // Read index
	rc   *sync.Cond // ReadCondition -- used when reader needs to
	w    int32      // Write index
	wc   *sync.Cond // WriteCondition -- used when writer has no more room
}

func NewBytePipe(size int32) *BytePipe {
	if size == 0 {
		size = 32368
	}
	return &BytePipe{
		buf:  make([]byte, size), // 32kb of room for buffering.
		r:    0,
		w:    0,
		rc:   sync.NewCond(&sync.Mutex{}),
		wc:   sync.NewCond(&sync.Mutex{}),
		size: size,
	}
}

func (bp *BytePipe) Len() int {
	r := atomic.LoadInt32(&bp.r)
	w := atomic.LoadInt32(&bp.w)
	if r <= w {
		return int(w - r)
	}

	return int((bp.size - r) + (r - w))
}

func (bp *BytePipe) Write(b []byte) {
	r := atomic.LoadInt32(&bp.r)
	w := atomic.LoadInt32(&bp.w)
	l := int32(len(b))

	if (w >= r && w+l <= bp.size) || (r > w && r-w > l) {
		// more than enough space to write entire thing into it.
		copy(bp.buf[w:], b)
		atomic.StoreInt32(&bp.w, w+l)
		bp.rc.L.Lock()
		bp.rc.Signal()
		bp.rc.L.Unlock()
		return
	}
	if (w == bp.size && r == 0) || w == r-1 {
		// no space eft until reader reads something
		bp.wc.L.Lock()
		if r == atomic.LoadInt32(&bp.r) {
			bp.wc.Wait()
		}
		bp.wc.L.Unlock()
		bp.Write(b)
		return
	}
	if w < r {
		// not enough space! write what we can now, and then call write again to wait.
		bidx := r - w - 1 // cant overwrite that last byte!
		copy(bp.buf[w:], b[:bidx])
		atomic.StoreInt32(&bp.w, r-1)
		bp.rc.L.Lock()
		bp.rc.Signal()
		bp.rc.L.Unlock()
		bp.Write(b[bidx:]) // Now write the rest
		return
	}
	// this means we dont have enought space at the end, write and try again.
	bidx := bp.size - w
	copy(bp.buf[w:], b[:bidx])
	if r == 0 {
		atomic.StoreInt32(&bp.w, bp.size)
	} else {
		atomic.StoreInt32(&bp.w, 0)
	}
	bp.rc.L.Lock()
	bp.rc.Signal()
	bp.rc.L.Unlock()
	bp.Write(b[bidx:]) // Now write the rest
}

func (bp *BytePipe) Read(b []byte) int {
	r := atomic.LoadInt32(&bp.r)
	w := atomic.LoadInt32(&bp.w)
	l := int32(len(b))

	if (r > w && r+l < bp.size) || (r < w && w-r > l) {
		// Have more bytes ready than len(b), just do it
		copy(b, bp.buf[r:r+l])
		atomic.StoreInt32(&bp.r, r+l)
		bp.wc.L.Lock()
		bp.wc.Signal()
		bp.wc.L.Unlock()
		return int(l)
	}
	if r > w {
		// we have more bytes, but need to wrap around to finish reading.
		copy(b, bp.buf[r:bp.size])
		atomic.StoreInt32(&bp.r, 0)
		bp.wc.L.Lock()
		bp.wc.Signal()
		bp.wc.L.Unlock()
		return bp.Read(b[bp.size-r:]) + int(bp.size-r)
	}
	if r == w {
		// We dont have anything to read, wait for now
		bp.rc.L.Lock()
		if w == atomic.LoadInt32(&bp.w) {
			bp.rc.Wait()
		}
		bp.rc.L.Unlock()
		return bp.Read(b)
	}
	// Finally, this means we have less bytes available than size of b, copy what we have.
	copy(b, bp.buf[r:w])
	atomic.StoreInt32(&bp.r, w)
	bp.wc.L.Lock()
	bp.wc.Signal()
	bp.wc.L.Unlock()
	return int(w - r)
}
