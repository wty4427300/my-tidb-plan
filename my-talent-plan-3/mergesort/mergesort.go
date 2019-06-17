package main

import "sync"

// MergeSort performs the merge sort algorithm.
// Please supplement this function to accomplish the home work.


//Three optimized versions are available(Mergesortv1, Mergesortv2, Mergesortv3)


//Use mergesort when less than Max
const max = 1 << 11
//merge phase
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


//Sequential mergesort
func mergesort(s []int64) {
	if len(s) > 1 {
		// Sequential
		middle := len(s) / 2
		mergesort(s[:middle])
		mergesort(s[middle:])
		merge(s, middle)
	}
}


//Sequential and Parallel mergesort
func Mergesortv2(s []int64) {
	len := len(s)
	if len > 1 {
		if len <= max {
			// Sequential
			mergesort(s)
		} else {
			// Parallel
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


//Sequential and Parallel mergesort
func Mergesortv3(s []int64) {
	len := len(s)

	if len > 1 {
		if len <= max {
			// Sequential
			mergesort(s)
		} else {
			// Parallel
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


//Parallel mergesort
func Mergesortv1(s []int64) {
	len := len(s)
	if len > 1 {
		// Parallel
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