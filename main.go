package main

import (
	"apiserver/router"
	"errors"
	_ "fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	// Create the Gin engine.
	g := gin.New()

	// gin middlewares
	middlewares := []gin.HandlerFunc{}

	// Routes.
	router.Load(
		// Cores.
		g,

		// Middlewares.
		middlewares...,
	)

	log.Printf("Start to listening the incoming requests on http address: %s", ":8080")
	log.Printf(http.ListenAndServe(":8080", g).Error())

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up", err)
		}
		log.Print("The router has been deployed successfully")
	}()
}

func pingServer() error {
	for i := 0; i < 10; i++ {
		// Ping the server by sending a GET request to '/health'
		resp, err := http.Get("http://localhost:8080/ping" + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a sec to continue the next ping.
		log.Print("Waiting for the server to start")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}
