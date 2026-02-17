package main

import (
	"context"
	"fmt"

	"github.com/sudankdk/offense/internal/tools/ffuf"
)

func main() {
	tool := &ffuf.Tool{}
	input := map[string]string{
		"url":      "https://sudankhadka.com.np/FUZZ",
		"wordlist": "./opt/ffuf/wordlist.txt",
	}
	res, err := tool.Run(context.Background(), input)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
