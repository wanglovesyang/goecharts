package goecharts

type SeriesMaker func(val interface{}, name string, tp string) *Series

func TruncatedSeriesMaker(maker SeriesMaker, trunc int32) SeriesMaker {
	return func(val interface{}, name string, tp string) *Series {
		if fv, suc := val.([]float32); suc {
			val = TruncFloatList(fv, trunc)
		}

		return maker(val, name, tp)
	}
}

func SeriesMakerWithMarkPoint(maker SeriesMaker, trunc int32) SeriesMaker {
	if trunc == 0 {
		trunc = 6
	}

	return func(val interface{}, name string, tp string) (ret *Series) {
		ret = maker(val, name, tp)
		if dataF, suc := val.([]float32); suc {
			min, max := minMaxFloat32InArray(dataF)
			ret.MarkPoint = &MarkPointModes{
				Data: []MarkPoint{
					MarkPoint{
						Name:       "Maximum",
						Symbol:     "pin",
						SymbolSize: 50,
						Label: &SeriesLabelModes{
							Normal: &SeriesLabel{
								Show: true,
								TextStyle: TextStyle{
									Color: "#fff",
								},
							},
						},
						Value: TruncFloat{dataF[max], trunc},
						XAxis: max,
						YAxis: dataF[max],
					},
					MarkPoint{
						Name:       "Minimum",
						Symbol:     "pin",
						SymbolSize: 50,
						Label: &SeriesLabelModes{
							Normal: &SeriesLabel{
								Show: true,
								TextStyle: TextStyle{
									Color: "#fff",
								},
							},
						},
						Value: TruncFloat{dataF[min], trunc},
						XAxis: min,
						YAxis: dataF[min],
					},
				},
			}
		}

		return
	}
}

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
