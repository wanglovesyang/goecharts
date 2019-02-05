package goecharts

import "encoding/json"

type CurveSettings struct {
	Title         string `json:"title"`
	Smooth        bool   `json:"smooth"`
	TuncPrecision int32  `json:"trunc_precision"`
}

func parseCurveSettings(s interface{}) (ret *CurveSettings, reterr error) {
	switch a := s.(type) {
	case CurveSettings:
		ret = &a
	case *CurveSettings:
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

func Curve(x interface{}, y interface{}, param interface{}) (ret *Chart) {
	var reterr error
	defer func() {
		if reterr != nil {
			panic(reterr)
		}
	}()

	bp, reterr := parseCurveSettings(param)
	if reterr != nil {
		return
	}

	xAxisData, reterr := extractXAxisData(x, y)
	if reterr != nil {
		return
	}

	curveMaker := DefaultSeries
	if bp.Smooth {
		curveMaker = SmoothedSeries
	}
	if bp.TuncPrecision > 0 {
		curveMaker = TruncatedSeriesMaker(curveMaker, bp.TuncPrecision)
	}

	series, reterr := extractSeries(x, y, curveMaker, "line")
	if reterr != nil {
		return
	}

	if xAxisData == nil {
		xAxisData = rangeNum(int32(len(series[0].Data.([]float32))))
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
