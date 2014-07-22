package pojo

type ClassifyResult struct {
	ret  string
	prob float64
}

func (classifyResult ClassifyResult) GetRet() string {
	return classifyResult.ret
}
func (classifyResult ClassifyResult) GetProb() float64 {
	return classifyResult.prob
}

type Pre struct {
	one string
	two string
}

func (pre Pre) GetOne() string {
	return pre.one
}

func (pre Pre) GetTwo() string {
	return pre.two
}

func (pre Pre) ToString() string {
	return pre.GetOne() + "-" + pre.GetTwo()
}

type Result struct {
	word string
	ch   string
}

func (result Result) GetWord() string {
	return result.word
}
func (result Result) GetCh() string {
	return result.ch
}

func (result Result) ToString() string {
	return result.GetWord() + "-" + result.GetCh()
}

type Tag struct {
	prefix Pre
	score  float64
	suffix string
}

func (tag Tag) GetPre() Pre {
	return tag.prefix
}

func (tag Tag) GetScore() float64 {
	return tag.score
}
func (tag Tag) Suffix() string {
	return tag.suffix
}

type WordTag struct {
	word string
	tag  string
}

func InitWT(word string, tag string) *WordTag {
	return &WordTag{word, tag}
}

func (wordTag WordTag) GetWord() string {
	return wordTag.word
}
func (wordTag WordTag) GetTag() string {
	return wordTag.tag
}

func (wordTag WordTag) ToString() string {
	return wordTag.GetTag() + "-" + wordTag.GetWord()
}
