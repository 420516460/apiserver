package main

import (
	"apiserver/config"
	"apiserver/router"
	"errors"
	_ "fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path")
) //TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	pflag.Parse()

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// Set gin mode
	gin.SetMode(viper.GetString("runmode"))

	// Create the Gin engine.
	g := gin.New()

	// gin middlewares
	middlewares := []gin.HandlerFunc{}

	// Routes, set Http Header
	router.Load(
		// Cores
		g,

		// Middlewares
		middlewares...,
	)

	// Ping the server to make sure the router is working
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up", err)
		}
		log.Print("The router has been deployed successfully")
	}()

	log.Printf("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Printf(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to '/health'
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Print("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}
