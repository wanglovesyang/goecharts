package goecharts

type SeriesMaker func(val interface{}, name string, tp string) *Series

var defaultColorSet = []string{
	"#c23531",
	"#2f4554",
	"#61a0a8",
	"#d48265",
	"#749f83",
	"#ca8622",
	"#bda29a",
	"#6e7074",
	"#546570",
	"#c4ccd3",
	"#f05b72",
	"#ef5b9c",
	"#f47920",
	"#905a3d",
	"#fab27b",
	"#2a5caa",
	"#444693",
	"#726930",
	"#b2d235",
	"#6d8346",
	"#ac6767",
	"#1d953f",
	"#6950a1",
	"#918597",
	"#f6f5ec",
}

func SmoothedSeries(data interface{}, name string, seriesType string) (ret *Series) {
	ret = DefaultSeries(data, name, seriesType)
	ret.Smooth = true

	return ret
}
