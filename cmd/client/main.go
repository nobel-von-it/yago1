package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// Client for testing my shortener
func main() {
	url := "http://localhost:8080/api/shorten"
	fmt.Println("URL:>", url)

	jsonStr := []byte(`{"url":"http://youtube.com"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
