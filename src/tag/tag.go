package tag

import (
	"container/list"
	"gonlp/pojo"
	"gonlp/util"
	"log"
	"regexp"
	"strings"
)

type Tag struct {
	Tnt *util.TnT
}

func InitTag() *Tag {
	tnt := util.NewTnT(0)
	tnt.Load("../tag/tag.marshal")
	return &Tag{tnt}
}

func (tag *Tag) Save(fname string) {
	tag.Tnt.Save(fname)
}

func (tag *Tag) Train(fileName string) {
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
		tag.Tnt.Train(wordTags)

	}
}

func (tag *Tag) Tag(words []string) []string {
	results := tag.TagAll(words)
	tags := []string{}
	for _, result := range results {
		tags = append(tags, result.GetCh())
	}
	return tags
}

func (tag *Tag) TagAll(words []string) []pojo.Result {
	return tag.Tnt.Tag(words)
}

//func main() {
//	tag := InitTag()
//	tag.Train("199801.txt")
//	tag.Save("tag.marshal")
//}
