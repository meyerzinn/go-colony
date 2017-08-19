package colony

import "testing"

func TestTypeColony(t *testing.T) {
	colony := NewTypeColony()
	t.Run("Insert", func(t *testing.T) {
		newT := new(Type)
		tp := colony.Insert(newT)
		if *tp != *newT {
			t.Fatalf("value of pointer returned from Insert does not equal the inserted value: (*newT) %v != (*tp) %v", *newT, *tp)
		}
	})
	t.Run("Delete", func(t *testing.T) {
		newT := new(Type)
		tp := colony.Insert(newT)
		colony.Delete(tp)
		tp2 := colony.Insert(newT)
		if tp != tp2 { // both should be allocated to the same spot
			t.Fatalf("deletion failed")
		}
	})
}

func BenchmarkTypeColony_Insert(b *testing.B) {
	benchmarks := []struct {
		name  string
		count int
	}{
		{"1", 1},
		{"10", 10},
		{"100", 100},
		{"1000", 1000},
		{"10000", 10000},
	}
	colony := NewTypeColony()
	newT := new(Type)
	b.ReportAllocs()
	for _, bm := range benchmarks {
		b.Run(bm.name, func(count int) func(*testing.B) {
			return func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					for j := 0; j < count; j++ {
						newT = colony.Insert(newT)
					}
				}
			}
		}(bm.count))
	}
}
