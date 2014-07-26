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

func (seg *Seg) InitSeg() {
	seg.segger = util.NewTnT(0)
}

func (seg *Seg) Save(fname string) {
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
			words := strings.Fields(line)
			line = strings.Join(words, "")
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

func (seg *Seg) Seg(sentence string) []string {
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
