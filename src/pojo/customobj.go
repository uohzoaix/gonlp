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
	One string
	Two string
}

func InitPre(one string, two string) Pre {
	return Pre{one, two}
}

func (pre *Pre) GetOne() string {
	return pre.One
}

func (pre *Pre) GetTwo() string {
	return pre.Two
}

func (pre Pre) ToString() string {
	return pre.GetOne() + "-" + pre.GetTwo()
}

type Result struct {
	word string
	ch   string
}

func InitResult(word string, ch string) *Result {
	return &Result{word, ch}
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
	Prefix *Pre
	Score  float64
	Suffix string
}

func InitTag(prefix Pre, score float64, suffix string) *Tag {
	return &Tag{&prefix, score, suffix}
}

func (tag Tag) GetPre() *Pre {
	return tag.Prefix
}

func (tag Tag) GetScore() float64 {
	return tag.Score
}
func (tag Tag) GetSuffix() string {
	return tag.Suffix
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

type StageValue struct {
	Score float64
	Value string
}

func InitSV(score float64, value string) *StageValue {
	return &StageValue{score, value}
}

func (sv *StageValue) GetScore() float64 {
	return sv.Score
}
func (sv *StageValue) GetValue() string {
	return sv.Value
}
