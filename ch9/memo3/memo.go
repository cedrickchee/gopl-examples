// Package memo provides a concurrency-safe memoization a function of
// type Func. Requests for different keys run concurrently.
// Concurrent requests for the same key result in duplicate work.
package memo

import "sync"

// A concurrent non-blocking cache.
// This is the problem of _memoizing_ a function, that is, caching the result of
// a function so that it need be computed only once. Our solution will be
// concurrency-safe and will avoid the contention associated with designs based
// on a single lock for the whole cache.

// A Memo caches the results of calling a Func.
type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]result
}

// Func is the type of the function to memoize.
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// Get is concurrency-safe.
// In this version of implementation, the calling goroutine acquires the lock
// twice: once for the lookup, and then a second time for the update if the
// lookup returned nothing. In between, other goroutines are free to use the
// cache.
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	memo.mu.Unlock()
	if !ok {
		res.value, res.err = memo.f(key)

		// Between the two critical sections, several goroutines
		// may race to compute f(key) and update the map.
		memo.mu.Lock()
		memo.cache[key] = res
		memo.mu.Unlock()
	}
	return res.value, res.err
}

// Notes
//
// A `Memo` instance holds the function `f` to memoize, of type `Func`, and the
// cache, which is a mapping from strings to `results`. Each result is simply
// the pair of results returned by a call to `f`--a value and an error. We’ll
// show several variations of `Memo` as the design progresses, but all will
// share these basic aspects.
