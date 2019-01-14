package goecharts

import "fmt"

func Float32ArrToStrArr(s []float32) (ret []string) {
	ret = make([]string, len(s))
	for i, ss := range s {
		ret[i] = fmt.Sprintf("%f", ss)
	}
	return
}
