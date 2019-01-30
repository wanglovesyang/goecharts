package goecharts

import (
	"crypto/md5"
	h "encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

var template = `
<script>
require.config({
	paths: {
		'echarts': '/nbextensions/echarts/echarts.min'
	}
});
</script>
<div id="${chart_id}" style="width:800px;height:400px;"></div>

<script>
require(['echarts'], function(echarts) {
	
var myChart_${chart_id} = echarts.init(document.getElementById('${chart_id}'), 'light', {renderer: 'canvas'});

var option_${chart_id} = 
${chart_opt};
myChart_${chart_id}.setOption(option_${chart_id});

});
</script>
`

type Chart struct {
	opt *ChartOption
}

func (e *Chart) RenderJupyter() (ret string) {
	optStr, err := json.Marshal(e.opt)
	if err != nil {
		ret = err.Error()
		return
	}

	chartID := time.Now().Format("2006-01-02 15:04:05")
	sm := md5.Sum([]byte(chartID))
	chartSig := h.EncodeToString(sm[0:16])

	ret = strings.Replace(ret, "${chart_opt}", string(optStr), -1)
	ret = strings.Replace(ret, "${chart_id}", chartSig, -1)

	Log(ret)
	return
}

type ChartOption struct {
	Title     []*Title  `json:"title"`
	ToolBox   *ToolBox  `json:"toolbox"`
	ToolTip   *ToolTip  `json:"tooltip"`
	Series    []*Series `json:"series"`
	Animation bool      `json:"animation"`
	XAxis     []*XAxis  `json:"xAxis"`
	YAxis     []*YAxis  `json:"yAxis"`
	Legend    []*Legend `json:"legend"`
}

/*
"title":{
            "text": "Bar chart",
            "subtext": "precipitation and evaporation one year",
            "left": "auto",
            "top": "auto",
            "textStyle": {
                "fontSize": 18
            },
            "subtextStyle": {
                "fontSize": 12
            }
        }
*/

type TextStyle struct {
	FontSize int32  `json:"frontSize"`
	Color    string `json:"color"`
}

type Title struct {
	Text         string     `json:"text"`
	SubText      string     `json:"subtext"`
	LeftMode     string     `json:"left"`
	TopMode      string     `json:"top"`
	TextStyle    *TextStyle `json:"textStyle"`
	SubTextStyle *TextStyle `json:"subTextStyle"`
}

func DefaultTitle(title string) *Title {
	return &Title{
		Text:     title,
		LeftMode: "auto",
		TopMode:  "auto",
		TextStyle: &TextStyle{
			FontSize: 18,
		},
	}
}

/*
 "toolbox": {
        "show": true,
        "orient": "vertical",
        "left": "95%",
        "top": "center",
        "feature": {
            "saveAsImage": {
                "show": true,
                "title": "save as image"
            },
            "restore": {
                "show": true,
                "title": "restore"
            },
            "dataView": {
                "show": true,
                "title": "data view"
            }
        }
    },
*/

type ToolBoxFeature struct {
	Show  bool   `json:"show"`
	Title string `json:"title"`
}

type ToolBox struct {
	Show        bool                      `json:"show"`
	Orientation string                    `json:"orient"`
	Left        string                    `json:"left"`
	Top         string                    `json:"top"`
	Feature     map[string]ToolBoxFeature `json:"feature"`
}

type AxisPointer struct {
	Type string `json:"type"`
}

func DefaultToolBox() *ToolBox {
	return &ToolBox{
		Show:        true,
		Orientation: "vertical",
		Left:        "95%",
		Top:         "center",
		Feature: map[string]ToolBoxFeature{
			"saveAsImage": ToolBoxFeature{
				Show:  true,
				Title: "save as image",
			},
			"restore": ToolBoxFeature{
				Show:  true,
				Title: "restore",
			},
			"dataView": ToolBoxFeature{
				Show:  true,
				Title: "data view",
			},
		},
	}
}

/*
"tooltip": {
        "trigger": "item",
        "triggerOn": "mousemove|click",
        "axisPointer": {
            "type": "line"
        },
        "textStyle": {
            "fontSize": 14
        },
        "backgroundColor": "rgba(50,50,50,0.7)",
        "borderColor": "#333",
        "borderWidth": 0
    },
*/

type ToolTip struct {
	Trigger         string      `json:"trigger"`
	TriggerOn       string      `json:"triggerOn"`
	AxisPointer     AxisPointer `json:"axisPointer"`
	TextStyle       TextStyle   `json:"textStyle"`
	BackGroundColor string      `json:"backGroundColor"`
	BorderColor     string      `json:"borderColor"`
	BorderWidth     int32       `json:"borderWidth"`
}

func DefaultToolTip() *ToolTip {
	return &ToolTip{
		Trigger:   "item",
		TriggerOn: "mousemove|click",
		AxisPointer: AxisPointer{
			Type: "line",
		},
		TextStyle: TextStyle{
			FontSize: 14,
		},
		BackGroundColor: "rgba(50,50,50,0.7)",
		BorderColor:     "#333",
		BorderWidth:     0,
	}
}

/*
{
            "type": "bar",
            "name": "precipitation",
            "data": [
                2.0,
                4.9,
                7.0,
                23.2,
                25.6,
                76.7,
                135.6,
                162.2,
                32.6,
                20.0,
                6.4,
                3.3
            ],
            "barCategoryGap": "20%",
            "label": {
                "normal": {
                    "show": false,
                    "position": "top",
                    "textStyle": {
                        "fontSize": 12
                    }
                },
                "emphasis": {
                    "show": true,
                    "textStyle": {
                        "fontSize": 12
                    }
                }
            },
            "markPoint": {
                "data": [
                    {
                        "type": "max",
                        "name": "Maximum",
                        "symbol": "pin",
                        "symbolSize": 50,
                        "label": {
                            "normal": {
                                "textStyle": {
                                    "color": "#fff"
                                }
                            }
                        }
                    },
                    {
                        "type": "min",
                        "name": "Minimum",
                        "symbol": "pin",
                        "symbolSize": 50,
                        "label": {
                            "normal": {
                                "textStyle": {
                                    "color": "#fff"
                                }
                            }
                        }
                    }
                ]
            },
            "markLine": {
                "data": [
                    {
                        "type": "average",
                        "name": "mean-Value"
                    }
                ],
                "symbolSize": 10
            },
            "seriesId": 5238410
        },
*/

type SeriesLabelModes struct {
	Normal   *SeriesLabel `json:"normal"`
	Emphasis *SeriesLabel `json:"emphasis"`
}

type SeriesLabel struct {
	Show      bool      `json:"show"`
	Position  string    `json:"top"`
	TextStyle TextStyle `json:"textStyle"`
}

type Series struct {
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Data        []float32         `json:"data"`
	CategoryGap string            `json:"barCategoryGap"`
	Label       *SeriesLabelModes `json:"label"`
	MarkPoint   *MarkPointModes   `json:"markPoint"`
	MarkLine    *MarkLineModes    `json:"markLine"`
}

type MarkPoint struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	Symbol     string            `json:"symbol"`
	SymbolSize int32             `json:"symbolSize"`
	Label      *SeriesLabelModes `json:"label"`
}

