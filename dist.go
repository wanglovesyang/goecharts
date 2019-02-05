package goecharts

import (
	"encoding/json"
	"fmt"

	gf "github.com/wanglovesyang/gframe"
)

type HistogramSettings struct {
	BinSize        int32  `json:"bin_size"`
	Normalize      bool   `json:"normalize"`
	Title          string `json:"title"`
	TruncPrecision int32  `json:"trunc_precision"`
}

func parseHistogramSettings(s interface{}) (ret *HistogramSettings, reterr error) {
	switch a := s.(type) {
	case HistogramSettings:
		ret = &a
	case *HistogramSettings:
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

func renderHistogram(data map[string][]float32, settings *HistogramSettings) (ret *Chart) {
	var displaySeriesNames []string

	opt := &ChartOption{
		Title:     []*Title{DefaultTitle(settings.Title)},
		ToolBox:   DefaultToolBox(),
		ToolTip:   DefaultToolTip(),
		Animation: true,
		YAxis: []*YAxis{
			DefaultYAxis(),
		},
		Color: defaultColorSet,
	}

	maker := DefaultSeries
	if settings.TruncPrecision > 0 {
		maker = TruncatedSeriesMaker(maker, settings.TruncPrecision)
	}

	hist, axis := alignedHistogram(data, settings.BinSize, settings.Normalize)
	opt.XAxis = []*XAxis{DefaultXAxis(axis, "category")}
	for k, h := range hist {
		series := maker(h, k, "bar")
		opt.Series = append(opt.Series, series)
		displaySeriesNames = append(displaySeriesNames, k)
	}

	opt.Legend = []*Legend{
		DefaultLegend(displaySeriesNames),
	}

	ret = &Chart{
		opt: opt,
	}
	return
}

func Histogram(data interface{}, param interface{}) (ret *Chart) {
	bp, err := parseHistogramSettings(param)
	if err != nil {
		panic(err)
	}

	mtx := make(map[string][]float32)
	switch d := data.(type) {
	case []float32:
		mtx["-"] = d
	case map[string][]float32:
		mtx = d
	case [][]float32:
		for i, l := range d {
			mtx[fmt.Sprintf("s%d", i)] = l
		}
	case *gf.DataFrame:
		cols := d.ValueColumnNames()
		colVals, _ := d.GetValColumns(cols...)
		for i := 0; i < len(cols); i++ {
			mtx[cols[i]] = colVals[i]
		}
	}

	ret = renderHistogram(mtx, bp)
	return
}
