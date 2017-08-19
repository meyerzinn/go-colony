//go:generate genny -in=colony_test.go -out=builtins_test.go gen "ValueType=BUILTINS"
package colony

import "testing"

func TestValueTypeColony(t *testing.T) {
	colony := NewValueTypeColony(1)
	t.Run("Insert", func(t *testing.T) {
		newT := new(ValueType)
		tp := colony.Insert(newT)
		if *tp != *newT {
			t.Fatalf("value of pointer returned from Insert does not equal the inserted value: (*newT) %v != (*tp) %v", *newT, *tp)
		}
	})
	t.Run("Delete", func(t *testing.T) {
		newT := new(ValueType)
		tp := colony.Insert(newT)
		colony.Delete(tp)
		tp2 := colony.Insert(newT)
		if tp != tp2 { // both should be allocated to the same spot
			t.Fatalf("deletion failed")
		}
	})
}

var ValueTypeBenchmarks = []struct {
	name  string
	count uint
}{
	{"1", 1},
	{"10", 10},
	{"100", 100},
	{"1000", 1000},
	{"10000", 10000},
	{"100000", 100000},
	{"1000000", 1000000},
}

func BenchmarkValueTypeColony_Insert(b *testing.B) {
	for _, bm := range ValueTypeBenchmarks {
		b.Run(bm.name, func(count uint) func(*testing.B) {
			return func(b *testing.B) {
				// setup
				colony := NewValueTypeColony(count)
				newValueType := new(ValueType)

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					for j := 0; uint(j) < count; j++ {
						newValueType = colony.Insert(newValueType)
					}
				}
			}
		}(bm.count))
	}
}

func BenchmarkValueTypeSlice(b *testing.B) {
	for _, bm := range ValueTypeBenchmarks {
		b.Run(bm.name, func(count uint) func(*testing.B) {
			return func(b *testing.B) {
				// setup
				 arr := make([]ValueType, count)
				newValueType := new(ValueType)

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					for j := 0; uint(j) < count; j++ {
						arr = append(arr, *newValueType)
					}
				}
			}
		}(bm.count))
	}
}
