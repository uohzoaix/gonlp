package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type MemFile struct {
}

func (memFile MemFile) loadFromMem(name string, obj interface{}) {
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
}

type Test struct {
	Name  string
	Teles []string
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
}
