package util

import "strconv"

type Strings []string

func (ss Strings) ToFloat64() []float64 {
	fs := make([]float64, 0, len(ss))
	for _, s := range ss {
		f, _ := strconv.ParseFloat(s, 64)
		fs = append(fs, f)
	}
	return fs
}
