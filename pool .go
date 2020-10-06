package uuid_pool

import (
	"github.com/google/uuid"
)

type UUIDPool struct {
	stop chan struct{}
	pool chan uuid.UUID
}

func NewUUIDPool(size uint) *UUIDPool {
	pool := &UUIDPool{
		pool: make(chan uuid.UUID, size),
		stop: make(chan struct{}, 1),
	}
	pool.fillPool()
	go pool.watch()
	return pool
}

func (p *UUIDPool) fillPool() {
	for i := 0; i < len(p.pool); i++ {
		p.pool <- uuid.New()
	}
}

func (p *UUIDPool) watch() {
	for {
		select {
		case <-p.stop:
			p.pool = nil
			p.stop = nil
			return
		default:
			// code here will be stuck trying to send on p.pool
			// until a uuid leaves the channel p.stop will never
			// be checked, fix this
			p.pool <- uuid.New()
		}
	}
}

func (p *UUIDPool) Get() uuid.UUID {
	return <-p.pool
}

func (p UUIDPool) Dissolve() {
	p.stop <- struct{}{}
	return
}
