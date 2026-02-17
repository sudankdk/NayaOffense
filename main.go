package main

import (
	"log"

	"github.com/sudankdk/offense/configs"
	_ "github.com/sudankdk/offense/internal/tools/ffuf"
	"github.com/sudankdk/offense/internal/transport"
	"github.com/sudankdk/offense/internal/transport/http/handler"
	"github.com/sudankdk/offense/internal/transport/http/router"
)

func main() {

	// load config
	cfg, err := configs.LoadConfig()

	//setup handlers
	h := handler.NewHandler()

	// setup routers

	r := router.NewRouter(h)

	// setup server\

	err = transport.StartServer(cfg.PORT, r)

	if err != nil {
		log.Fatalf("Error in starting the server: %v", err)
	}

	// tool := &ffuf.Tool{}
	// input := map[string]string{
	// 	"url":      "https://sudankhadka.com.np/FUZZ",
	// 	"wordlist": "./opt/ffuf/wordlist.txt",
	// }
	// res, err := tool.Run(context.Background(), input)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res)
}
