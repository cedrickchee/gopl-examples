// Package memo provides a concurrency-safe memoization a function of
// a function. Requests for different keys proceed in parallel.
// Concurrent requests for the same key block until the first completes.
// This implementation uses a Mutex.
package memo

import "sync"

// A concurrent non-blocking cache.
// This is the problem of _memoizing_ a function, that is, caching the result of
// a function so that it need be computed only once. Our solution will be
// concurrency-safe and will avoid the contention associated with designs based
// on a single lock for the whole cache.

// In the version of `Memo` below, each map element is a pointer to an `entry`
// struct. Each `entry` contains the memoized result of a call to the function
// `f`, as before, but it additionally contains a channel called `ready`. Just
// after the `entry`'s `result` has been set, this channel will be closed, to
// _broadcast_ to any other goroutines that it is now safe for them to
// read the result from the `entry`.

// A Memo caches the results of calling a Func.
type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]*entry
}

// Func is the type of the function to memoize.
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

// A call to Get now involves acquiring the mutex lock that guards the `cache`
// map, looking in the map for a pointer to an existing `entry`, allocating and
// inserting a new `entry` if none was found, then releasing the lock. If there
// was an existing `entry`, its value is not necessarily ready yet--another
// goroutine could still be calling the slow function `f`--so the calling
// goroutine must wait for the `entry`'s "ready" condition before it reads the
// `entry`'s `result`. It does this by reading a value from the `ready` channel,
// since this operation blocks until the channel is closed.
//
// If there was no existing `entry`, then by inserting a new "not ready" `entry`
// into the map, the current goroutine becomes responsible for invoking the slow
// function, updating the `entry`, and broadcasting the readiness of the new
// `entry` to any other goroutines that might (by then) be waiting for it.
//
// Notice that the variables `e.res.value` and `e.res.err` in the `entry` are
// shared among multiple goroutines. The goroutine that creates the `entry` sets
// their values, and other goroutines read their values once the "ready"
// condition has been broadcast. Despite being accessed by multiple goroutines,
// no mutex lock is necessary. The closing of the `ready` channel _happens
// before_ any other goroutine receives the broadcast event, so the write to
// those variables in the first goroutine _happens before_ they are read by
// subsequent goroutines. There is no data race.
func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// This is the first request for this key.
		// This goroutine becomes responsible for computing
		// the value and broadcasting the ready condition.
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()

		e.res.value, e.res.err = memo.f(key)

		close(e.ready) // broadcast ready condition
	} else {
		// This is a repeat request for this key.
		memo.mu.Unlock()

		<-e.ready // wait for ready condition
	}
	return e.res.value, e.res.err
}

// Notes
//
// A `Memo` instance holds the function `f` to memoize, of type `Func`, and the
// cache, which is a mapping from strings to `results`. Each result is simply
// the pair of results returned by a call to `f`--a value and an error. We’ll
// show several variations of `Memo` as the design progresses, but all will
// share these basic aspects.
