package main

import (
	"fmt"
)

var mp map[string]string = make(map[string]string)

func main() {
	v, ok := mp["hello"]
	if !ok {
		mp["hello"] = "world"
	} else {
		fmt.Print(v)
	}

	for k, v := range mp {
		fmt.Println(k, v)
	}
}
