package util

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
)

type MemFile struct {
}

func loadFromMem(fname string, tnt *TnT) {
	//data := make(map[string]interface{})
	//t := reflect.TypeOf(obj)
	//v := reflect.ValueOf(obj)
	//for i := 0; i < t.NumField(); i++ {
	//	//根据索引获取属性 StructField
	//	f := t.Field(i)
	//	//返回属性的对应值 相当于取值操作
	//	val := v.Field(i).Interface()
	//	if v.Field(i).Kind() == reflect.Slice {
	//		data[f.Name] = []reflect.Value{reflect.ValueOf(val)}
	//	} else if reflect.TypeOf(v.Field(i)).Name() == "Base" {
	//		data[f.Name] = (*Base)(unsafe.Pointer(&val))
	//	} else {
	//		data[f.Name] = val
	//	}
	//}
	SaveToFile(tnt, fname)
}

func SaveToFile(obj interface{}, fname string) {
	log.Println("SaveToFile() ", fname)
	fout, err := os.Create(fname)
	defer fout.Close()
	if err != nil {
		fmt.Println(fname, err)
		return
	}
	encoder := json.NewEncoder(fout)
	encoder.Encode(obj)
	//ret, _ := json.Marshal(obj)
	//fout.Write([]byte(ret))
}

func LoadFromFile(fname string) []byte {
	fin, err := os.Open(fname)
	defer fin.Close()
	if err != nil {
		fmt.Println(fname, err)
		return nil
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
	return result
}

//TODO:field.set改为调用struct的set方法
func loadToMem(result []byte, obj interface{}) {
	data := make(map[string]interface{})
	json.Unmarshal(result, data)
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(&obj).Elem()
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
							field.Set(reflect.ValueOf(val.([]string)))
						}
					}
				}
				break
			}
		}
	}
}

func bayesLoadToMem(result []byte, obj interface{}) {

}

func GetReader(fileName string) (*os.File, *bufio.Reader) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	return file, bufio.NewReader(file)
}

func ReadFile(reader *bufio.Reader, handle func(string)) error {
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
