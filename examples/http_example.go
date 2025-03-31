package main

import (
	"fmt"

	"github.com/pramithamj/microcomms/pkg/http"
)

func main() {
	client := http.NewClient("https://jsonplaceholder.typicode.com")
	resp, err := client.Get("/posts/1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Response:", resp.String())
}
