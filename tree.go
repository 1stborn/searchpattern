package searchpattern

import "strings"

type radixTree struct {
	leafs         radix
	caseSensitive bool
	value         interface{}
}

type radix map[uint8]*radixTree

const (
	skipOne      = '?'
	skipInfinite = '*'
)

func CaseSensitive() *radixTree {
	return &radixTree{caseSensitive: true}
}

func CaseInsensitive() *radixTree {
	return new(radixTree)
}

func (r *radixTree) Add(pattern string, v interface{}) {
	if r.caseSensitive {
		r.add(pattern, v)
	} else {
		r.add(strings.ToLower(pattern), v)
	}
}

func (r *radixTree) add(search string, v interface{}) {
	if search == "*" || search == "" {
		r.value = v

		return
	} else if r.leafs == nil {
		r.leafs = make(radix)
	}

	c := search[0]

	if v, ok := r.leafs[c]; ok {
		v.add(search[1:], v)
	} else {
		rt := new(radixTree)
		rt.add(search[1:], v)
		r.leafs[c] = rt
	}
	r.leafs[c] = r.leafs[c]
}

func (r *radixTree) Find(search string) interface{} {
	if len(search) == 0 {
		return nil
	} else if f := r.find(search); len(f) > 0 {
		return f[0]
	}

	return nil
}

func (r *radixTree) find(search string) (found []interface{}) {
	 if r.leafs == nil {
		found = append(found, r.value)
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
