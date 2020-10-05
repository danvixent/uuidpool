package uuid_pool

import (
	"github.com/google/uuid"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestNewUUIDPool(t *testing.T) {
	type args struct {
		size uint16
	}
	tests := []struct {
		name string
		args args
		want *UUIDPool
	}{
		{
			"new1",
			args{16},
			&UUIDPool{
				pool: make([]uuid.UUID, 16),
				off:  15,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUUIDPool(tt.args.size)

			// they mustn't be equal, because the got's pool should have been filled
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("want pool length %v, got pool length %v", len(tt.want.pool), len(got.pool))
			}

			if tt.want.off != got.off {
				t.Errorf("want offset %d, got offset %d", tt.want.off, got.off)
			}
		})
	}
}

func TestUUIDPool_Get(t *testing.T) {
	pool := NewUUIDPool(4)
	got, want := pool.Get(), uuid.UUID{}
	if got == want {
		t.Errorf("empty uuid in pool: %v", got)
	}

	time.Sleep(time.Second)
	if pool.pool[pool.off+1] == got {
		t.Errorf("uuid not automatically replaced: %v", pool.pool[pool.off+1])
	}
}

func BenchmarkUUIDPool_Get(b *testing.B) {
	pool := NewUUIDPool(math.MaxUint16)

	var uu uuid.UUID
	for i := 0; i < b.N; i++ {
		uu = pool.Get()
	}
	_ = uu.String()
}
