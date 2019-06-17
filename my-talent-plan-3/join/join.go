package main

import (
	"github.com/pingcap/tidb/util/mvmap"
	"runtime"
	"strconv"
	"sync"
	"unsafe"
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
	rowIDs := make(chan []int64,runtime.NumCPU())
	for _, row := range tbl1 {
		go myProbe(hashtable, row, offset1, rowIDs)
	}
	for i := 0; i < len(tbl1); i++ {
		rowID := <-rowIDs
		for _, id := range rowID {
			v, err := strconv.ParseUint(tbl0[id][0], 10, 64)
			if err != nil {
				panic("JoinExample panic\n" + err.Error())
			}
			sum += v
		}
	}
	return sum
}

func myProbe(hashtable *mvmap.MVMap, row []string, offset []int, rowIDs chan<- []int64) {
	var keyHash []byte
	var vals [][]byte
	var rowID []int64
	for _, off := range offset {
		keyHash = append(keyHash, []byte(row[off])...)
	}
	vals = hashtable.Get(keyHash, vals)
	for _, val := range vals {
		rowID = append(rowID, *(*int64)(unsafe.Pointer(&val[0])))
	}
	rowIDs <- rowID;
}