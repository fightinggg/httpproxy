package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type HttpReq struct {
	Method string
	Url    string
	Header map[string][]string
	Body   string
}

type HttpRsp struct {
	Status int
	Header map[string][]string
	Body   string
}

func index(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}
	req := HttpReq{}
	err = json.Unmarshal(data, &req)
	if err != nil {
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	request, err := http.NewRequest(req.Method, req.Url, bytes.NewReader([]byte(req.Body)))

	for k, vs := range req.Header {
		for _, v := range vs {
			request.Header.Add(k, v)
		}
	}

	resp, err := http.DefaultClient.Do(request)

	w.WriteHeader(resp.StatusCode)
	for k, vs := range resp.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}
	_, _ = w.Write(data)

}

// CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build httpproxy.go
func main() {
	http.HandleFunc("/", index)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
