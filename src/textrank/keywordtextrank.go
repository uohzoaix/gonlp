package textrank

import (
	"math"
	"sort"
)

type KeyWordTextRank struct {
	Docs    [][]string
	Words   map[string][]string
	Vertex  map[string]float64
	D       float64
	MaxIter int
	Mindiff float64
	Top     PairList
}

type Pair struct {
	Key   string
	Value float64
}

type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }

func InitKWRank(docs [][]string) *KeyWordTextRank {
	return &KeyWordTextRank{docs, make(map[string][]string), make(map[string]float64), 0.85, 200, 0.001, []Pair{}}
}

func (kwRank *KeyWordTextRank) Solve() {
	for _, doc := range kwRank.Docs {
		que := []string{}
		for _, word := range doc {
			if value, ok := kwRank.Words[word]; !ok {
				value = []string{}
				kwRank.Words[word] = value
				kwRank.Vertex[word] = 1.0
			}
			que = append(que, word)
			if len(que) > 5 {
				que = append(que[:0], que[1:]...)
			}
			for _, w1 := range que {
				for _, w2 := range que {
					if w1 == w2 {
						continue
					}
					kwRank.Words[w1] = append(kwRank.Words[w1], w2)
					kwRank.Words[w2] = append(kwRank.Words[w2], w1)
				}
			}
		}
	}
	iterNum := 0
	for {
		if iterNum < kwRank.MaxIter {
			iterNum++
			m := make(map[string]float64)
			maxDiff := 0.0
			for key, value := range kwRank.Words {
				m[key] = 1 - kwRank.D
				for _, j := range value {
					if key == j || len(kwRank.Words[j]) == 0 {
						continue
					}
					m[key] = m[key] + (kwRank.D / float64(len(kwRank.Words[j])) * kwRank.Vertex[j])
				}
				if math.Abs(m[key])-kwRank.Vertex[key] > maxDiff {
					maxDiff = math.Abs(m[key]) - kwRank.Vertex[key]
				}
			}
			kwRank.Vertex = m
			if maxDiff <= kwRank.Mindiff {
				break
			}
		}
	}
	kwRank.Top = sortMapByValue(kwRank.Vertex)
}

func (kwRank *KeyWordTextRank) TopIndex(limit int) []string {
	indexes := []string{}
	num := 0
	for _, pair := range kwRank.Top {
		if num == limit {
			break
		}
		indexes = append(indexes, pair.Key)
		num++
	}
	return indexes
}

func sortMapByValue(m map[string]float64) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
	}
	sort.Sort(p)
	return p
}
