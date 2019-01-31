package goecharts

var funcmap = make(map[string]string)

func addFunc(k, v string) {
	funcmap[k] = v
}

func LeakMinFunc() string {
	addFunc("leak_min", "function(v){return v.min - Math.abs(v.min) * 0.05}")
	return "$leak_min$"
}

func LeakMaxFunc() string {
	addFunc("leak_max", "function(v){return v.max + Math.abs(v.max) * 0.05}")
	return "$leak_max$"
}
