package binpack

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkKnapsackPack(b *testing.B) {
	for _, count := range []int{20, 84, 239} {
		b.Run(fmt.Sprintf("items=%d", count), func(b *testing.B) {
			items := make([]Item[int], count)
			bins := make([]*Bin[int, int], count)
			for i := range count {
				length := float64(1000 + (i%6)*1000)
				items[i] = Item[int]{Original: i, Length: length}
				bins[i] = NewBin[int, int](i, 12000)
			}

			algorithm := NewKnapsackAlgorithm[int, int]()
			b.ReportAllocs()
			for b.Loop() {
				for _, bin := range bins {
					bin.Reset()
				}
				if _, err := algorithm.Pack(context.Background(), items, bins); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