type MarkPointModes struct {
	Data []MarkPoint `json:"data"`
}

type MarkLine struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type MarkLineModes struct {
	Data       []MarkLine `json:"data"`
	SymbolSize int32      `json:"symbolSize"`
}

func DefaultSeries(data []float32, name string, seriesType string) *Series {
	return &Series{
		Type:        seriesType,
		Name:        name,
		Data:        data,
		CategoryGap: "20%",
		Label: &SeriesLabelModes{
			Normal: &SeriesLabel{
				Show:     true,
				Position: "top",
				TextStyle: TextStyle{
					FontSize: 12,
				},
			},
			Emphasis: &SeriesLabel{
				Show:     true,
				Position: "top",
				TextStyle: TextStyle{
					FontSize: 12,
				},
			},
		},
		MarkPoint: &MarkPointModes{
			Data: []MarkPoint{
				MarkPoint{
					Type:       "max",
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
				},
				MarkPoint{
					Type:       "min",
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
				},
			},
		},
		MarkLine: &MarkLineModes{
			Data: []MarkLine{
				MarkLine{
					Type: "average",
					Name: "mean-value",
				},
				MarkLine{},
			},
			SymbolSize: 10,
		},
	}
}

/*
"xAxis": [
        {
            "show": true,
            "nameLocation": "middle",
            "nameGap": 25,
            "nameTextStyle": {
                "fontSize": 14
            },
            "axisTick": {
                "alignWithLabel": false
            },
            "inverse": false,
            "boundaryGap": true,
            "type": "category",
            "splitLine": {
                "show": false
            },
            "axisLine": {
                "lineStyle": {
                    "width": 1
                }
            },
            "axisLabel": {
                "interval": "auto",
                "rotate": 0,
                "margin": 8,
                "textStyle": {
                    "fontSize": 12
                }
            },
            "data": [
                "Jan",
                "Feb",
                "Mar",
                "Apr",
                "May",
                "Jun",
                "Jul",
                "Aug",
                "Sep",
                "Oct",
                "Nov",
                "Dec"
            ]
        }
    ],
*/

