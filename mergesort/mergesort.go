package main

import "sync"

// MergeSort performs the merge sort algorithm.
// Please supplement this function to accomplish the home work.
const max = 1 << 11

func merge(s []int64, middle int) {
	helper := make([]int64, len(s))
	copy(helper, s)

	helperLeft := 0
	helperRight := middle
	current := 0
	high := len(s) - 1

	for helperLeft <= middle-1 && helperRight <= high {
		if helper[helperLeft] <= helper[helperRight] {
			s[current] = helper[helperLeft]
			helperLeft++
		} else {
			s[current] = helper[helperRight]
			helperRight++
		}
		current++
	}

	for helperLeft <= middle-1 {
		s[current] = helper[helperLeft]
		current++
		helperLeft++
	}
}



func mergesort(s []int64) {
	if len(s) > 1 {
		middle := len(s) / 2
		mergesort(s[:middle])
		mergesort(s[middle:])
		merge(s, middle)
	}
}



func Mergesortv2(s []int64) {
	len := len(s)
	if len > 1 {
		if len <= max { // Sequential
			mergesort(s)
		} else { // Parallel
			middle := len / 2
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				Mergesortv2(s[:middle])
			}()
			go func() {
				defer wg.Done()
				Mergesortv2(s[middle:])
			}()
			wg.Wait()
			merge(s, middle)
		}
	}
}



func Mergesortv3(s []int64) {
	len := len(s)

	if len > 1 {
		if len <= max {
			mergesort(s)
		} else {
			middle := len / 2
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				Mergesortv3(s[:middle])
			}()
			Mergesortv3(s[middle:])
			wg.Wait()
			merge(s, middle)
		}
	}
}



func Mergesortv1(s []int64) {
	len := len(s)

	if len > 1 {
		middle := len / 2
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			Mergesortv1(s[:middle])
		}()
		go func() {
			defer wg.Done()
			Mergesortv1(s[middle:])
		}()
		wg.Wait()
		merge(s, middle)
	}
}