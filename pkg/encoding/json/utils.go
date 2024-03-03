package json

import (
	"fmt"
	"strconv"
)

func unquoteString(s string) string {
	value, err := strconv.Unquote(s)
	if err != nil {
		fmt.Println("WARN: failed to unquote string: ", err)
		return s
	}
	return value
}

//func unquoteString(s string) string {
//	start := 1
//	end := len(s) - 1
//	if s[0] == '"' && s[end] == '"' {
//		return s[start:end]
//	}
//	return s
//}
