package treesort_test

import (
	"math/rand"
	"sort"
	"testing"

	"gopl.io/ch4/treesort"
)

// Run:
// $ go test -v -run Sort gopl.io/ch4/treesort
func TestSort(t *testing.T) {
	data := make([]int, 50)
	for i := range data {
		data[i] = rand.Int() % 50
	}
	t.Logf("not sorted data: %v", data)
	treesort.Sort(data)
	if !sort.IntsAreSorted(data) {
		t.Errorf("not sorted: %v", data)
	}
	t.Logf("sorted data: %v", data)
}
