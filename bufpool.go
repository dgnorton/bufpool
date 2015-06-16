package bufpool

import (
	"container/list"
	"fmt"
	"sync"
)

type Pool struct {
	bufs *list.List
	in chan []byte
	out chan []byte
	stop chan struct{}
	wg sync.WaitGroup
	takes int
	gives int
}

func NewPool(bufsz, cnt int) *Pool {
	p := &Pool{
		bufs: list.New(),
		in: make(chan []byte),
		out: make(chan []byte),
		stop: make(chan struct{}),
	}

	for n := 0; n < cnt; n++ {
		p.bufs.PushBack(make([]byte, 0, bufsz))
	}

	p.wg.Add(1)

	go func() {
		defer p.wg.Done()
		for {
			if p.bufs.Len() == 0 {
				p.bufs.PushFront(make([]byte, 0, bufsz))
			}

			bo := p.bufs.Front()

			select {
			case <-p.stop:
				return
			case bi := <-p.in:
				p.bufs.PushFront(bi)
			case p.out <- bo.Value.([]byte):
				p.bufs.Remove(bo)
			}
		}
	}()

	return p
}

func (p *Pool) Take() []byte {
	p.takes++
	return <-p.out
}

func (p *Pool) Give(b []byte) {
	p.gives++
	p.in <- b
}

func (p *Pool) Stats() string {
	stats := fmt.Sprintf("len = %d, takes = %d, gives = %d", p.bufs.Len(), p.takes, p.gives)
	return stats
}

func (p *Pool) Close() {
	close(p.stop)
	p.wg.Wait()
}