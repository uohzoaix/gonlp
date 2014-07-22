package util

type common interface {
	add(key string, value int)
	getData()
}

type Base struct {
	Data        map[string]float64
	Total, None float64
}

func InitBase() *Base {
	return &Base{make(map[string]float64), 0.0, 0}
}

func (base *Base) Exist(key string) bool {
	_, ok := base.Data[key]
	return ok
}
func (base *Base) GetSum() float64 {
	return base.Total
}
func (base *Base) Get(key string) float64 {
	if base.Exist(key) {
		return base.Data[key]
	}
	return base.None
}
func (base *Base) Frequency(key string) float64 {
	return base.Get(key) / base.Total
}
func (base *Base) Samples() []float64 {
	v := make([]float64, 0, len(base.Data))
	for _, value := range base.Data {
		v = append(v, value)
	}
	return v
}

type AddOne struct {
	Base
}

func InitAddOne() *AddOne {
	return &AddOne{Base{make(map[string]float64), 0.0, 1}}
}

func (addOne *AddOne) Add(key string, value int) {
	addOne.Total += float64(value)
	if !addOne.Exist(key) {
		addOne.Data[key] = 1
		addOne.Total += float64(1)
	}
	addOne.Data[key] = addOne.Data[key] + float64(value)
}

func (addOne *AddOne) GetData() map[string]float64 {
	return addOne.Data
}

type GoodTuring struct {
	Base
	handled bool
}

func InitGoodTuring() *GoodTuring {
	return &GoodTuring{Base{make(map[string]float64), 0.0, 0}, false}
}

func (goodTuring *GoodTuring) Add(key string, value int) {
	if !goodTuring.Exist(key) {
		goodTuring.Data[key] = float64(0)
	}
	goodTuring.Data[key] = goodTuring.Data[key] + float64(value)
}

func (goodTuring *GoodTuring) GetData() map[string]float64 {
	return goodTuring.Data
}

func (goodTuring *GoodTuring) Get(key string) float64 {
	if !goodTuring.handled {
		goodTuring.handled = true
		result := Calc(goodTuring.Data)
		if f, ok := result[0].(float64); ok {
			goodTuring.None = f
		}
		if m, ok := result[1].(map[string]float64); ok {
			goodTuring.Data = m
		}
		goodTuring.Total = float64(0)
		for _, value := range goodTuring.Data {
			goodTuring.Total += value
		}
	}
	if !goodTuring.Exist(key) {
		return goodTuring.None
	}
	return goodTuring.Data[key]
}

type Normal struct {
	Base
}

func InitNormal() *Normal {
	return &Normal{Base{make(map[string]float64), 0.0, 0}}
}

func (normal *Normal) Add(key string, value int) {
	if !normal.Exist(key) {
		normal.Data[key] = float64(0)
	}
	normal.Data[key] = normal.Data[key] + float64(value)
	normal.Total += float64(value)
}

func (normal *Normal) GetData() map[string]float64 {
	return normal.Data
}
