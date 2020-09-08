package false_sharing

import (
	"math/rand"
	"sync"
	"testing"
)

const (
	limitOfMetrics  = 5000
	limitOfParallel = 64
)

var (
	o       sync.Once
	metrics [][]int
)

func initMetrics() {
	// to guarantee the same metrics
	rand.NewSource(314159265359)
	// gen the metrics
	for i := 0; i < limitOfMetrics; i++ {
		rows := make([]int, 0, limitOfMetrics)
		for j := 0; j < limitOfMetrics; j++ {
			rows = append(rows, rand.Int())
		}
		metrics = append(metrics, rows)
	}
}

func BenchmarkRowTraversal(b *testing.B) {
	o.Do(initMetrics)
	numRow := len(metrics)
	numCol := len(metrics[0])
	b.ResetTimer()
	for i := 0; i < numRow; i++ {
		for j := 0; j < numCol; j++ {
			gg := metrics[i][j]
			gg++
		}
	}
}

func BenchmarkColTraversal(b *testing.B) {
	o.Do(initMetrics)
	numRow := len(metrics)
	numCol := len(metrics[0])
	b.ResetTimer()

	for j := 0; j < numCol; j++ {
		for i := 0; i < numRow; i++ {
			gg := metrics[i][j]
			gg++
		}
	}
}

func BenchmarkOddTraversal(b *testing.B) {
	var wg sync.WaitGroup
	o.Do(initMetrics)
	ret := make([]int, limitOfParallel)

	b.ResetTimer()
	for tID := 0; tID < limitOfParallel; tID++ {
		wg.Add(1)
		go func(threadNo int) {
			defer wg.Done()
			for i, row := range metrics {
				if i%limitOfParallel != 0 {
					continue
				}
				for _, v := range row {
					if v%2 != 0 {
						ret[threadNo]++
					}
				}
			}
		}(tID)
	}
	wg.Wait()

	numOdds := 0
	for _, r := range ret {
		numOdds += r
	}
}

func BenchmarkOddTraversalOpt(b *testing.B) {
	var wg sync.WaitGroup
	o.Do(initMetrics)
	ret := make([]int, limitOfParallel)

	b.ResetTimer()
	for tID := 0; tID < limitOfParallel; tID++ {
		wg.Add(1)
		go func(threadNo int) {
			defer wg.Done()
			for i, row := range metrics {
				if i%limitOfParallel != 0 {
					continue
				}
				cnt := 0
				for _, v := range row {
					if v%2 != 0 {
						cnt++
					}
				}
				ret[threadNo] += cnt
			}
		}(tID)
	}
	wg.Wait()

	numOdds := 0
	for _, r := range ret {
		numOdds += r
	}
}

func BenchmarkOddTraversalOldSchool(b *testing.B) {
	o.Do(initMetrics)
	numRow := len(metrics)
	numCol := len(metrics[0])

	var wg sync.WaitGroup
	ret := make([]int, limitOfParallel)
	b.ResetTimer()

	for tID := 0; tID < limitOfParallel; tID++ {
		wg.Add(1)
		go func(threadNo int) {
			defer wg.Done()
			for i := 0; i < numRow; i++ {
				if i%limitOfParallel != 0 {
					continue
				}
				for j := 0; j < numCol; j++ {
					if metrics[i][j]%2 != 0 {
						ret[threadNo]++
					}
				}
			}
		}(tID)
	}
	wg.Wait()

	numOdds := 0
	for _, r := range ret {
		numOdds += r
	}
}

func BenchmarkOddTraversalOldSchoolOpt(b *testing.B) {
	o.Do(initMetrics)
	numRow := len(metrics)
	numCol := len(metrics[0])

	var wg sync.WaitGroup
	ret := make([]int, limitOfParallel)
	b.ResetTimer()

	for tID := 0; tID < limitOfParallel; tID++ {
		wg.Add(1)
		go func(threadNo int) {
			defer wg.Done()
			for i := 0; i < numRow; i++ {
				if i%limitOfParallel != 0 {
					continue
				}
				cnt := 0
				for j := 0; j < numCol; j++ {
					if metrics[i][j]%2 != 0 {
						cnt++
					}
				}
				if cnt != 0 {
					ret[threadNo] += cnt
				}
			}
		}(tID)
	}
	wg.Wait()

	numOdds := 0
	for _, r := range ret {
		numOdds += r
	}
}
