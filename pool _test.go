package uuidpool

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestNewUUIDPool(t *testing.T) {
	type args struct {
		size uint
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
				pool: make(chan uuid.UUID, 16),
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
		})
	}
}

func TestUUIDPool_Get(t *testing.T) {
	pool := NewUUIDPool(4)
	got, want := pool.Get(), uuid.UUID{}
	if got == want {
		t.Errorf("empty uuid in pool: %v", got)
	}
}

func BenchmarkUUIDPool_Get(b *testing.B) {
	var pool = NewUUIDPool(10)
	for i := 0; i < b.N; i++ {
		pool.Get()
	}
}
