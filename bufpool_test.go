package bufpool

import (
	"testing"
)

func TestNewPool(t *testing.T) {
	t.Parallel()

	p := NewPool(0, 0)
	if p == nil {
		t.Error("NewPool(0, 0) failed")
	}
	p.Close()

	p = NewPool(0, 1)
	if p == nil {
		t.Error("NewPool(0, 1) failed")
	}
	p.Close()

	p = NewPool(1000, 1000)
	if p == nil {
		t.Error("NewPool(1000, 1000) failed")
	}
	p.Close()
}

func TestTakeGive(t *testing.T) {
	t.Parallel()

	p := NewPool(0, 0)

	b1 := p.Take()
	if cap(b1) != 0 {
		t.Errorf("\n\texp = %d\n\tgot = %d\n", 0, cap(b1))
	}

	b2 := p.Take()
	if cap(b2) != 0 {
		t.Errorf("\n\texp = %d\n\tgot = %d\n", 0, cap(b2))
	}

	p.Give(b1)
	p.Give(b2)
	p.Close()

	p = NewPool(1000, 0)

	b1 = p.Take()
	if cap(b1) != 1000 {
		t.Errorf("\n\texp = %d\n\tgot = %d\n", 0, cap(b1))
	}

	p.Give(b1)
	p.Close()

	p = NewPool(1000, 100)
	bufs := [][]byte{}
	for i := 0; i < 100; i++ {
		bufs = append(bufs, p.Take())
	}
	p.Close()
}