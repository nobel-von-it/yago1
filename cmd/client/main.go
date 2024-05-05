package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Client for testing my shortener
func main() {
	endpoint := "http://localhost:8080/"
	data := url.Values{}

	fmt.Println("enter a long URL")
	reader := bufio.NewReader(os.Stdin)
	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	long = strings.TrimSuffix(long, "\n")
	data.Set("address", long)

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}

	cli := &http.Client{}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := cli.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println("Code: ", res.StatusCode)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
