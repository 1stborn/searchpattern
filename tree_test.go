package searchpattern_test

import "testing"
import "github.com/1stborn/searchpattern"

func TestRadixTree(t *testing.T) {
	rtree := searchpattern.CaseInsensitive()

	for i, v := range []string{"aol.com", "love.com", "y?m.com", "g*mes.com", "wow.com"} {
		rtree.Add(v, i)
	}

	if v, ok := rtree.Find("love.com").(int); ok {
		if v != 1 {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	if v, ok := rtree.Find("yim.com").(int); ok {
		if v != 2 {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	if v, ok := rtree.Find("grames.com").(int); ok {
		if v != 3 {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	if _, ok := rtree.Find("aoling.com").(int); ok {
		t.Fail()
	}
}
