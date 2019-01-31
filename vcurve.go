package goecharts

import (
	"fmt"
	"strings"

	gf "github.com/wanglovesyang/gframe"
)

func pairValues(x interface{}, y interface{}) (ret [][2][]float32, names []string, reterr error) {
	switch xx := x.(type) {
	case [][]float32:
		if yy, suc := y.([][]float32); suc {
			if len(yy) != len(xx) {
				reterr = fmt.Errorf("when using parallel arrays, parameter x and y shall have the same dimensions")
				return
			}

			ret = make([][2][]float32, len(yy))
			names = make([]string, len(yy))
			for i := 0; i < len(yy); i++ {
				ret[i][0] = xx[i]
				ret[i][1] = yy[i]
				names[i] = fmt.Sprintf("s-%d", i)
			}
		} else {
			reterr = fmt.Errorf("when using parallel arrays, x and y shall all be float32 array (list)")
		}
	case []float32:
		if yy, suc := y.([]float32); suc {
			ret = make([][2][]float32, 1)
			names = make([]string, 1)
			ret[0][0] = xx
			ret[0][1] = yy
			names[0] = "s-0"
		} else {
			reterr = fmt.Errorf("when using parallel arrays, x and y shall all be float32 array (list)")
		}
	case []string:
		yy, suc := y.(*gf.DataFrame)
		if !suc {
			reterr = fmt.Errorf("y should be dataframe when string-arrayed x is provided")
			return
		}

		fldPairs, err := detectFieldPairs(xx)
		if err != nil {
			reterr = err
			return
		}

		names = make([]string, len(fldPairs))
		ret = make([][2][]float32, len(fldPairs))
		for i, p := range fldPairs {
			cols, err := yy.GetValColumns(p[0:2]...)
			if err != nil {
				return
			}

			ret[i] = [2][]float32{cols[0], cols[1]}
			names[i] = fmt.Sprintf("%s on %s", p[1], p[0])
		}
	}

	return
}

func detectFieldPairs(x []string) (ret [][2]string, reterr error) {
	cnt := 0
	for _, f := range x {
		if strings.Contains(f, ",") {
			cnt++
		}
	}

	if cnt == len(x) {
		ret = make([][2]string, 0, len(x))
		for _, f := range x {
			eles := strings.Split(f, ",")
			ret = append(ret, [2]string{eles[0], eles[1]})
		}
	} else {
		if len(x)%2 != 0 {
			reterr = fmt.Errorf("when using field pairs, please use comma splitted field spec or give a array which can be mod with 2")
			return
		}

		ret = make([][2]string, 0, len(x)/2)
		for i := 0; i < len(x); i += 2 {
			ret = append(ret, [2]string{x[i], x[i+1]})
		}
	}

	return
}

func trainsposePairs(s [2][]float32) (ret [][2]float32, reterr error) {
	l := len(s[0])
	if len(s[1]) != l {
		reterr = fmt.Errorf("lengths in the pair is inconsistent")
		return
	}

	ret = make([][2]float32, l)
	for i := 0; i < l; i++ {
		ret[i] = [2]float32{s[0][i], s[1][i]}
	}

	return
}

func trainsposePairsList(s [][2][]float32) (ret [][][2]float32, reterr error) {
	ret = make([][][2]float32, len(s))
	for i, ss := range s {
		if ret[i], reterr = trainsposePairs(ss); reterr != nil {
			return
		}
	}

	return
}

func VCurve(x interface{}, y interface{}, param interface{}) (ret *Chart) {
	var reterr error
	defer func() {
		if reterr != nil {
			panic(reterr)
		}
	}()

	data, names, err := pairValues(x, y)
	if err != nil {
		reterr = err
		return
	}

	dataPairs, reterr := trainsposePairsList(data)
	if reterr != nil {
		return
	}

	bp, reterr := parseCurveSettings(param)
	if reterr != nil {
		return
	}

	xAxis := DefaultXAxis(nil, "value")
	title := DefaultTitle(bp.Title)

	curveMaker := DefaultSeries
	if bp.Smooth {
		curveMaker = SmoothedSeries
	}

	serieses := make([]*Series, 0, len(dataPairs))
	for i, pairs := range dataPairs {
		series := curveMaker(pairs, names[i], "line")
		serieses = append(serieses, series)
	}

	ret = &Chart{
		opt: &ChartOption{
			Title:   []*Title{title},
			ToolBox: DefaultToolBox(),
			ToolTip: DefaultToolTip(),
			Series:  serieses,
			Legend: []*Legend{
				DefaultLegend(names),
			},
			XAxis: []*XAxis{
				xAxis,
			},
			YAxis: []*YAxis{
				DefaultYAxis(),
			},
			Color: defaultColorSet,
		},
	}

	return
}
