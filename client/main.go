package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Result struct {
	Value float64 `json:"value"`
}

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
	responseBody, error := io.ReadAll(response.Body)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()
	var result Result
	error = json.Unmarshal(responseBody, &result)
	if error != nil {
		panic(error)
	}
	saveResultOnFile(result)
}

func saveResultOnFile(result Result) {
	file, error := os.Create("cotacao.txt")
	if error != nil {
		panic(error)
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("Dólar:%f", result.Value))
}
