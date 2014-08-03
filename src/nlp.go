package gonlp

import (
	"gonlp/normal"
	"gonlp/pojo"
	"gonlp/seg"
	"gonlp/sentiment"
	"gonlp/sim"
	"gonlp/tag"
	"gonlp/textrank"
	"log"
)

type NLP struct {
	doc        string
	bm25       *sim.Bm25
	seg        *seg.Seg
	normal     *normal.Normal
	sentiment  *sentiment.Sentiment
	tag        *tag.Tag
	textrank   *textrank.TextRank
	kvTextRank *textrank.KeyWordTextRank
}

func InitNlp(doc interface{}) *NLP {
	nlp := &NLP{}
	if v, ok := doc.(string); ok {
		nlp.doc = v
	} else if v, ok := doc.([][]string); ok {
		nlp.bm25 = sim.InitBm25(v)
	} else {
		log.Println("暂不支持第三种数据格式")
	}
	nlp.seg = seg.InitSeg()
	nlp.normal = normal.InitNormal()
	nlp.sentiment = sentiment.InitSentiment()
	nlp.tag = tag.InitTag()
	return nlp
}

func (nlp *NLP) Words() []string {
	return nlp.seg.Seg(nlp.doc)
}

func (nlp *NLP) Sentences() []string {
	return nlp.normal.GetSentence(nlp.doc)
}

func (nlp *NLP) Han() string {
	return nlp.normal.Zh2hans(nlp.doc)
}

func (nlp *NLP) Pinyin() string {
	words := nlp.Words()
	ret := ""
	for _, word := range words {
		ret += nlp.normal.GetPinyin(word) + " "
	}
	return ret
}

func (nlp *NLP) Sentiments() float64 {
	return nlp.sentiment.Classify(nlp.doc)
}

func (nlp *NLP) Tags() []pojo.Result {
	words := nlp.Words()
	tags := nlp.tag.Tag(words)
	result := []pojo.Result{}
	if len(words) == len(tags) {
		for i := 0; i < len(words); i++ {
			result = append(result, pojo.InitResult(words[i], tags[i]))
		}
	}
	return result
}

func (nlp *NLP) Tf() [](map[string]int) {
	return nlp.bm25.GetF()
}

func (nlp *NLP) Idf() map[string]float64 {
	return nlp.bm25.GetIdf()
}

func (nlp *NLP) Sim(doc []string) []float64 {
	return nlp.bm25.Simall(doc)
}

func (nlp *NLP) Summary(limit ...int) []string {
	size := 0
	if limit == nil || len(limit) == 0 {
		size = 5
	} else {
		size = limit[0]
	}
	doc := []([]string){}
	sents := nlp.Sentences()
	for _, sent := range sents {
		words := nlp.seg.Seg(sent)
		words = nlp.normal.FilterStop(words)
		doc = append(doc, words)
	}
	nlp.textrank = textrank.InitTR(doc)
	nlp.textrank.Solve()
	ret := []string{}
	indexes := nlp.textrank.TopIndex(size)
	for _, index := range indexes {
		ret = append(ret, sents[index])
	}
	return ret
}

func (nlp *NLP) Keywords(limit ...int) []string {
	size := 0
	if limit == nil || len(limit) == 0 {
		size = 5
	} else {
		size = limit[0]
	}
	doc := [][]string{}
	sents := nlp.Sentences()
	for _, sent := range sents {
		words := nlp.seg.Seg(sent)
		words = nlp.normal.FilterStop(words)
		doc = append(doc, words)
	}
	nlp.kvTextRank = textrank.InitKWRank(doc)
	nlp.kvTextRank.Solve()
	ret := []string{}
	indexes := nlp.kvTextRank.TopIndex(size)
	for _, index := range indexes {
		ret = append(ret, index)
	}
	return ret
}
