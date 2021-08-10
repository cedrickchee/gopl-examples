// Package memo provides a concurrency-unsafe
// memoization of a function of type Func.
package memo

// A concurrent non-blocking cache.
// This is the problem of _memoizing_ a function, that is, caching the result of
// a function so that it need be computed only once. Our solution will be
// concurrency-safe and will avoid the contention associated with designs based
// on a single lock for the whole cache.

// This version is the first draft of the cache.

// A Memo caches the results of calling a Func.
type Memo struct {
	f     Func
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

// NOTE: not concurrency-safe!
func (memo *Memo) Get(key string) (interface{}, error) {
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
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
