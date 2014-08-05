package util

import (
	"strings"
)

func ToString(tuple []string) string {
	return strings.Join(tuple, "-")
}

func FromString(str string) []string {
	strArr := strings.Split(str, "-")
	result := []string{}
	for _, val := range strArr {
		if val != "" {
			result = append(result, strings.Trim(val, " "))
		}
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
