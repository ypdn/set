package set

import "testing"

func TestStuff(test *testing.T) {
	s := New[string]()
	s.Put("zürich")
	s.Put("sofia")
	s.Put("dresden")
	s.Put("turin")
	s.Put("lyon")
	s.Put("london")

	t := New[string]()
	t.Put("york")
	t.Put("zürich")
	t.Put("turin")
	t.Put("london")

	u := New[string]()
	u.Put("bern")
	u.Put("london")

	i := Intersection(s, t, u)
	test.Log(i)
	if !i.Has("london") || i.Len() != 1 {
		test.Fail()
	}

	test.Log(Union(s, t, u))
	test.Log(Intersection(s, t))
	test.Log(s.Difference(t))
	if !s.Equal(s.Copy()) {
		test.Fail()
	}
	if s.Equal(i) {
		test.Fail()
	}
	ct := t.Copy()
	ct.Remove("york")
	if !ct.Subset(s) {
		test.Fail()
	}
}
