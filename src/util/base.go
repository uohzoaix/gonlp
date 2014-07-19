package util

import (
	"../util"
	"fmt"
)

type common interface {
	add(key string, value int)
	getData()
}

type Base struct {
	Data        map[string]float64
	Total, None float64
}

func (base Base) Init() Base {
	base.Data = make(map[string]float64)
	base.Total = 0.0
	base.None = 0
	return base
}

func (base Base) exist(key string) bool {
	_, ok := base.Data[key]
	return ok
}
func (base Base) getSum() float64 {
	return base.Total
}
func (base Base) get(key string) float64 {
	if base.exist(key) {
		return base.Data[key]
	}
	return base.None
}
func (base Base) frequency(key string) float64 {
	return base.get(key) / base.Total
}
func (base Base) samples() []float64 {
	v := make([]float64, 0, len(base.Data))
	for _, value := range base.Data {
		v = append(v, value)
	}
	return v
}

type AddOne struct {
	Base
}

func (addOne AddOne) Init() AddOne {
	addOne.Init()
	addOne.None = 1
	return addOne
}

func (addOne AddOne) add(key string, value int) {
	addOne.Total += float64(value)
	if !addOne.exist(key) {
		addOne.Data[key] = 1
		addOne.Total += float64(1)
	}
	addOne.Data[key] = addOne.Data[key] + float64(value)
}

func (addOne AddOne) getData() map[string]float64 {
	return addOne.Data
}

type GoodTuring struct {
	Base
	handled bool
}

func (goodTuring GoodTuring) Init() GoodTuring {
	goodTuring.Init()
	goodTuring.handled = false
	return goodTuring
}

func (goodTuring GoodTuring) add(key string, value int) {
	if !goodTuring.exist(key) {
		goodTuring.Data[key] = float64(0)
	}
	goodTuring.Data[key] = goodTuring.Data[key] + float64(value)
}

func (goodTuring GoodTuring) getData() map[string]float64 {
	return goodTuring.Data
}

func (goodTuring GoodTuring) get(key string) float64 {
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
	if !goodTuring.exist(key) {
		return goodTuring.None
	}
	return goodTuring.Data[key]
}

type Normal struct {
	Base
}

func (normal Normal) add(key string, value int) {
	if !normal.exist(key) {
		normal.Data[key] = float64(0)
	}
	normal.Data[key] = normal.Data[key] + float64(value)
	normal.Total += float64(value)
}

func (normal Normal) getData() map[string]float64 {
	return normal.Data
}
