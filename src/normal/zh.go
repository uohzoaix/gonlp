package normal

import (
	"gonlp/util"
	"log"
	"strings"
)

type ZH struct {
	Zh2hans map[string]string
	Maxl    int
}

func (zh *ZH) InitZH() {
	zh.Zh2hans = make(map[string]string)
	for key, _ := range zh.Zh2hans {
		if len(key) > zh.Maxl {
			zh.Maxl = len(key)
		}
	}
}

func (zh *ZH) initHans(hansFile string) {
	file, reader := util.GetReader(hansFile)
	defer file.Close()
	if reader != nil {
		err := util.ReadFile(reader, func(line string) {
			words := strings.Split(strings.Trim(line, " "), "|")
			if len(words) == 1 {
				zh.Zh2hans[words[0]] = ""
			} else {
				zh.Zh2hans[words[0]] = words[1]
			}
		})
		if err != nil {
			log.Println("normal load hansFile err", err)
		}
	}
}

func (zh *ZH) Transfer(sentence string) string {
	var pos = 0
	var end = 0
	var ret = ""
	for pos < len(sentence) {
		find := false
		for i := zh.Maxl; i > 0; i-- {
			if pos+i > len(sentence) {
				end = len(sentence)
			} else {
				end = pos + i
			}
			word := sentence[pos:end]
			if val, ok := zh.Zh2hans[word]; ok {
				ret += val
				pos += i
				find = true
				break
			}
		}
		if !find {
			ret += sentence[pos : pos+1]
			pos += 1
		}
	}
	return ret
}
