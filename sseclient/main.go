package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
)

func main() {
	port := flag.Int("port", 8081, "remote port")
	flag.Parse()

	url := fmt.Sprintf("http://localhost:%d/events", *port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return
	}
	req.Header.Set("Connection", "Keep-Alive")

	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 2,
	}
	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	bufReader := bufio.NewReader(resp.Body)
	for {
		rawLine, readErr := bufReader.ReadBytes('\n')
		if readErr != nil {
			fmt.Printf("bufReader.ReadBytes() error: %v\n", readErr)
			break
		}

		if string(rawLine) == "[DONE]" {
			break
		}
		fmt.Printf("|||response: %s", string(rawLine))
		fmt.Printf("|||statusCode: %v\n\n", resp.StatusCode)
	}

	bodys, _ := io.ReadAll(resp.Body)
	fmt.Printf("StatusCode:%v, bodys:%s\n", resp.StatusCode, bodys)
	//fmt.Printf("body: %v\n", string(body))
}
