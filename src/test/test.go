package main

import (
	"bufio"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"gonlp/pojo"
	"gonlp/util"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
	"unsafe"
)

type Test struct {
	Name  string
	Teles []string
}

func (test *Test) SetName(name string) {
	test.Name = name
}

func (test *Test) SetTeles(teles []string) {
	test.Teles = teles
}

func main() {
	arr := make([]string, 0)
	arr = append(arr, "fdf")
	arr = append(arr, "gfgg")
	test := Test{"aa", arr}
	t := reflect.TypeOf(test)
	v := reflect.ValueOf(test)
	fmt.Println(reflect.TypeOf(test).Name())
	test1 := (*Test)(unsafe.Pointer(&test))
	//test.(Test)
	fmt.Println(test1.Name)
	for i := 0; i < t.NumField(); i++ {
		//根据索引获取属性 StructField
		f := t.Field(i)
		//返回属性的对应值 相当于取值操作
		val := v.Field(i).Interface()
		testval := [2]string{"fdf", "123"}
		fmt.Println([]reflect.Value{reflect.ValueOf(testval)})
		fmt.Println(v.Field(1).Kind() == reflect.Slice)
		//data[f.Name] = val
		fmt.Printf("%6s : %v = %v \n", f.Name, f.Type, val)
	}

	fname := "testfile.txt"
	fout, err := os.Create(fname)
	defer fout.Close()
	if err != nil {
		fmt.Println(fname, err)
		return
	}
	ret, _ := json.Marshal(test)
	fout.Write([]byte(ret))

	fin, err := os.Open(fname)
	defer fin.Close()
	if err != nil {
		fmt.Println(fname, err)
		return
	}
	var result, temp []byte
	var beginPos int
	buf := make([]byte, 2048*10000)
	for {
		num, _ := fin.Read(buf)
		if 0 == num {
			break
		}
		if len(result) > 0 {
			temp = result
			result = make([]byte, beginPos+num)
			copy(result[0:len(temp)], temp[:])
		} else {
			result = make([]byte, num)
		}
		copy(result[beginPos:beginPos+num], buf[:num])
		//os.Stdout.Write(buf[:num])
	}
	//fmt.Println(result)
	testfunc(ret, test)

	data := make(map[string]interface{})
	var testinter interface{}
	testinter = test
	json.Unmarshal([]byte(ret), &testinter)
	fmt.Println("testinter", testinter)
	t = reflect.TypeOf(test)
	v = reflect.ValueOf(&test).Elem()
	for key, val := range data {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Name == key {
				field := v.FieldByName(key)
				if field.IsValid() {
					if field.CanSet() {
						if field.Kind() == reflect.String {
							field.Set(reflect.ValueOf(val.(string)))
						} else if field.Kind() == reflect.Slice {
							val, ok := val.([]interface{})
							if ok {
							}
							fmt.Println(ok, val)
						}
					}
				}
				break
			}
		}
	}
	fmt.Println(test.Teles)

	addOne := util.InitAddOne()
	fmt.Println(addOne.None)
	fmt.Println(addOne.Exist("fdf"))

	dicts := list.New()
	wordTags := list.New()
	wordTags.PushBack(pojo.InitWT("fdf", "gfgg"))
	dicts.PushBack(wordTags)
	fmt.Println(dicts.Front().Value.(*list.List).Front().Value.(*pojo.WordTag).ToString())

	testmaplist := make(map[string]*list.List)
	testmaplist["fdf"] = list.New()

	tags := make([]pojo.Tag, 1)
	tags = append(tags, pojo.InitTag(pojo.InitPre("BOS", "BOS"), 0.0, ""))

	teststring := "你好啊,abcd"
	for _, val := range teststring {
		fmt.Println(string(val))
	}

	teststring = "  fdf fd gfhgh "
	fmt.Println(strings.Fields(teststring))
	fmt.Println(regexp.MustCompile("[^\\s]+").FindAllString(teststring, -1))

	fmt.Println(regexp.MustCompile("([\u4E00-\u9FA5]+)").FindAllString("你好 fdf 费大幅度发", -1))
	fmt.Println(strings.Trim("    fdf     ", " "))

	testsliceappend := []string{"dff", "dgf"}
	testappended := []string{"fdg", "ghhh"}
	testsliceappend = append(testsliceappend, testappended...)
	fmt.Println(testsliceappend)

}

func testfunc(bytes []byte, obj Test) (err error) {
	json.Unmarshal([]byte(bytes), &obj)
	fmt.Println("obj", obj.Name)
	file, reader := getReader("testfile.txt")
	defer file.Close()
	err = readFile(reader, func(line string) {
		if len(line) > 0 {
			//runes := utils.ToRunes(line)
			//dict[runes[0]] = runes[0]
			fmt.Println(line)
		}
	})
	return
}

func getReader(fileName string) (*os.File, *bufio.Reader) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil
	}
	return file, bufio.NewReader(file)
}

func readFile(reader *bufio.Reader, handle func(string)) error {
	if reader != nil {
		for {
			line, isPrefix, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if isPrefix {
				return errors.New("Error: unexcepted long line.")
			}
			handle(string(line))
		}
	}
	return nil
}
