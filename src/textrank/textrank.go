package textrank

import (
	"gonlp/sim"
	"gonlp/util"
	"math"
	"strconv"
)

type TextRank struct {
	Docs      [][]string
	Bm25      *sim.Bm25
	Num       int
	D         float64
	Weight    [][]float64
	WeightSum []float64
	Vertex    []float64
	MaxIter   int
	MinDiff   float64
	Top       util.PairList
}

func InitTR(docs [][]string) *TextRank {
	return &TextRank{docs, sim.InitBm25(docs), len(docs), 0.85, [][]float64{}, []float64{}, []float64{}, 200, 0.001, []util.Pair{}}
}

func (textRank *TextRank) Solve() {
	for _, doc := range textRank.Docs {
		scores := textRank.Bm25.Simall(doc)
		textRank.Weight = append(textRank.Weight, scores)
		sum := 0.0
		for _, score := range scores {
			sum += score
		}
		textRank.WeightSum = append(textRank.WeightSum, sum)
		textRank.Vertex = append(textRank.Vertex, 1.0)
	}
	iterNum := 0
	for {
		if iterNum < textRank.MaxIter {
			iterNum++
			m := []float64{}
			maxDiff := 0.0
			for i := 0; i < textRank.Num; i++ {
				m = append(m, 1-textRank.D)
				for j := 0; j < textRank.Num; j++ {
					if j == i || textRank.WeightSum[j] == 0.0 {
						continue
					}
					val := m[len(m)-1]
					m = append(m[:len(m)-1], m[len(m):]...)
					m = append(m, val+textRank.D*textRank.Weight[i][j]/textRank.WeightSum[j]*textRank.Vertex[j])
				}
				if math.Abs(m[len(m)-1]-textRank.Vertex[i]) > maxDiff {
					maxDiff = math.Abs(m[len(m)-1] - textRank.Vertex[i])
				}
			}
			textRank.Vertex = m
			if maxDiff <= textRank.MinDiff {
				break
			}
		}
	}
	result := make(map[string]float64)
	for k, v := range textRank.Vertex {
		result[strconv.Itoa(k)] = v
	}
	textRank.Top = util.SortMapByValue(result)
}

func (textRank *TextRank) TopIndex(limit int) []int {
	indexes := []int{}
	num := 0
	for _, pair := range textRank.Top {
		if num == limit {
			break
		}
		v, _ := strconv.Atoi(pair.Key)
		indexes = append(indexes, v)
		num++
	}
	return indexes
}

func (textRank *TextRank) Tops(limit int) [][]string {
	docs := [][]string{}
	num := 0
	for _, pair := range textRank.Top {
		if num == limit {
			break
		}
		v, _ := strconv.Atoi(pair.Key)
		docs = append(docs, textRank.Docs[v])
		num++
	}
	return docs
}
