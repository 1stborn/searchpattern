package searchpattern

import "strings"

type RadixTree struct {
	leafs         radix
	caseSensitive bool
	value         interface{}
}

type radix map[uint8]*RadixTree

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
		r.add(pattern, v)
	} else {
		r.add(strings.ToLower(pattern), v)
	}
}

func (r *RadixTree) add(search string, v interface{}) {
	if search == "*" || search == "" {
		r.value = v

		return
	} else if r.leafs == nil {
		r.leafs = make(radix)
	}

	c := search[0]

	if l, ok := r.leafs[c]; ok {
		l.add(search[1:], v)
	} else {
		rt := new(RadixTree)
		rt.add(search[1:], v)
		r.leafs[c] = rt
	}

	r.leafs[c] = r.leafs[c]
}

func (r *RadixTree) Find(search string) interface{} {
	if len(search) == 0 {
		return nil
	} else if f := r.find(search); len(f) > 0 {
		return f[0]
	}

	return nil
}

func (r *RadixTree) find(search string) (found []interface{}) {
	if len(search) == 0 {
		if r.leafs == nil {
			return append(found, r.value)
		} else {
			return
		}
	}

	if !r.caseSensitive {
		search = strings.ToLower(search)
	}

	current := search[0]

	if leaf, ok := r.leafs[current]; ok {
		found = append(found, leaf.find(search[1:])...)
	}

	if leaf, ok := r.leafs[skipOne]; ok {
		found = append(found, leaf.find(search[1:])...)
	}

	if leaf, ok := r.leafs[skipInfinite]; ok {
		for i := 0; i < len(search); i++ {
			current := search[i]
			if leaf, ok := leaf.leafs[current]; ok {
				found = append(found, leaf.find(search[i+1:])...)
			}
		}
	}

	return
}