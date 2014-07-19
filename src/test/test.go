package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"unsafe"
)

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

	fname := "testfile.txt"
	fout, err := os.Create(fname)
	defer fout.Close()
	if err != nil {
		fmt.Println(fname, err)
		return
	}
	ret, _ := json.Marshal(test)
	fout.Write([]byte(ret))

}
