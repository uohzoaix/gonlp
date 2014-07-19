package util

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"unsafe"
)

type MemFile struct {
}

func (memFile MemFile) loadFromMem(fname string, obj interface{}) {
	data := make(map[string]interface{})
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for i := 0; i < t.NumField(); i++ {
		//根据索引获取属性 StructField
		f := t.Field(i)
		//返回属性的对应值 相当于取值操作
		val := v.Field(i).Interface()
		if v.Field(i).Kind() == reflect.Slice {
			data[f.Name] = []reflect.Value{reflect.ValueOf(val)}
		} else if reflect.TypeOf(v.Field(i)).Name() == "Base" {
			data[f.Name] = (*Base)(unsafe.Pointer(&val))
		} else {
			data[f.Name] = val
		}
	}
	memFile.saveToFile(data, fname)
}

func (memFile MemFile) saveToFile(data map[string]interface{}, fname string) {
	fout, err := os.Create(fname)
	defer fout.Close()
	if err != nil {
		fmt.Println(fname, err)
		return
	}
	ret, _ := json.Marshal(data)
	fout.Write([]byte(ret))
}

func (memFile MemFile) loadFromFile(fname string, obj interface{}) []byte {
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
	return result
}
