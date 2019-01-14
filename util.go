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
