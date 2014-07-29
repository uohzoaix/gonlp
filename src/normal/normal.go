package normal

import (
	"gonlp/util"
	"log"
	"regexp"
	"strings"
)

type Normal struct {
	Stop   []string
	Pinyin map[string]string
}

func InitNormal() Normal {
	return Normal{initStop("stopwords.txt"), initPinyin("pinyin.txt")}
}

func initStop(stopFile string) []string {
	result := []string{}
	file, reader := util.GetReader(stopFile)
	defer file.Close()
	if reader != nil {
		err := util.ReadFile(reader, func(line string) {
			line = strings.Trim(line, " ")
			result = append(result, line)
		})
		if err != nil {
			log.Println("normal load stopFile err", err)
		}
	}
	return result
}

func initPinyin(pyFile string) map[string]string {
	result := make(map[string]string)
	file, reader := util.GetReader(pyFile)
	defer file.Close()
	if reader != nil {
		err := util.ReadFile(reader, func(line string) {
			words := strings.SplitN(strings.Trim(line, " "), " ", 2)
			result[words[0]] = words[1]
		})
		if err != nil {
			log.Println("normal load pinyinFile err", err)
		}
	}
	return result
}

func (normal *Normal) FilterStop(words []string) []string {
	filters := []string{}
	for _, word := range words {
		find := false
		for _, val := range normal.Stop {
			if word == val {
				find = true
			}
		}
		if !find {
			filters = append(filters, word)
		}
	}
	return filters
}

func (normal *Normal) Zh2hans(sentence string) string {
	zh := &ZH{}
	zh.InitZH()
	return zh.Transfer(sentence)
}

func (normal *Normal) GetSentence(doc string) []string {
	sentences := []string{}
	delimiter := regexp.MustCompile("[，。？！；]")
	for _, line := range strings.SplitN(doc, "\r\n", -1) {
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}
		line = delimiter.ReplaceAllString(line, " ")
		for _, sent := range regexp.MustCompile("[^\\s]+").FindAllString(line, -1) {
			sent = strings.Trim(sent, " ")
			if sent == "" {
				continue
			}
			sentences = append(sentences, sent)
		}
	}
	return sentences
}

func (normal *Normal) GetPinyin(word string) string {
	for key, val := range normal.Pinyin {
		if key == word {
			return val
		}
	}
	var ret = ""
	for _, w := range word {
		for key, val := range normal.Pinyin {
			if key == string(w) {
				ret += val
			}
		}
	}
	return ret
}
