package common

import "fmt"

func MapToString(m map[string]interface{}) string {
	var s string
	for k, v := range m {
		s += fmt.Sprintf("[%v]: %v ", k, v)
	}

	return s
}
