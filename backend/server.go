package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/config"
	"social-network/routes"
)


func main()  {
	mux := http.NewServeMux()

	// handler := &routes.Handler{
	// 	// Cntrlrs: ,
	// }

	routes.Routes(mux, handler)

	server := http.Server{
		Addr:    "0.0.0.0"+config.Port,
		Handler: mux,
	}

	fmt.Print("http://localhost"+config.Port)
	log.Fatal(server.ListenAndServe())
}