package util

import (
	"fmt"
	"strings"
)

func ToString(tuple []string) string {
	return strings.Join(tuple, "-")
}

func main() {
	test := make([]string, 0)
	test = append(test, "fdf")
	test = append(test, "  ")
	test = append(test, "fdf")
	test = append(test, "fdf")
	fmt.Println(ToString(test))
}
