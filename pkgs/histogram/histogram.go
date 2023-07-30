package histogram

import (
	"sort"
)

var (
	percentiles = []int{25, 50, 75, 100}
)

type Histogram struct {
	Entries []*Entry `json:"entries"`
}

type Entry struct {
	Bin   int `json:"bin"`
	Count int `json:"count"`
}

func Percentile(data []float64) *Histogram {
	sort.Float64s(data)

	percentileMax := make(map[int]float64)
	for _, p := range percentiles {
		idx := (len(data) - 1) * p / 100
		percentileMax[p] = data[idx]
	}

	rs := make(map[int]int)
	for _, fare := range data {
		for _, p := range percentiles {
			if fare <= percentileMax[p] {
				rs[p]++
				break
			}
		}
	}

	var histogram Histogram
	for k, v := range rs {
		histogram.Entries = append(histogram.Entries, &Entry{Bin: k, Count: v})
	}

	return &histogram
}
