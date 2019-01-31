package goecharts

import (
	"fmt"
	"os"
)

func Float32ArrToStrArr(s []float32) (ret []string) {
	ret = make([]string, len(s))
	for i, ss := range s {
		ret[i] = fmt.Sprintf("%f", ss)
	}
	return
}

func Log(s string) {
	fmt.Fprintln(os.Stderr, s)
}

func minMaxFloat32(s []float32) (min, max int32) {
	for i, v := range s {
		if v < s[min] {
			min = int32(i)
		}

		if v > s[max] {
			max = int32(i)
		}
	}

	return
}

type TruncFloat struct {
	v     float32
	trunc int32
}

func (t *TruncFloat) MarshalJSON() ([]byte, error) {
	fmtS := fmt.Sprintf("%%.%df", t.trunc)
	return []byte(fmt.Sprintf(fmtS, t.v)), nil
}

func histgoram(data []float32, binSize int32, norm bool) (ret []float32, axis []float32) {
	ret = make([]float32, binSize)
	axis = make([]float32, binSize+1)
	min, max := minMaxFloat32(data)
	minV, maxV := data[min], data[max]
	binStride := (maxV - minV) / float32(binSize)

	getId := func(v float32) int32 {
		id := int32((v - minV) / float32(binStride))
		if id < 0 {
			id = 0
		} else if id >= binSize {
			id = binSize - 1
		}

		return id
	}

	for _, v := range data {
		ret[getId(v)]++
	}

	for i, v := range ret {
		ret[i] = v / float32(len(data))
		axis[i] = minV + float32(i)*binStride
	}
	axis[binSize] = maxV

	return
}

func alignedHistogram(data map[string][]float32, binSize int32) (ret map[string][]float32, axis []float32) {
	var mins []float32
	var maxs []float32
	for k, l := range data {
		min, max := minMaxFloat32(l)
		minV, maxV := l[min], l[max]
		mins = append(mins, minV)
		maxs = append(maxs, maxV)

	}

}
