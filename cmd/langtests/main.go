package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
)

func main() {
	// data содержит данные в формате gob

	// напишите код, который декодирует data в массив строк
	// 1. создайте буфер `bytes.NewBuffer(data)` для передачи в декодер
	// 2. создайте декодер `dec := gob.NewDecoder(buf)`
	// 3. определите `make([]string, 0)` для получения декодированного слайса
	// 4. декодируйте данные используя функцию `dec.Decode`
}
