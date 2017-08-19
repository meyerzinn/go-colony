package colony

import (
	"github.com/cheekybits/genny/generic"
	"github.com/willf/bitset"
	"sync"
	"unsafe"
)

type Type generic.Type

type TypeColony struct {
	entry *colonyGroupType
}

// NewTypeColony returns a new colony of Type's.
func NewTypeColony() *TypeColony {
	return &TypeColony{
		entry: newTypeGroup(nil),
	}
}

// Iterate sends pointers to all instances of Type in the colony to the given channel.
func (c *TypeColony) Iterate() <-chan *Type {
	ch := make(chan *Type)
	var wg sync.WaitGroup
	for g := c.entry; g != nil; g = g.next {
		wg.Add(1)
		go func(g *colonyGroupType) {
			g.l.RLock()
			for i, e := g.index.NextSet(0); e; i, e = g.index.NextSet(i + 1) {
				ch <- &g.data[i]
			}
			g.l.RUnlock()
			wg.Done()
		}(g)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}

func (c *TypeColony) Insert(t *Type) (tp *Type) {
	return c.entry.Insert(t)
}

func (c *TypeColony) Delete(tp *Type) {
	c.entry.Delete(tp)
}

func newTypeGroup(previous *colonyGroupType) *colonyGroupType {
	var g colonyGroupType
	if previous != nil {
		g.data = make([]Type, len(previous.data)*2)
		g.index = bitset.New(uint(len(previous.data) * 2))
	} else {
		g.data = make([]Type, 2)
		g.index = bitset.New(2)
	}
	g.next = nil
	g.l = &sync.RWMutex{}
	g.minPtr = uintptr(unsafe.Pointer(&g.data[0]))
	g.maxPtr = uintptr(unsafe.Pointer(&g.data[len(g.data)-1]))
	return &g
}

type colonyGroupType struct {
	data     []Type
	index    *bitset.BitSet
	maxPtr   uintptr
	minPtr   uintptr
	next     *colonyGroupType
	previous *colonyGroupType

	l *sync.RWMutex
}

func (g *colonyGroupType) Insert(t *Type) (tp *Type) {
	g.l.Lock()
	if i, e := g.index.NextClear(0); e {
		g.data[i] = *t
		g.index.Set(i)
		tp = &g.data[i]
		g.l.Unlock()
		return
	}
	if g.next == nil {
		g.next = newTypeGroup(g)
	}
	g.l.Unlock()
	return g.next.Insert(t)
}

func (g *colonyGroupType) Delete(tp *Type) {
	if uintptr(unsafe.Pointer(tp)) > g.maxPtr { // hack to determine if a pointer points to this array
		g.next.Delete(tp)
	}
	g.l.Lock()
	for i := 0; i < len(g.data); i++ {
		if tp == &g.data[i] {
			g.index.Clear(uint(i))
			//if !g.index.Any() {
			// TODO: if a group has no more elements, then we should de-allocate it.
			//}
			g.l.Unlock()
			return
		}
	}
	g.l.Unlock()
	if g.next != nil {
		g.next.Delete(tp)
	}
}
