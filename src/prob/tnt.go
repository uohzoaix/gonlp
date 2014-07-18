package prob

import (
	"fmt"
)

type TnT struct {
	num           int
	l1            float64
	l2            float64
	l3            float64
	status        []string
	wd, eos, eosd AddOne
	uni, bi, tri  Normal
	word          map[string][]string
	trans         map[interface{}]float64
}

func (tnt TnT) init(num int) TnT {
	tnt.num = num
	tnt.l1 = 0.0
	tnt.l2 = 0.0
	tnt.l3 = 0.0
	tnt.status = make([]string, 0)
	tnt.wd = AddOne{Base{make(map[string]float64), 0.0, 1}}
	tnt.eos = AddOne{Base{make(map[string]float64), 0.0, 1}}
	tnt.eosd = AddOne{Base{make(map[string]float64), 0.0, 1}}
	tnt.uni = Normal{Base{make(map[string]float64), 0.0, 0}}
	tnt.bi = Normal{Base{make(map[string]float64), 0.0, 0}}
	tnt.tri = Normal{Base{make(map[string]float64), 0.0, 0}}
	tnt.word = make(map[string][]string)
	tnt.trans = make(map[interface{}]float64)
	return tnt
}

func (tnt TnT) getNum() int {
	return tnt.num
}

func main() {
	tnt := TnT{}
	tnt = tnt.init(10)
	fmt.Println(tnt.getNum())
}
