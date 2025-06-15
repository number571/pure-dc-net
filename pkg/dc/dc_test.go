package dc

import (
	"testing"
)

func TestDCNetGenerate(t *testing.T) {
	t.Parallel()

	var (
		g12 = NewHGenerator([]byte("secret-key-1"))
		g23 = NewHGenerator([]byte("secret-key-2"))
		g31 = NewHGenerator([]byte("secret-key-3"))
	)

	node1 := NewDCNet(0, g12, g31)
	node2 := NewDCNet(0, g12, g23)
	node3 := NewDCNet(0, g23, g31)

	for range 100 {
		result := node1.Generate() ^ node2.Generate() ^ node3.Generate()
		if result != 0 {
			t.Error("generate failed without msg")
			break
		}
	}

	msg := byte(0x71)
	for range 100 {
		result := (msg ^ node1.Generate()) ^ node2.Generate() ^ node3.Generate()
		if result != msg {
			t.Error("generate failed with msg")
			break
		}
	}
}
