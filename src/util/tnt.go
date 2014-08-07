package util

import (
	"container/list"
	"encoding/json"
	//"fmt"
	"gonlp/pojo"
	//"log"
	"math"
	"sort"
	"strings"
)

type TnT struct {
	Num           int
	L1            float64
	L2            float64
	L3            float64
	Status        []string
	Wd, Eos, Eosd *AddOne
	Uni, Bi, Tri  *Normal
	Word          map[string]([]string)
	Trans         map[string]float64
}

func NewTnT(num int) *TnT {
	if num == 0 {
		num = 100
	}
	return &TnT{num, 0.0, 0.0, 0.0, []string{}, InitAddOne(), InitAddOne(), InitAddOne(), InitNormal(), InitNormal(), InitNormal(), make(map[string]([]string)), make(map[string]float64)}
}

func (tnt TnT) getNum() int {
	return tnt.Num
}

func (tnt *TnT) Save(fname string) {
	SaveToFile(tnt, fname)
}

func (tnt *TnT) Load(fname string) {
	result := LoadFromFile(fname)
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
	if !tnt.Eosd.Exist(tag) {
		return math.Log((float64)(1.0 / len(tnt.Status)))
	}
	return math.Log(tnt.Eos.Get(tag+"-EOS")) - math.Log(tnt.Eosd.Get(tag))
}

func (tnt *TnT) Train(data *list.List) {
	now := []string{}
	now = append(now, "BOS")
	now = append(now, "BOS")
	for wtList := data.Front(); wtList != nil; wtList = wtList.Next() {
		tnt.Bi.Add("BOS-BOS", 1)
		tnt.Uni.Add("BOS", 2)
		newList := wtList.Value.(*list.List)
		for wt := newList.Front(); wt != nil; wt = wt.Next() {
			wordTag := wt.Value.(pojo.WordTag)
			now = append(now, wordTag.GetTag())
			tupleStr := strings.Join(now[1:], "-")
			tnt.Status = append(tnt.Status, wordTag.GetTag())
			tnt.Wd.Add(wordTag.ToString(), 1)
			tnt.Eos.Add(tupleStr, 1)
			tnt.Eosd.Add(wordTag.GetTag(), 1)
			tnt.Uni.Add(wordTag.GetTag(), 1)
			tnt.Bi.Add(tupleStr, 1)
			tnt.Tri.Add(strings.Join(now, "-"), 1)
			if _, ok := tnt.Word[wordTag.GetWord()]; !ok {
				tnt.Word[wordTag.GetWord()] = []string{}
			}
			if !Contains(tnt.Word[wordTag.GetWord()], wordTag.GetTag()) {
				tnt.Word[wordTag.GetWord()] = append(tnt.Word[wordTag.GetWord()], wordTag.GetTag())
			}
			now = now[1:]
		}
		tnt.Eos.Add(now[len(now)-1]+"-EOS", 1)
	}
	var tl1, tl2, tl3 float64
	for _, val := range tnt.Tri.Samples() {
		now = FromString(val)
		c3 := tnt.tntDiv(tnt.Tri.Get(val)-1, tnt.Bi.Get(ToString(now[0:2]))-1)
		c2 := tnt.tntDiv(tnt.Bi.Get(ToString(now[1:]))-1, tnt.Uni.Get(now[1])-1)
		c1 := tnt.tntDiv(tnt.Uni.Get(now[2])-1, tnt.Uni.GetSum()-1)
		result := tnt.Tri.Get(ToString(now))
		if c3 >= c1 && c3 >= c2 {
			tl3 += result
		} else if c2 >= c1 && c2 >= c3 {
			tl2 += result
		} else if c1 >= c2 && c1 >= c3 {
			tl1 += result
		}
	}

	tnt.L1 = tl1 / (tl1 + tl2 + tl3)
	tnt.L2 = tl2 / (tl1 + tl2 + tl3)
	tnt.L3 = tl3 / (tl1 + tl2 + tl3)
	newStatus := &StringSet{make(map[string]bool)}
	for _, val := range tnt.Status {
		newStatus.Add(val)
	}
	newStatus.Add("BOS")
	var uni float64
	var bi float64
	var tri float64
	for key1, _ := range newStatus.set {
		for key2, _ := range newStatus.set {
			for _, val := range tnt.Status {
				if key1 == "BOS" && key2 == "BOS" && val == "s" {
					//fmt.Println("aa")
				}
				uni = tnt.L1 * tnt.Uni.Frequency(val)
				bi = tnt.tntDiv(tnt.L2*tnt.Bi.Get(key2+"-"+val), tnt.Uni.Get(key2))
				tri = tnt.tntDiv(tnt.L3*tnt.Tri.Get(key1+"-"+key2+"-"+val), tnt.Bi.Get(key1+"-"+key2))
				tnt.Trans[key1+"-"+key2+"-"+val] = math.Log(uni + bi + tri)
			}
		}
	}
}

func (tnt *TnT) Tag(data []string) []pojo.Result {
	//log.Println("tnt", tnt)
	tags := []pojo.Tag{}
	tags = append(tags, pojo.InitTag(pojo.InitPre("BOS", "BOS"), 0.0, ""))
	stage := make(map[string]pojo.StageValue)

	for _, ch := range data {
		stage = make(map[string]pojo.StageValue)
		samples := tnt.Status
		if val, ok := tnt.Word[ch]; ok {
			samples = val
		}
		for _, s := range samples {
			wd := math.Log(tnt.Wd.Get(s+"-"+ch)) - math.Log(tnt.Uni.Get(s))
			for _, tag := range tags {
				p := tag.GetScore() + wd + tnt.Trans[tag.GetPre().ToString()+"-"+s]
				pre := pojo.Pre{tag.GetPre().GetTwo(), s}.ToString()
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
			tags = append(tags, pojo.Tag{pojo.InitPre(strings.Split(key, "-")[0], strings.Split(key, "-")[1]), val.GetScore(), val.GetValue()})
		}
		temp := make(map[float64]pojo.Tag)
		for _, val := range tags {
			temp[val.GetScore()] = val
		}
		mk := []float64{}
		i := 0
		for k, _ := range temp {
			mk = append(mk, k)
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
		score := val.GetScore() + tnt.getEos(strings.Split(key, "-")[1])
		tags = append(tags, pojo.Tag{pojo.InitPre(strings.Split(key, "-")[0], strings.Split(key, "-")[1]), score, val.GetValue()})
	}
	temp := make(map[float64]pojo.Tag)
	for _, val := range tags {
		temp[val.GetScore()] = val
	}
	mk := []float64{}
	for k, _ := range temp {
		mk = append(mk, k)
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
