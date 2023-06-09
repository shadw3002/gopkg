package lscq

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestCacheRemap16Byte(t *testing.T) {
	t.Parallel()

	t.Run("scqsize must be 2^n", func(t *testing.T) {
		if scqsize&(scqsize-1) != 0 {
			panic("unexpected")
		}
	})

	t.Run("remapStep must be 2*k + 1", func(t *testing.T) {
		if !(1 < remapStep && remapStep%2 == 1) {
			panic("unexpected")
		}
	})

	t.Run("remapStep must be large enough", func(t *testing.T) {
		if !(remapStep*unsafe.Sizeof(*new(scqNodePointer)) >= cacheLineSize) {
			panic("unexpected")
		}
	})

	t.Run("one to one map", func(t *testing.T) {
		visited := make([]bool, scqsize)
		for i := uint64(0); i < scqsize; i++ {
			visited[cacheRemap16Byte(i)] = true
		}
		for i := uint64(0); i < scqsize; i++ {
			if !visited[i] {
				panic("unexpected")
			}
		}
	})

	t.Run("shuffle cacheline access", func(t *testing.T) {
		var (
			pre     = -1
			preline = -1
			nowline = -1
		)
		for i := 0; i < scqsize+1; i++ {
			line := int(cacheRemap16Byte(uint64(0)) / 4)
			nowline = line
			if preline == nowline {
				fmt.Printf("%d %d is in the same line.", pre, i)
			}
			pre = i
			preline = nowline
		}
	})
}
