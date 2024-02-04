// The set package provides a concurrency-safe set implementation.
package set

import (
	"fmt"
	"strings"
	"sync"
)

// Set implements a set.
type Set[T comparable] struct {
	m  map[T]struct{}
	mu sync.RWMutex
}

// New returns a new Set.
func New[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]struct{})}
}

// Intersection returns the intersection of sets.
func Intersection[T comparable](sets ...*Set[T]) *Set[T] {
	i := New[T]()
	if len(sets) == 0 {
		return i
	}
	smallest := sets[0]

	for _, s := range sets {
		s.mu.RLock()
		defer s.mu.RUnlock()

		if len(s.m) < len(smallest.m) {
			smallest = s
		}
	}

next:
	for e := range smallest.m {
		for _, s := range sets {
			if _, ok := s.m[e]; !ok {
				continue next
			}
		}
		i.m[e] = struct{}{}
	}
	return i
}

// Union returns the union of sets.
func Union[T comparable](sets ...*Set[T]) *Set[T] {
	u := New[T]()
	for _, s := range sets {
		s.Range(func(e T) bool {
			u.m[e] = struct{}{}
			return true
		})
	}
	return u
}

// Put puts e into s.
func (s *Set[T]) Put(e T) {
	s.mu.Lock()
	s.m[e] = struct{}{}
	s.mu.Unlock()
}

// Delete deletes e from s.
func (s *Set[T]) Delete(e T) {
	s.mu.Lock()
	delete(s.m, e)
	s.mu.Unlock()
}

// Has returns true if s has e.
func (s *Set[T]) Has(e T) bool {
	s.mu.RLock()
	_, ok := s.m[e]
	s.mu.RUnlock()
	return ok
}

// Len returns the number of members of s.
func (s *Set[T]) Len() int {
	s.mu.RLock()
	l := len(s.m)
	s.mu.RUnlock()
	return l
}

// Range calls f for each member of s, so long as it returns true.
func (s *Set[T]) Range(f func(T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for e := range s.m {
		if !f(e) {
			break
		}
	}
}

// Copy returns a copy of s.
func (s *Set[T]) Copy() *Set[T] {
	c := New[T]()
	s.Range(func(e T) bool {
		c.m[e] = struct{}{}
		return true
	})
	return c
}

// Diff returns the result of s-t.
func (s *Set[T]) Diff(t *Set[T]) *Set[T] {
	t.mu.RLock()
	defer t.mu.RUnlock()

	d := New[T]()
	s.Range(func(e T) bool {
		if _, ok := t.m[e]; !ok {
			d.m[e] = struct{}{}
		}
		return true
	})
	return d
}

// Subset returns true if s is a subset of t.
func (s *Set[T]) Subset(t *Set[T]) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t.mu.RLock()
	defer t.mu.RUnlock()

	for e := range s.m {
		if _, ok := t.m[e]; !ok {
			return false
		}
	}
	return true
}

// Equal returns true if s and t have the same members.
func (s *Set[T]) Equal(t *Set[T]) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t.mu.RLock()
	defer t.mu.RUnlock()

	if len(s.m) != len(t.m) {
		return false
	}
	for e := range s.m {
		if _, ok := t.m[e]; !ok {
			return false
		}
	}
	return true
}

// String returns the string representation of s.
func (s *Set[T]) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var sb strings.Builder
	sb.WriteRune('{')

	var i int
	for e := range s.m {
		i++
		fmt.Fprint(&sb, e)
		if i < len(s.m) {
			sb.WriteString(", ")
		}
	}
	sb.WriteRune('}')
	return sb.String()
}
