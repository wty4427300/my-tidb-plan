package main

import (
	"github.com/pingcap/tidb/util/mvmap"
	"runtime"
	"strconv"
	"sync"
)



// Join accepts a join query of two relations, and returns the sum of
// relation0.col0 in the final result.
// Input arguments:
//   f0: file name of the given relation0
//   f1: file name of the given relation1
//   offset0: offsets of which columns the given relation0 should be joined
//   offset1: offsets of which columns the given relation1 should be joined
// Output arguments:
//   sum: sum of relation0.col0 in the final result
func Join(f0, f1 string, offset0, offset1 []int) (sum uint64) {
	var(
		tbl0      [][]string
		tbl1      [][]string
		hashtable *mvmap.MVMap
		wg        sync.WaitGroup
	)
	wg.Add(2)
	go func() {
		tbl0=readCSVFileIntoTbl(f0)
		hashtable= buildHashTable(tbl0, offset0)
		wg.Done()
	}()
	go func() {
		tbl1=readCSVFileIntoTbl(f1)
		wg.Done()
	}()
	wg.Wait()
	task:=make(chan struct{},runtime.NumCPU())
	for _, row := range tbl1 {
		go func(r []string) {
			task<- struct{}{}
			rowIDs := probe(hashtable, row, offset1)
			for _, id := range rowIDs {
				v, err := strconv.ParseUint(tbl0[id][0], 10, 64)
				if err != nil {
					panic("JoinExample panic\n" + err.Error())
				}
				sum += v
			}
		}(row)
		<-task
	}
	return sum
}
