package version_sort

import (
//	"sort"
	"strings"
//	"math"
	"strconv"
)

type ByVersion []string

func (va ByVersion) Len() int {
	return len(va)
}

func (va ByVersion) Swap(i, j int) {
	va[i], va[j] = va[j], va[i]
}

func intize(str string, ret chan int32) {
	defer close(ret)
	numstrs := strings.Split(str, ".")
	for _, v := range(numstrs) {
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return
		}
		ret<- int32(i)
	}
}

func (va ByVersion) Less(i, j int) bool {
	ic := make(chan int32)
	jc := make(chan int32)
	go intize(va[i], ic)
	go intize(va[j], jc)

	for {
		iv, iok := <-ic
		jv, jok := <-jc
		if !jok {
			// !jok comes first because if they
			// both have failed or if jok fails first,
			// iok is not less than jok
			// (if they fail at the same time, they are
			// equal)
			return false
		}
		if !iok {
			return true
		}
		if iv < jv {
			return true
		}
		if jv < iv {
			return false
		}
	}
	return false // arbitrary.  we can never reach this.
}