package util

import (
	"container/list"
	"encoding/json"
	"gonlp/pojo"
	"math"
	"strings"
)

type TnT struct {
	num           int
	l1            float64
	l2            float64
	l3            float64
	status        []string
	wd, eos, eosd *AddOne
	uni, bi, tri  *Normal
	word          map[string]*list.List
	trans         map[interface{}]float64
}

func NewTnT(num int) *TnT {
	return &TnT{num, 0.0, 0.0, 0.0, make([]string, 0), InitAddOne(), InitAddOne(), InitAddOne(), InitNormal(), InitNormal(), InitNormal(), make(map[string]*list.List), make(map[interface{}]float64)}
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

func (tnt *TnT) train(data *list.List) {
	now := []string{}
	now = append(now, "BOS")
	now = append(now, "BOS")
	for wtList := data.Front(); wtList != nil; wtList = wtList.Next() {
		tnt.bi.Add("BOS-BOS", 1)
		tnt.uni.Add("BOS", 2)
		newList := wtList.Value.(*list.List)
		for wt := newList.Front(); wt != nil; wt = wt.Next() {
			wordTag := wt.Value.(*pojo.WordTag)
			now = append(now, wordTag.GetTag())
			tupleStr := strings.Join(now[1:], "-")
			tnt.status = append(tnt.status, wordTag.GetTag())
			tnt.wd.Add(wordTag.ToString(), 1)
			tnt.eos.Add(tupleStr, 1)
			tnt.eosd.Add(wordTag.GetTag(), 1)
			tnt.uni.Add(wordTag.GetTag(), 1)
			tnt.bi.Add(tupleStr, 1)
			tnt.tri.Add(strings.Join(now, "-"), 1)
			if _, ok := tnt.word[wordTag.GetWord()]; !ok {
				tnt.word[wordTag.GetWord()] = list.New()
			}
			if !Contains(tnt.word[wordTag.GetWord()], wordTag.GetTag()) {
				tnt.word[wordTag.GetWord()].PushBack(wordTag.GetTag)
			}
			now = now[1:]
		}
	}
}

func Contains(list *list.List, elem interface{}) bool {
	for t := list.Front(); t != nil; t = t.Next() {
		if t.Value == elem {
			return true
		}
	}
	return false
}

//func ToString(slice []string) string {
//	result
//	for _, val := range slice {

//	}
//}

//func main() {
//	tnt := TnT{}
//	tnt = tnt.init(10)
//	fmt.Println(tnt.getNum())
//}
