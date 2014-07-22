package util

import (
	"encoding/json"
	"math"
)

type TnT struct {
	num           int
	l1            float64
	l2            float64
	l3            float64
	status        []string
	wd, eos, eosd *AddOne
	uni, bi, tri  *Normal
	word          map[string][]string
	trans         map[interface{}]float64
}

func NewTnT(num int) *TnT {
	return &TnT{num, 0.0, 0.0, 0.0, make([]string, 0), InitAddOne(), InitAddOne(), InitAddOne(), InitNormal(), InitNormal(), InitNormal(), make(map[string][]string), make(map[interface{}]float64)}
}

func (tnt TnT) getNum() int {
	return tnt.num
}

func (tnt *TnT) save(fname string) {
	saveToFile(tnt, fname)
}

func (tnt *TnT) load(fname string) {
	result := loadFromFile(fname, tnt)
	if result != nil {
		//MemFile.loadToMem(result, tnt)
		json.Unmarshal(result, &tnt)
	}
}

func (tnt *TnT) tntDiv(v1 float64, v2 float64) float64 {
	if v2 == 0.0 {
		return v2
	}
	return v1 / v2
}

func (tnt *TnT) getEos(tag string) float64 {
	if tnt.eosd.Exist(tag) {
		return math.Log((float64)(1.0 / len(tnt.status)))
	}
	return math.Log(tnt.eos.Get(tag+"-EOS")) - math.Log(tnt.eosd.Get(tag))
}

//func main() {
//	tnt := TnT{}
//	tnt = tnt.init(10)
//	fmt.Println(tnt.getNum())
//}
