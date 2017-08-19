//go:generate genny -in=colony.go -out=builtins.go gen "ValueType=BUILTINS"
package colony

import (
	"github.com/cheekybits/genny/generic"
	"sync"
	"unsafe"
)

type ValueType generic.Type

// ValueTypeColony represents a colony of ValueTypes.
type ValueTypeColony struct {
	entry *colonyGroupValueType
}

// NewValueTypeColony returns a new colony of ValueTypes.
func NewValueTypeColony(size uint) *ValueTypeColony {
	return &ValueTypeColony{
		entry: newValueTypeGroup(nil, size),
	}
}

// Insert returns a pointer from the colony and initializes it with the provided data.
func (c *ValueTypeColony) Insert(t *ValueType) (tp *ValueType) {
	return c.entry.Insert(t)
}

// Delete returns a pointer to the colony.
func (c *ValueTypeColony) Delete(tp *ValueType) {
	c.entry.Delete(tp)
}

func newValueTypeGroup(previous *colonyGroupValueType, size uint) *colonyGroupValueType {
	var g colonyGroupValueType
	if size == 0 {
		size = 8
	}
	g.data = make([]ValueType, size)
	g.free = make(chan *ValueType, size)
	for i := 0; i < len(g.data); i++ {
		g.free <- &g.data[i]
	}
	g.next = nil
	g.l = &sync.RWMutex{}
	g.previous = previous
	g.minPtr = uintptr(unsafe.Pointer(&g.data[0]))
	g.maxPtr = uintptr(unsafe.Pointer(&g.data[len(g.data)-1]))
	return &g
}

type colonyGroupValueType struct {
	data     []ValueType
	free     chan *ValueType
	maxPtr   uintptr
	minPtr   uintptr
	next     *colonyGroupValueType
	previous *colonyGroupValueType

	l *sync.RWMutex
}

func (g *colonyGroupValueType) Insert(t *ValueType) (tp *ValueType) {
	select {
	case tp = <-g.free:
		return
	default:
		if g.next == nil {
			g.next = newValueTypeGroup(g, uint(len(g.data)*2))
		}
		return g.next.Insert(t)
	}
}

func (g *colonyGroupValueType) Delete(tp *ValueType) {
	tpu := uintptr(unsafe.Pointer(tp))
	if tpu < g.minPtr || tpu > g.maxPtr { // hack to determine if a pointer points to this array
		g.next.Delete(tp)
	}
	g.free <- tp
	//if !g.index.Any() {
	// TODO: if a group has no more elements, then we should de-allocate it.
	//}
	return
}
