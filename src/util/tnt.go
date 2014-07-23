package util

import (
	"container/list"
	"encoding/json"
	"fmt"
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
	return &TnT{num, 0.0, 0.0, 0.0, []string{}, InitAddOne(), InitAddOne(), InitAddOne(), InitNormal(), InitNormal(), InitNormal(), make(map[string]*list.List), make(map[interface{}]float64)}
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
		tnt.eos.Add(now[len(now)-1]+"-EOS", 1)
	}
	var tl1, tl2, tl3 float64
	for _, val := range tnt.tri.Samples() {
		now = FromString(val)
		c3 := tnt.tntDiv(tnt.tri.Get(val)-1, tnt.bi.Get(ToString(now[0:2]))-1)
		c2 := tnt.tntDiv(tnt.bi.Get(ToString(now[1:]))-1, tnt.uni.Get(now[1])-1)
		c1 := tnt.tntDiv(tnt.uni.Get(now[2])-1, tnt.uni.GetSum()-1)
		result := tnt.tri.Get(ToString(now))
		if c3 >= c1 && c3 >= c2 {
			tl3 += result
		} else if c2 >= c1 && c2 >= c3 {
			tl2 += result
		} else if c1 >= c2 && c1 >= c3 {
			tl1 += result
		}
	}

	tnt.l1 = tl1 / (tl1 + tl2 + tl3)
	tnt.l2 = tl2 / (tl1 + tl2 + tl3)
	tnt.l3 = tl3 / (tl1 + tl2 + tl3)
	newStatus := &StringSet{make(map[string]bool)}
	for _, val := range tnt.status {
		newStatus.Add(val)
	}
	newStatus.Add("BOS")
	for key1, _ := range newStatus.set {
		for key2, _ := range newStatus.set {
			for _, val := range tnt.status {
				if key1 == "BOS" && key2 == "BOS" && val == "s" {
					fmt.Println("aa")
				}
				uni := tnt.l1 * tnt.uni.Frequency(val)
				bi := tnt.tntDiv(tnt.l2*tnt.bi.Get(key2+"-"+val), tnt.uni.Get(key2))
				tri := tnt.tntDiv(tnt.l3*tnt.tri.Get(key1+"-"+key2+"-"+val), tnt.bi.Get(key1+"-"+key2))
				tnt.trans[key1+"-"+key2+"-"+val] = math.Log(uni + bi + tri)
			}
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

type StringSet struct {
	set map[string]bool
}

func (set *StringSet) Add(str string) bool {
	_, found := set.set[str]
	set.set[str] = true
	return !found
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
