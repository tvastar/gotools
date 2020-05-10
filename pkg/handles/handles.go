// Package handles provides a simple handle-table implementation for
// use with cgo
package handles

import "sync"

// Table tracks arbitrary values by assigning them handles (unique IDs).
type Table struct {
	sync.Mutex
	max   uintptr
	items map[uintptr]interface{}
}

// Add adds an arbitrary item into the table returning a handle.
func (t *Table) Add(v interface{}) uintptr {
	t.Lock()
	defer t.Unlock()

	if t.items == nil {
		t.items = map[uintptr]interface{}{}
	}

	result := t.max
	t.max++
	t.items[result] = v
	return result
}

// Get returns the item for a handle.
func (t *Table) Get(h uintptr) (v interface{}, ok bool) {
	t.Lock()
	defer t.Unlock()

	v, ok = t.items[h]
	return v, ok
}

// Delete removes a handle from the table.
func (t *Table) Delete(h uintptr) bool {
	t.Lock()
	defer t.Unlock()
	_, exists := t.items[h]
	delete(t.items, h)
	return exists
}

// Size returns the number of handles.
func (t *Table) Size() int {
	t.Lock()
	defer t.Unlock()
	return len(t.items)
}
