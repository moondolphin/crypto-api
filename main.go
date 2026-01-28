// @title Crypto API
// @version 1.0
// @description API de cotizaciones de criptomonedas
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer " followed by a space and the JWT token.
package main

import (
	"log"
	"os"

	nethttp "net/http"

	_ "github.com/moondolphin/crypto-api/docs"

	"github.com/moondolphin/crypto-api/bootstrap"
	"github.com/moondolphin/crypto-api/config"
)

func main() {
	r, err := bootstrap.Start()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	port := config.Getenv("HTTP_PORT", "8080")
	log.Fatal(nethttp.ListenAndServe(":"+port, r))
}
