package goecharts

import (
	"encoding/json"
	"fmt"
	"reflect"

	gf "github.com/wanglovesyang/gframe"
)

type BarSettings struct {
	Title          string `json:"title"`
	TruncPrecision int32  `json:"trunc_precision"`
}

func parseBarSettings(s interface{}) (ret *BarSettings, reterr error) {
	switch a := s.(type) {
	case BarSettings:
		ret = &a
	case *BarSettings:
		ret = a
	default:
		jsonData, err := json.Marshal(s)
		if err != nil {
			reterr = err
			return
		}
		if err := json.Unmarshal(jsonData, &ret); err != nil {
			reterr = err
			return
		}
	}

	return
}

func rangeNum(size int32) (ret []int32) {
	ret = make([]int32, size)
	for i := int32(0); i < size; i++ {
		ret[i] = int32(i)
	}
	return
}

func extractXAxisData(x, y interface{}) (ret interface{}, reterr error) {
	if x == nil {
		return nil, nil
	}

	switch xx := x.(type) {
	case string:
		if dy, suc := y.(*gf.DataFrame); !suc {
			reterr = fmt.Errorf("argument y should be a dataframe when x is given as a string")
		} else {
			if valy, err := dy.GetIdColumns(xx); err == nil {
				ret = valy[0]
				return
			}

			if valy, err := dy.GetValColumns(xx); err == nil {
				ret = valy[0]
				return
			}

			reterr = fmt.Errorf("columns %s does not exist in dataframe", xx)
		}
	case *gf.DataFrame:
		if xx.Shape()[1] != 1 {
			reterr = fmt.Errorf("cannot adopt data frame with multiple columns as x label")
			return
		}

		if col, err := xx.GetIdColumns(xx.Columns()...); err == nil {
			ret = col
			return
		}

		if col, err := xx.GetValColumns(xx.Columns()...); err == nil {
			ret = col
			return
		}

		reterr = fmt.Errorf("are you kidding me?")
	default:
		if tp := reflect.TypeOf(x).Kind(); tp == reflect.Slice {
			ret = x
			return
		}

		reterr = fmt.Errorf("Invalid type of x axis")
	}

	return
}

func extractSeries(x, y interface{}, maker SeriesMaker, tp string) (ret []*Series, reterr error) {
	switch yy := y.(type) {
	case []float32:
		ret = []*Series{maker(yy, "-", tp)}
	case map[string][]float32:
		for k, v := range yy {
			ret = append(ret, maker(v, k, tp))
		}
	case [][]float32:
		for k, v := range yy {
			ret = append(ret, maker(v, fmt.Sprintf("s%d", k), tp))
		}
	case *gf.DataFrame:
		xx := x.(string)

		columns := yy.ValueColumnNames()
		if len(columns) == 0 {
			reterr = fmt.Errorf("no value columns in the dataframe")
			return
		}

		var cols [][]float32
		if cols, reterr = yy.GetValColumns(columns...); reterr != nil {
			return
		}

		for i, c := range cols {
			if columns[i] != xx {
				ret = append(ret, maker(c, columns[i], tp))
			}
		}
	default:
		reterr = fmt.Errorf("unsupported y type")
	}

	return
}

func Bar(x interface{}, y interface{}, param interface{}) (ret *Chart) {
	var reterr error
	defer func() {
		if reterr != nil {
			panic(reterr)
		}
	}()

	bp, reterr := parseBarSettings(param)
	if reterr != nil {
		return
	}

	xAxisData, reterr := extractXAxisData(x, y)
	if reterr != nil {
		return
	}

	maker := DefaultSeries
	if bp.TruncPrecision > 0 {
		maker = TruncatedSeriesMaker(maker, bp.TruncPrecision)
	}

	series, reterr := extractSeries(x, y, DefaultSeries, "bar")
	if reterr != nil {
		return
	}

	xAxis := DefaultXAxis(xAxisData, "category")
	title := DefaultTitle(bp.Title)

	var displaySeriesNames []string
	for _, s := range series {
		displaySeriesNames = append(displaySeriesNames, s.Name)
	}

	ret = &Chart{
		opt: &ChartOption{
			Title:   []*Title{title},
			ToolBox: DefaultToolBox(),
			ToolTip: DefaultToolTip(),
			Series:  series,
			Legend: []*Legend{
				DefaultLegend(displaySeriesNames),
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
