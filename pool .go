package uuid_pool

import (
	"github.com/google/uuid"
)

// UUIDPool maintains an internal pool
// google/uuid.UUID objects, UUIDPool
// should be created with NewUUIDPool
type UUIDPool struct {
	stop chan struct{}
	pool chan uuid.UUID
}

// NewUUIDPool returns a new UUIDPool whose intenal
// buffer capacity set to size
func NewUUIDPool(size uint) *UUIDPool {
	pool := &UUIDPool{
		pool: make(chan uuid.UUID, size),
		stop: make(chan struct{}),
	}
	pool.fillPool()
	go pool.watch()
	return pool
}

// fillPool initializes the pool to half capacity
func (p *UUIDPool) fillPool() {
	for i := 0; i < len(p.pool)/2; i++ {
		p.pool <- uuid.New()
	}
}

// watch makes sure that p.pool is always full
// it also looks out for p.stop
func (p *UUIDPool) watch() {
	for {
		select {
		// keep trying to send on p.pool
		case p.pool <- uuid.New():

		// when p.stop is closed, exit the loop
		case <-p.stop:
			p.pool = nil
			p.stop = nil
			return
		}
	}
}

// Get returns the uuid in front of the queue
func (p *UUIDPool) Get() uuid.UUID {
	return <-p.pool
}

// Dissolve causes the pool to stop generation
// and release internal pools
//
// After Dissolve is called, Get must not
// be called.
func (p *UUIDPool) Dissolve() {
	close(p.stop)
}
