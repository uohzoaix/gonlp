package util

import (
	"math"
	"sort"
)

func getz(r []float64, nr []float64) []float64 {
	z := []float64{}
	z = append(z, 2*nr[0]/r[1])
	for i := 0; i < len(nr)-2; i++ {
		z = append(z, 2*nr[i+1]/(r[i+2]-r[i]))
	}
	z = append(z, nr[len(nr)-1]/(r[len(r)-1]-r[len(r)-2]))
	return z
}

func leastSquare(rd []float64, zd []float64) []float64 {
	result := []float64{}
	var sumX, sumY, sumXY, square, b float64
	for _, value := range rd {
		sumX += value
	}
	for _, value := range zd {
		sumY += value
	}
	meanX := sumX / float64(len(rd))
	meanY := sumY / float64(len(zd))
	for i := 0; i < len(rd); i++ {
		square += float64(math.Pow(rd[i]-meanX, 2))
		sumXY += (rd[i] - meanX) * (zd[i] - meanY)
	}
	b = sumXY / square
	result = append(result, meanY-b*meanX)
	result = append(result, b)
	return result
}

func Calc(data map[string]float64) []interface{} {
	var r, rd, zd, nr, prob, z, values []float64
	for _, value := range data {
		values = append(values, value)
	}
	sort.Float64s(values)
	for _, value := range values {
		if len(r) == 0 || r[len(r)-1] != value {
			r = append(r, value)
			nr = append(nr, float64(1))
		} else {
			if len(nr) == 0 {
				nr = append(nr, float64(1))
			} else {
				nr[len(nr)-1] = nr[len(nr)-1] + 1
			}
		}
	}
	var total float64
	rr := make(map[float64]int)
	for pos, value := range r {
		if pos < len(nr) {
			total += value * nr[pos]
		}
		rr[value] = pos
		rd = append(rd, math.Log(value))
	}
	z = getz(r, nr)
	for _, value := range z {
		zd = append(zd, math.Log(value))
	}
	square := leastSquare(rd, zd)
	var useGoodTuring = false
	nr = append(nr, math.Exp(square[0]+square[1]*math.Log(r[len(r)-1]+1)))
	for i := 0; i < len(r); i++ {
		var goodTuring = (r[i] + 1) * (math.Exp(square[1] * (math.Log(r[i]+1) - math.Log(r[i]))))
		var turing = goodTuring
		if i+1 < len(r) {
			turing = (r[i] + 1) * nr[i+1] / nr[i]
		}
		var diff = math.Pow(math.Pow(r[i]+1, 2)/nr[i]*nr[i+1]/nr[i]*(1+nr[i+1]/nr[i]), 0.5) * 1.65
		if !useGoodTuring && math.Abs(goodTuring-turing) > diff {
			prob = append(prob, turing)
		} else {
			useGoodTuring = true
			prob = append(prob, goodTuring)
		}
	}
	var sump = float64(0)
	for pos, value := range nr {
		if pos < len(prob) {
			sump += value * prob[pos]
		}
	}
	for i := 0; i < len(prob); i++ {
		prob[i] = (1 - nr[0]/total) * prob[i] / sump
	}
	mixResult := make([]interface{}, 0)
	result := make(map[string]float64)
	for key, value := range data {
		result[key] = prob[rr[value]]
	}
	mixResult = append(mixResult, nr[0]/total/total)
	mixResult = append(mixResult, result)
	return mixResult
}
