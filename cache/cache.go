package cache

import (
	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/store"
)

var Cache *cache.Cache
var BareCache *ristretto.Cache

func Init() {
	var err error

	BareCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})

	if err != nil {
		panic(err)
	}

	ristrettoStore := store.NewRistretto(BareCache, &store.Options{
		Cost: 1,
	})

	Cache = cache.New(ristrettoStore)
}
