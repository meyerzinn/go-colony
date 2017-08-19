//go:generate genny -in=colony.go -out=builtins.go gen "ValueType=complex64"
package colony

import (
	"github.com/cheekybits/genny/generic"
	"sync"
	"unsafe"
)

type ValueType generic.Type

type ValueTypeColony struct {
	entry *colonyGroupValueType
}

// NewValueTypeColony returns a new colony of ValueType's.
func NewValueTypeColony(size uint) *ValueTypeColony {
	return &ValueTypeColony{
		entry: newValueTypeGroup(nil, size),
	}
}

//// Iterate sends pointers to all instances of ValueType in the colony to the given channel.
//func (c *ValueTypeColony) Iterate() <-chan *ValueType {
//	ch := make(chan *ValueType)
//	var wg sync.WaitGroup
//	for g := c.entry; g != nil; g = g.next {
//		wg.Add(1)
//		go func(g *colonyGroupValueType) {
//			g.l.RLock()
//			//for i, e := g.index.NextSet(0); e; i, e = g.index.NextSet(i + 1) {
//			//	ch <- &g.data[i]
//			//}
//			for _, f := range g.free {
//
//			}
//			g.l.RUnlock()
//			wg.Done()
//		}(g)
//	}
//	go func() {
//		wg.Wait()
//		close(ch)
//	}()
//	return ch
//}

func (c *ValueTypeColony) Insert(t *ValueType) (tp *ValueType) {
	return c.entry.Insert(t)
}

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
	//if i, e := g.index.NextClear(0); e {
	//	g.data[i] = *t
	//	g.index.Set(i)
	//	tp = &g.data[i]
	//	g.l.Unlock()
	//	return
	//}
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
