package forestGen

import (
	"fmt"
	"testing"
)

func TestGen(t *testing.T) {
	m := Generate(100, 100)
	fmt.Print(m.String())
}
