package directedPath

import (
	"fmt"
	"testing"
)

func TestGen(t *testing.T) {
	m := Generate(150, 75)
	fmt.Print(m.String())
}
