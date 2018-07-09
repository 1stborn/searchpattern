package searchpattern

import (
	"sort"
	"strings"
)

type RadixTree struct {
	leafs         radix
	caseSensitive bool
	value         interface{}
}

type radix map[rune]*RadixTree

type result struct {
	weight int
	value   interface{}
}

type results []result

func (r results) Len() int {
	return len(r)
}

func (r results) Less(i, j int) bool {
	return r[i].weight < r[j].weight
}

func (r results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

const (
	skipOne      = '?'
	skipInfinite = '*'
)

func CaseSensitive() *RadixTree {
	return &RadixTree{caseSensitive: true}
}

func CaseInsensitive() *RadixTree {
	return new(RadixTree)
}

func (r *RadixTree) Add(pattern string, v interface{}) {
	if r.caseSensitive {
		r.add([]rune(pattern), v)
	} else {
		r.add([]rune(strings.ToLower(pattern)), v)
	}
}

func (r *RadixTree) add(search []rune, v interface{}) {
	if l := len(search); l == 0 {
		r.value = v
		return
	} else if l == 1 && search[0] == '*' {
		r.value = v
		return
	} else if r.leafs == nil {
		r.leafs = make(radix)
	}

	c := []rune(search)[0]

	if l, ok := r.leafs[c]; ok {
		l.add(search[1:], v)
	} else {
		rt := &RadixTree{caseSensitive: r.caseSensitive}
		rt.add(search[1:], v)
		r.leafs[c] = rt
	}

	r.leafs[c] = r.leafs[c]
}

func (r *RadixTree) Find(search string) interface{} {
	if len(search) == 0 {
		return nil
	} else if !r.caseSensitive {
		search = strings.ToLower(search)
	}

	f := r.find([]rune(search), 0)

	switch len(f) {
	case 0:
		return nil
	case 1:
		return f[0].value
	default:
		sort.Sort(f)
		return f[0].value
	}
}

func (r *RadixTree) FindFirst(search string) interface{} {
	if len(search) == 0 {
		return nil
	} else if !r.caseSensitive {
		search = strings.ToLower(search)
	}

	if f := r.find([]rune(search), 0); len(f) > 0 {
		return f[0].value
	}

	return nil
}

func (r *RadixTree) FindAll(search string, fn func(v interface{})) {
	if len(search) == 0 {
		return
	} else if !r.caseSensitive {
		search = strings.ToLower(search)
	}

	for _, v := range r.find([]rune(search), 0) {
		fn(v.value)
	}
}

func (r *RadixTree) find(search []rune, weight int) (found results) {
	if len(search) == 0 {
		return
	} else if r.leafs == nil {
		found = append(found, result{
			weight:  weight,
			value:   r.value,
		})
	}

	var current = search[0]

	if leaf, ok := r.leafs[current]; ok {
		found = append(found, leaf.find(search[1:], weight + 1)...)
	}

	if leaf, ok := r.leafs[skipOne]; ok {
		found = append(found, leaf.find(search[1:], weight - 1)...)
	}

	if leaf, ok := r.leafs[skipInfinite]; ok {
		var skip = 0
		for i := 0; i < len(search); i++ {
			current := search[i]
			skip++
			if leaf, ok := leaf.leafs[current]; ok {
				found = append(found, leaf.find(search[i+1:], weight - skip)...)
			}
		}
	}

	return
}
