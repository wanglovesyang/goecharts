package goecharts

type SeriesMaker func(val []float32, name string, tp string) *Series
