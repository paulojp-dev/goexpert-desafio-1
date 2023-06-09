package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	request, error := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if error != nil {
		panic(error)
	}
	response, error := http.DefaultClient.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()
	io.Copy(os.Stdout, response.Body)
}