type XAxis struct {
	Show          bool        `json:"show"`
	NameLoc       string      `json:"nameLocation"`
	NameGap       float32     `json:"nameGap"`
	NameTextStyle *TextStyle  `json:"nameTextStyle"`
	Inverse       bool        `json:"inverse"`
	BoundaryGap   bool        `json:"boundaryGap"`
	Type          string      `json:"type"`
	SplitLine     *SplitLine  `json:"splitLine"`
	AxisLine      *AxisLine   `json:"axisLine"`
	AxisLabel     *AxisLabel  `json:"axisLabel"`
	Data          interface{} `json:"data"`
}

type SplitLine struct {
	Show bool
}

type AxisTick struct {
	AlignWithLabel bool `json:"aligWithLabel"`
}

type AxisLine struct {
	LineStyle *LineStyle `json:"lineStyle"`
}

type LineStyle struct {
	Width float32 `json:"lineStyle"`
}

type AxisLabel struct {
	Interval  string     `json:"interval"`
	Formater  string     `json:"formater"`
	Rotate    float32    `json:"rotate"`
	Margin    float32    `json:"margin"`
	TextStyle *TextStyle `json:"textSyle"`
}

func DefaultXAxis(data interface{}, tp string) *XAxis {
	return &XAxis{
		Data:    data,
		Show:    true,
		NameLoc: "middle",
		NameGap: 25,
		NameTextStyle: &TextStyle{
			FontSize: 14,
		},
		Inverse:     false,
		BoundaryGap: true,
		Type:        tp,
		SplitLine: &SplitLine{
			Show: true,
		},
		AxisLine: &AxisLine{
			LineStyle: &LineStyle{
				Width: 1,
			},
		},
		AxisLabel: &AxisLabel{
			Interval: "auto",
			Rotate:   0,
			Margin:   8,
			TextStyle: &TextStyle{
				FontSize: 12,
			},
		},
	}
}

/*
"yAxis": [
        {
            "show": true,
            "nameLocation": "middle",
            "nameGap": 25,
            "nameTextStyle": {
                "fontSize": 14
            },
            "axisTick": {
                "alignWithLabel": false
            },
            "inverse": false,
            "boundaryGap": true,
            "type": "value",
            "splitLine": {
                "show": true
            },
            "axisLine": {
                "lineStyle": {
                    "width": 1
                }
            },
            "axisLabel": {
                "interval": "auto",
                "formatter": "{value} ",
                "rotate": 0,
                "margin": 8,
                "textStyle": {
                    "fontSize": 12
                }
            }
        }
    ],
*/

type YAxis struct {
	Show          bool       `json:"show"`
	NameLoc       string     `json:"nameLoc"`
	NameGap       float32    `json:"nameGap"`
	NameTextStyle *TextStyle `json:"nameTextStyle"`
	AxisTick      *AxisTick  `json:"axisTick"`
	Inverse       bool       `json:"inverse"`
	BoundaryGap   bool       `json:"boundaryGap"`
	Type          string     `json:"type"`
	SplitLine     *SplitLine `json:"splitLine"`
	AxisLine      *AxisLine  `json:"axisLine"`
	AxisLabel     *AxisLabel `json:"axisLabel"`
}

func DefaultYAxis() *YAxis {
	return &YAxis{
		Show:    true,
		NameLoc: "middle",
		NameGap: 25,
		NameTextStyle: &TextStyle{
			FontSize: 14,
		},
		AxisTick: &AxisTick{
			AlignWithLabel: false,
		},
		Inverse:     false,
		BoundaryGap: true,
		Type:        "value",
		SplitLine: &SplitLine{
			Show: true,
		},
		AxisLine: &AxisLine{
			LineStyle: &LineStyle{
				Width: 1,
			},
		},
		AxisLabel: &AxisLabel{
			Interval: "auto",
			Formater: "{value}",
			Rotate:   0,
			Margin:   8,
			TextStyle: &TextStyle{
				FontSize: 12,
			},
		},
	}
}

/*
"legend": [
        {
            "data": [
                "precipitation",
                "evaporation"
            ],
            "selectedMode": "multiple",
            "show": true,
            "left": "center",
            "top": "top",
            "orient": "horizontal",
            "textStyle": {
                "fontSize": 12
            }
        }
    ],
*/

type Legend struct {
	Data        []string   `json:"data"`
	SelectMode  string     `json:"selectMode"`
	Show        bool       `json:"show"`
	Left        string     `json:"center"`
	Top         string     `json:"top"`
	Orientation string     `json:"orient"`
	TextStyle   *TextStyle `json:"textStyle"`
}

func DefaultLegend(names []string) *Legend {
	return &Legend{
		Data:        names,
		SelectMode:  "multiple",
		Show:        true,
		Left:        "center",
		Top:         "top",
		Orientation: "horizontal",
		TextStyle: &TextStyle{
			FontSize: 12,
		},
	}
}
