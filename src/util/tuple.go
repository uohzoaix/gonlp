package util

import (
	"strings"
)

func ToString(tuple []string) string {
	return strings.Join(tuple, "-")
}

func FromString(str string) []string {
	strArr := strings.Split(str, "-")
	result := make([]string, len(strArr))
	for _, val := range strArr {
		result = append(result, val)
	}
	return result
}

//func main() {
//	test := make([]string, 0)
//	test = append(test, "fdf")
//	test = append(test, "  ")
//	test = append(test, "fdf")
//	test = append(test, "fdf")
//	fmt.Println(ToString(test))
//}
