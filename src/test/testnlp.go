package main

import (
	"../"
	"log"
)

func main() {
	nlp := gonlp.InitNlp("这个东西真心很赞")
	log.Println(nlp.Words())
	log.Println(nlp.Tags())
	log.Println(nlp.Sentiments())
	log.Println(nlp.Pinyin())
	nlp = gonlp.InitNlp("「繁體字」「繁體中文」的叫法在臺灣亦很常見。")
	log.Println(nlp.Han())
	text := "自然语言处理是计算机科学领域与人工智能领域中的一个重要方向。\r\n它研究能实现人与计算机之间用自然语言进行有效通信的各种理论和方法。\r\n自然语言处理是一门融语言学、计算机科学、数学于一体的科学。\r\n因此，这一领域的研究将涉及自然语言，即人们日常使用的语言，\r\n所以它与语言学的研究有着密切的联系，但又有重要的区别。\r\n自然语言处理并不是一般地研究自然语言，\r\n而在于研制能有效地实现自然语言通信的计算机系统，\r\n特别是其中的软件系统。因而它是计算机科学的一部分。"
	nlp = gonlp.InitNlp(text)
	log.Println(nlp.Keywords(5))
	log.Println(nlp.Summary(3))
	log.Println(nlp.Sentences())
	test := [][]string{}
	test1 := []string{}
	test1 = append(test1, "这篇")
	test1 = append(test1, "文章")
	test = append(test, test1)
	test1 = []string{}
	test1 = append(test1, "那篇")
	test1 = append(test1, "论文")
	test = append(test, test1)
	test1 = []string{}
	test1 = append(test1, "这个")
	test = append(test, test1)
	nlp = gonlp.InitNlp(test)
	log.Println(nlp.Tf())
	log.Println(nlp.Idf())
	test1 = []string{}
	test1 = append(test1, "文章")
	log.Println(nlp.Sim(test1))
}
