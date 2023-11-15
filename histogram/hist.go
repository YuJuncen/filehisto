package histogram

import (
	"log"
	"math"
)

// XYer wraps the Len and XY methods.
type XYer interface {
	// Len returns the number of x, y pairs.
	Len() int

	// XY returns an x, y pair.
	XY(int) (x, y float64)
}

// Histogram holds a count of values partionned over buckets.
type Histogram struct {
	// Min is the size of the smallest bucket.
	Min int
	// Max is the size of the biggest bucket.
	Max int
	// Count is the total size of all buckets.
	Count int
	// Total Y value of this histogram.
	FreqCount float64
	// Buckets over which values are partionned.
	Buckets []Bucket
}

// Bucket counts a partion of values.
type Bucket struct {
	// Count is the number of values represented in the bucket.
	Count int
	// Min is the low, inclusive bound of the bucket.
	Min float64
	// Max is the high, exclusive bound of the bucket. If
	// this bucket is the last bucket, the bound is inclusive
	// and contains the max value of the histogram.
	Max float64
}

// Hist creates an histogram partionning input over `bins` buckets.
func Hist(bins int, input XYer) Histogram {
	if input.Len() == 0 || bins == 0 {
		return Histogram{}
	}

	min, _ := input.XY(0)
	max := min
	for i := 0; i < input.Len(); i++ {
		val, _ := input.XY(i)
		min = math.Min(min, val)
		max = math.Max(max, val)
	}

	if min == max {
		return Histogram{
			Min:     input.Len(),
			Max:     input.Len(),
			Count:   input.Len(),
			Buckets: []Bucket{{Count: input.Len(), Min: min, Max: max}},
		}
	}

	scale := (max - min) / float64(bins)
	buckets := make([]Bucket, bins)
	for i := range buckets {
		bmin, bmax := float64(i)*scale+min, float64(i+1)*scale+min
		buckets[i] = Bucket{Min: bmin, Max: bmax}
	}

	minC, maxC := 0, 0
	totalC := float64(0)
	for i := 0; i < input.Len(); i++ {
		val, cnt := input.XY(i)
		minx := float64(min)
		xdiff := val - minx
		bi := imin(int(xdiff/scale), len(buckets)-1)
		if bi < 0 || bi >= len(buckets) {
			log.Panicf("bi=%d\tval=%f\txdiff=%f\tscale=%f\tlen(buckets)=%d", bi, val, xdiff, scale, len(buckets))
		}
		buckets[bi].Count += int(cnt)
		totalC += cnt
		minC = imin(minC, buckets[bi].Count)
		maxC = imax(maxC, buckets[bi].Count)
	}

	return Histogram{
		Min:       minC,
		Max:       maxC,
		Count:     input.Len(),
		FreqCount: totalC,
		Buckets:   buckets,
	}
}

// PowerHist creates an histogram partionning input over buckets of power
// `pow`.
func PowerHist(power float64, input []float64) Histogram {
	if len(input) == 0 || power <= 0 {
		return Histogram{}
	}

	minx, maxx := input[0], input[0]
	for _, val := range input {
		minx = math.Min(minx, val)
		maxx = math.Max(maxx, val)
	}

	fromPower := math.Floor(logbase(minx, power))
	toPower := math.Floor(logbase(maxx, power))

	buckets := make([]Bucket, int(toPower-fromPower)+1)
	for i, bkt := range buckets {
		bkt.Min = math.Pow(power, float64(i)+fromPower)
		bkt.Max = math.Pow(power, float64(i+1)+fromPower)
		buckets[i] = bkt
	}

	minC := 0
	maxC := 0
	for _, val := range input {
		powAway := logbase(val, power) - fromPower
		bi := int(math.Floor(powAway))
		buckets[bi].Count++
		minC = imin(buckets[bi].Count, minC)
		maxC = imax(buckets[bi].Count, maxC)
	}

	return Histogram{
		Min:     minC,
		Max:     maxC,
		Count:   len(input),
		Buckets: buckets,
	}
}

// Scale gives the scaled count of the bucket at idx, using the
// provided scale func.
func (h Histogram) Scale(s ScaleFunc, idx int) float64 {
	bkt := h.Buckets[idx]
	scale := s(h.Min, h.Max, bkt.Count)
	return scale
}

// ScaleFunc is the type to implement to scale an histogram.
type ScaleFunc func(min, max, value int) float64

// Linear builds a ScaleFunc that will linearly scale the values of
// an histogram so that they do not exceed width.
func Linear(width int) ScaleFunc {
	return func(min, max, value int) float64 {
		if min == max {
			return 1
		}
		return float64(value-min) / float64(max-min) * float64(width)
	}
}

func logbase(a, base float64) float64 {
	return math.Log2(a) / math.Log2(base)
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
