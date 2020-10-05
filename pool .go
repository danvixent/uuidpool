package uuid_pool

import (
	"github.com/google/uuid"
	"sync"
)

type UUIDPool struct {
	get  sync.Mutex
	repl sync.Mutex
	pool []uuid.UUID
	off  uint16
}

func NewUUIDPool(size uint16) *UUIDPool {
	pool := &UUIDPool{
		pool: make([]uuid.UUID, size),
		off:  size - 1,
	}
	pool.fillPool()
	return pool
}

func (p *UUIDPool) Get() uuid.UUID {
	p.get.Lock()
	u := p.pool[p.off]
	go p.replace(p.off)

	if p.off == 0 {
		p.off = uint16(len(p.pool))
	}

	p.off--
	p.get.Unlock()
	return u
}

func (p *UUIDPool) fillPool() {
	for i := range p.pool {
		p.pool[i] = uuid.New()
	}
}

func (p *UUIDPool) replace(pos uint16) {
	p.pool[pos] = uuid.New()
}
