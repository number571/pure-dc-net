package syncmap

import "sync"

type syncMap struct {
	mp map[string]byte
	mt *sync.RWMutex
}

func NewSyncMap() ISyncMap {
	return &syncMap{
		mp: make(map[string]byte),
		mt: &sync.RWMutex{},
	}
}

func (p *syncMap) Store(k string, v byte) {
	p.mt.Lock()
	defer p.mt.Unlock()

	p.mp[k] = v
}

func (p *syncMap) Size() int {
	p.mt.RLock()
	defer p.mt.RUnlock()

	return len(p.mp)
}

func (p *syncMap) Sum() byte {
	p.mt.RLock()
	defer p.mt.RUnlock()

	s := byte(0)
	for _, b := range p.mp {
		s ^= b
	}
	return s
}

func (p *syncMap) Clear() {
	p.mt.Lock()
	defer p.mt.Unlock()

	p.mp = make(map[string]byte)
}
