package sentiment

import (
	"gonlp/classification"
	"gonlp/normal"
	"gonlp/seg"
	"gonlp/util"
	"log"
	"strings"
)

type Sentiment struct {
	Bayes  classification.Bayes
	Normal normal.Normal
	Seg    seg.Seg
}

func InitSentiment() *Sentiment {
	classifier := classification.InitBayes()
	classifier.Load("sentiment.marshal")
	return &Sentiment{classifier, normal.InitNormal(), seg.InitSeg()}
}

func (sentiment *Sentiment) Save(fname string) {
	sentiment.Bayes.Save(fname)
}

func (sentiment *Sentiment) Handle(doc string) []string {
	return sentiment.Normal.FilterStop(sentiment.Seg.Seg(doc))
}

func (sentiment *Sentiment) Train(negFile string, posFile string) {
	negDocs := []string{}
	posDocs := []string{}
	file, reader := util.GetReader(negFile)
	defer file.Close()
	if reader != nil {
		err := util.ReadFile(reader, func(line string) {
			line = strings.Trim(line, " ")
			negDocs = append(negDocs, line)
		})
		if err != nil {
			log.Println("读取" + negFile + "出错")
		}
	}
	file, reader = util.GetReader(posFile)
	defer file.Close()
	if reader != nil {
		err := util.ReadFile(reader, func(line string) {
			line = strings.Trim(line, " ")
			posDocs = append(posDocs, line)
		})
		if err != nil {
			log.Println("读取" + posFile + "出错")
		}
	}
	data := []([]interface{}){}
	for _, sent := range negDocs {
		words := sentiment.Handle(sent)
		arr := make([]interface{}, 2)
		arr[0] = words
		arr[1] = "neg"
		data = append(data, arr)
	}
	for _, sent := range posDocs {
		words := sentiment.Handle(sent)
		arr := make([]interface{}, 2)
		arr[0] = words
		arr[1] = "pos"
		data = append(data, arr)
	}
	sentiment.Bayes.Train(data)
}
