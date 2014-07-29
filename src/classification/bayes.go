package classification

import (
	"encoding/json"
	"gonlp/pojo"
	"gonlp/util"
	"log"
	"math"
)

type Bayes struct {
	D     map[string]util.AddOne
	Total float64
}

func InitBayes() Bayes {
	return Bayes{make(map[string]util.AddOne), 0.0}
}

func (bayes *Bayes) Save(fname string) {
	data := make(map[string]interface{})
	data["total"] = bayes.Total
	probdata := make(map[string]util.AddOne)
	for key, val := range bayes.D {
		probdata[key] = val
	}
	data["d"] = probdata
	util.SaveToFile(data, fname)
}

func (bayes *Bayes) Load(fname string) {
	result := util.LoadFromFile(fname)
	if result != nil {
		json.Unmarshal(result, &bayes)
	} else {
		log.Println("bayes读取" + fname + "文件出错！")
	}
}

func (bayes *Bayes) Train(data []([]interface{})) {
	for _, d := range data {
		c := d[1].(string)
		if _, ok := bayes.D[c]; !ok {
			bayes.D[c] = *util.InitAddOne()
		}
		for _, word := range d[0].([]string) {
			bayes.D[c].Add(word, 1)
		}
	}
	for key, _ := range bayes.D {
		bayes.Total = bayes.D[key].GetSum()
	}
}

func (bayes *Bayes) Classify(x []string) pojo.ClassifyResult {
	tmp := make(map[string]float64)
	for key, _ := range bayes.D {
		tmp[key] = 0.0
		for _, word := range x {
			tmp[key] = tmp[key] + math.Log(bayes.D[key].GetSum()) - math.Log(bayes.Total) + math.Log(bayes.D[key].Frequency(word))
		}
	}
	ret := ""
	prob := 0.0
	for key1, _ := range bayes.D {
		now := 0.0
		for otherKey, _ := range bayes.D {
			now += math.Exp(tmp[otherKey] - tmp[key1])
		}
		now = 1.0 / now
		if now > prob {
			ret = key1
			prob = now
		}
	}
	return pojo.ClassifyResult{ret, prob}
}
