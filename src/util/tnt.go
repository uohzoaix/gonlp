package util

import (
	"container/list"
	"encoding/json"
	"fmt"
	"gonlp/pojo"
	"math"
	"sort"
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
	word          map[string]([]string)
	trans         map[interface{}]float64
}

func NewTnT(num int) *TnT {
	if num == 0 {
		num = 100
	}
	return &TnT{num, 0.0, 0.0, 0.0, []string{}, InitAddOne(), InitAddOne(), InitAddOne(), InitNormal(), InitNormal(), InitNormal(), make(map[string]([]string)), make(map[interface{}]float64)}
}

func (tnt TnT) getNum() int {
	return tnt.num
}

func (tnt *TnT) Save(fname string) {
	saveToFile(tnt, fname)
}

func (tnt *TnT) Load(fname string) {
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

func (tnt *TnT) Train(data *list.List) {
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
				tnt.word[wordTag.GetWord()] = []string{}
			}
			if !Contains(tnt.word[wordTag.GetWord()], wordTag.GetTag()) {
				tnt.word[wordTag.GetWord()] = append(tnt.word[wordTag.GetWord()], wordTag.GetTag())
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

func (tnt *TnT) Tag(data []string) []pojo.Result {
	tags := make([]pojo.Tag, tnt.getNum())
	tags = append(tags, pojo.InitTag(pojo.InitPre("BOS", "BOS"), 0.0, ""))
	stage := make(map[pojo.Pre]pojo.StageValue)

	for _, ch := range data {
		samples := tnt.status
		if val, ok := tnt.word[ch]; ok {
			samples = val
		}
		for _, s := range samples {
			wd := math.Log(tnt.wd.Get(s+"-"+ch)) - math.Log(tnt.uni.Get(s))
			for _, tag := range tags {
				p := tag.GetScore() + wd + tnt.trans[tag.GetPre().ToString()+"-"+s]
				pre := pojo.Pre{tag.GetPre().GetTwo(), s}
				if sv, ok := stage[pre]; !ok || p > sv.GetScore() {
					if tag.GetSuffix() == "" {
						stage[pre] = pojo.StageValue{p, s}
					} else {
						stage[pre] = pojo.StageValue{p, tag.GetSuffix() + "-" + s}
					}
				}
			}
		}
		tags = []pojo.Tag{}
		for key, val := range stage {
			tags = append(tags, pojo.Tag{key, val.GetScore(), val.GetValue()})
		}
		temp := make(map[float64]pojo.Tag)
		for _, val := range tags {
			temp[val.GetScore()] = val
		}
		mk := make([]float64, len(temp))
		i := 0
		for k, _ := range temp {
			mk[i] = k
			i++
		}
		sort.Sort(sort.Reverse(sort.Float64Slice(mk)))
		tags = []pojo.Tag{}
		i = 0
		for _, v := range mk {
			if i < tnt.getNum() {
				tags = append(tags, temp[v])
			}
			i++
		}
	}
	tags = []pojo.Tag{}
	for key, val := range stage {
		score := val.GetScore() + tnt.getEos(key.GetTwo())
		tags = append(tags, pojo.Tag{key, score, val.GetValue()})
	}
	temp := make(map[float64]pojo.Tag)
	for _, val := range tags {
		temp[val.GetScore()] = val
	}
	mk := make([]float64, len(temp))
	i := 0
	for k, _ := range temp {
		mk[i] = k
		i++
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(mk)))
	tags = []pojo.Tag{}
	for _, v := range mk {
		tags = append(tags, temp[v])
	}
	results := []pojo.Result{}
	tagArr := strings.Split(tags[0].GetSuffix(), "-")
	if len(tagArr) != len(data) {
		panic("出错了！")
	}
	for i, val := range data {
		results = append(results, pojo.Result{val, tagArr[i]})
	}
	return results
}

func Contains(arr []string, elem interface{}) bool {
	for _, val := range arr {
		if val == elem.(string) {
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
