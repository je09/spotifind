package spotifind

import (
	"sort"
)

// Pair is a data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value int
}

// PairList is a slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// SortStyleMap sorts a map by value, and returns a slice of keys in the sorted order.
func SortStyleMap(styleMap map[string]int) []string {
	p := make(PairList, len(styleMap))
	i := 0
	for k, v := range styleMap {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))

	sortedStyles := make([]string, len(styleMap))
	for i, pair := range p {
		sortedStyles[i] = pair.Key
	}

	return sortedStyles
}
