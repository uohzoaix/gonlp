package sim

import (
	"math"
)

type Bm25 struct {
	D     int
	Avgdl float64
	Docs  []([]string)
	F     [](map[string]int)
	Df    map[string]int
	Idf   map[string]float64
	K1    float64
	B     float64
}

func InitBm25(docs []([]string)) *Bm25 {
	f := [](map[string]int){}
	df := make(map[string]int)
	idf := make(map[string]float64)
	var sum = 0
	for _, doc := range docs {
		sum += len(doc)
		tmp := make(map[string]int)
		for _, word := range doc {
			if val, ok := tmp[word]; ok {
				tmp[word] = val + 1
			} else {
				tmp[word] = 1
			}
		}
		f = append(f, tmp)
		for key, _ := range tmp {
			if value, ok := df[key]; ok {
				df[key] = value + 1
			} else {
				df[key] = 1
			}
		}
	}
	for k, v := range df {
		idf[k] = math.Log(float64(len(docs)-v)+0.5) - math.Log(float64(v)+0.5)
	}
	return &Bm25{len(docs), float64(sum) / float64(len(docs)), docs, f, df, idf, 1.5, 0.75}
}

func (bm25 *Bm25) GetF() [](map[string]int) {
	return bm25.F
}

func (bm25 *Bm25) GetIdf() map[string]float64 {
	return bm25.Idf
}

func (bm25 *Bm25) Sim(doc []string, index int) float64 {
	score := 0.0
	for _, word := range doc {
		if _, ok := bm25.F[index][word]; !ok {
			continue
		}
		score += bm25.Idf[word] * float64(bm25.F[index][word]) * (bm25.K1 + 1) / (float64(bm25.F[index][word]) + bm25.K1*(1-bm25.B+bm25.B*float64(bm25.D)/bm25.Avgdl))
	}
	return score
}

func (bm25 *Bm25) Simall(doc []string) []float64 {
	scores := []float64{}
	for i := 0; i < bm25.D; i++ {
		scores = append(scores, bm25.Sim(doc, i))
	}
	return scores
}
