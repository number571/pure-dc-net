package dc

import "sync"

type totalizer struct {
	mtx   *sync.Mutex
	blist []byte
}

func NewTotalizer() ITotalizer {
	return &totalizer{
		mtx:   &sync.Mutex{},
		blist: make([]byte, 0, 256),
	}
}

func (p *totalizer) Store(b ...byte) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.blist = append(p.blist, b...)
}

func (p *totalizer) Size() int {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	return len(p.blist)
}

func (p *totalizer) Sum() byte {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	s := byte(0)
	for _, bx := range p.blist {
		s ^= bx
	}

	p.blist = p.blist[:0]
	return s
}
