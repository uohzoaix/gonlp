package seg

import (
	"container/list"
	"gonlp/pojo"
	"gonlp/util"
	"log"
	"regexp"
	"strings"
)

type Seg struct {
	segger *util.TnT
}

func InitSeg() *Seg {
	segger := util.NewTnT(0)
	segger.Load("../seg/seg.marshal")
	return &Seg{segger}
}

func (seg *Seg) Save(fname string) {
	log.Println("save to ", fname)
	seg.segger.Save(fname)
}

func (seg *Seg) Load(fname string) {
	seg.segger.Load(fname)
}

func (seg *Seg) Train(fileName string) {
	wordTags := list.New()
	file, reader := util.GetReader(fileName)
	defer file.Close()
	if reader != nil {
		err := util.ReadFile(reader, func(line string) {
			line = strings.Trim(line, " ")
			wts := list.New()
			wordTags.PushBack(wts)
			for _, str := range regexp.MustCompile("[^\\s]+").FindAllString(line, -1) {
				if strings.Trim(str, " ") != "" {
					wts.PushBack(pojo.WordTag{strings.Split(str, "/")[0], strings.Split(str, "/")[1]})
				}
			}
		})
		if err != nil {
			log.Println("seg train err", err)
		}
		seg.segger.Train(wordTags)
	}
}

func (seg *Seg) SegOne(sentence string) []string {
	ret := []string{}
	data := []string{}
	for _, ch := range sentence {
		data = append(data, string(ch))
	}
	var oneres string
	results := seg.segger.Tag(data)
	for _, result := range results {
		if result.GetCh() == "s" {
			ret = append(ret, result.GetWord())
		} else if result.GetCh() == "e" {
			oneres += result.GetWord()
			ret = append(ret, oneres)
			oneres = ""
		} else {
			oneres += result.GetWord()
		}
	}
	return ret
}

func (seg *Seg) Seg(sent string) []string {
	words := []string{}
	for _, str := range regexp.MustCompile("([\u4E00-\u9FA5]+)").FindAllString(sent, -1) {
		str = strings.Trim(str, " ")
		if str == "" {
			continue
		}
		hzRegexp := regexp.MustCompile("[\u4e00-\u9fa5]")
		if ok := hzRegexp.MatchString(str); ok {
			words = append(words, seg.SegOne(str)...)
		} else {
			for _, word := range regexp.MustCompile("[^\\s]+").FindAllString(str, -1) {
				word = strings.Trim(word, " ")
				if word != "" {
					words = append(words, word)
				}
			}
		}
	}
	return words
}

func main() {
	seg := InitSeg()
	//seg.Train("data.txt") // 主要是用来放置一些简单快速的中文分词和词性标注的程序
	//seg.Save("seg1.marshal")
	log.Println(seg.Seg("这个东西真心很赞"))
}
