package goecharts

import (
	"fmt"
	"math"
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

func minMaxFloat32InArray(s []float32) (min, max int32) {
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

func minFloat32(a, b float32) (ret float32) {
	if a > b {
		return b
	}

	return a
}

func maxFloat32(a, b float32) (ret float32) {
	if a < b {
		return b
	}

	return a
}

type TruncFloat struct {
	v     float32
	trunc int32
}

func (t *TruncFloat) MarshalJSON() ([]byte, error) {
	fmtS := fmt.Sprintf("%%.%df", t.trunc)
	return []byte(fmt.Sprintf(fmtS, t.v)), nil
}

func TruncFloatList(s []float32, trunc int32) (ret []TruncFloat) {
	ret = make([]TruncFloat, len(s))
	for i, ss := range s {
		ret[i] = TruncFloat{
			v:     ss,
			trunc: trunc,
		}
	}

	return
}

func histgoram(data []float32, binSize int32, norm bool) (ret []float32, axis []float32) {
	ret = make([]float32, binSize)
	axis = make([]float32, binSize+1)
	min, max := minMaxFloat32InArray(data)
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

func alignedHistogram(data map[string][]float32, binSize int32, normalize bool) (ret map[string][]float32, axis []float32) {
	var min, max float32 = math.MaxFloat32, -math.MaxFloat32
	for _, l := range data {
		minID, maxID := minMaxFloat32InArray(l)
		minV, maxV := l[minID], l[maxID]
		min = minFloat32(min, minV)
		max = maxFloat32(max, maxV)
	}

	binStride := (max - min) / float32(binSize)

	getId := func(v float32) int32 {
		id := int32((v - min) / float32(binStride))
		if id < 0 {
			id = 0
		} else if id >= binSize {
			id = binSize - 1
		}

		return id
	}

	ret = make(map[string][]float32)
	for k, s := range data {
		rv := make([]float32, binSize)
		for _, v := range s {
			rv[getId(v)]++
		}
		ret[k] = rv
	}

	if normalize {
		for k, rv := range ret {
			norm := float32(len(data[k]))
			for i, v := range rv {
				rv[i] = v / norm
			}
		}
	}

	axis = make([]float32, binSize+1)
	for i := int32(0); i < binSize; i++ {
		axis[i] = min + float32(i)*binStride
	}
	axis[binSize] = max

	return
}
