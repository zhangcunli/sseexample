package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"
)

type ResultStu struct {
	Username  string `json:"username,omitempty"`
	Time      string `json:"time,omitempty"`
	Text      string `json:"text,omitempty"`
	SteamDone bool   `json:"stream_done,omitempty"`
}

type CustomResult struct {
	Success   bool        `json:"success"`
	ErrorCode string      `json:"errorCode,omitempty"`
	ErrorMsg  string      `json:"errorMsg,omitempty"`
	Result    interface{} `json:"result,omitempty"`
}

func main() {
	errFlag := flag.Bool("errFlag", false, "return err")
	noResult := flag.Bool("noResult", false, "return err")
	flag.Parse()
	fmt.Printf(">>>main errFlag:%+v\n", *errFlag)

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Printf(">>>headers:%+v\n", r.Header)
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		// Set the necessary headers to allow for Server-Sent Events
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		startTime := time.Now().Unix()
		defer func() {
			endTime := time.Now().Unix()
			fmt.Printf("startTime:%d, endTime:%d, delta:%v\n", startTime, endTime, endTime-startTime)
		}()

		index := 0
		for {
			index++

			if *errFlag {
				if index == 5 {
					customResult := CustomResult{
						Success:   false,
						ErrorCode: "1000",
						ErrorMsg:  "failed",
					}
					customResults, _ := json.Marshal(customResult)
					resstr := fmt.Sprintf("data: %s\n", customResults)
					fmt.Fprintf(w, resstr)
					fmt.Printf("2.server write message, resstr:%s\n", resstr)
					break
				}
			}

			if *noResult {
				if index == 5 {
					customResult := CustomResult{
						Success: true,
					}
					customResults, _ := json.Marshal(customResult)
					resstr := fmt.Sprintf("data: %s\n", customResults)
					fmt.Fprintf(w, resstr)
					fmt.Printf("2.server write message, resstr:%s\n", resstr)
					//break
				}
			}

			if index > 5 {
				customResult := CustomResult{
					Success: true,
					Result: ResultStu{
						SteamDone: true,
					},
				}
				customResults, _ := json.Marshal(customResult)
				resstr := fmt.Sprintf("data: %s\n", customResults)
				fmt.Fprintf(w, resstr)
				fmt.Printf("3.server write message, resstr:%s\n", resstr)
				break
			}

			result := ResultStu{
				Username: "bobby",
				Time:     time.Now().Format("15:04:05"),
				Text:     "Request received",
			}

			customResult := CustomResult{
				Success: true,
				Result:  result,
			}
			customResults, _ := json.Marshal(customResult)
			resstr := fmt.Sprintf("data: %s\n", customResults)
			// Write some data to the client
			_, err := fmt.Fprintf(w, resstr)
			if err != nil {
				fmt.Printf("writing data err:%s, clientAddr:%s\n", err, r.RemoteAddr)
				return
			}

			// Flush the data immediately instead of buffering it for later.
			flusher.Flush()

			fmt.Printf("1.server write message, resstr:%s\n", resstr)

			// Pause for a second before the next iteration.
			time.Sleep(time.Duration(1) * time.Second)
			//time.Sleep(time.Duration(100) * time.Millisecond)
		}
	})

	httpServer := http.Server{
		Addr:         ":8081",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	httpServer.SetKeepAlivesEnabled(true)
	httpServer.ListenAndServe()
}
