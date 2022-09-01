package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"aggregator/middleware"
	"aggregator/providers/infura"
	"aggregator/services"

	"github.com/gorilla/mux"
)

var ()

func main() {
	var wait time.Duration
	var host, port, endpoint, key, version string

	// Override default environment settings
	flag.StringVar(&host, "listening-ip", os.Getenv("HOST"), "self listening ip")
	flag.StringVar(&port, "listening-port", os.Getenv("PORT"), "self listening port")
	flag.StringVar(&endpoint, "service-endpoint", os.Getenv("INFURA_ENDPOINT"), "endpoint provided by infura service")
	flag.StringVar(&key, "service-key", os.Getenv("INFURA_KEY"), "key provided by infura service")
	flag.StringVar(&version, "api-version", os.Getenv("API_VERSION"), "the version exposed in api uri eg. /v1/route")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Default mux router
	router := mux.NewRouter()
	// Initialize provider RPC service
	provider, err := infura.New(endpoint + "/" + key)
	if err != nil {
		log.Fatal(err)
	}

	// Build subrouter for API version.
	log.Printf("running api version: %s", version)
	v1 := router.PathPrefix("/" + version).Subrouter()
	v1.Use(mux.CORSMethodMiddleware(v1))
	v1.Use(middleware.Headers())

	// Routes works as service group switch.
	routes := services.New(v1, provider)
	// Functional routes works as service breaker.
	routes.
		EnableBlocks().      // Turn on blocks services
		EnableTransactions() // Turn on transactions services
	// In this case Send transactions service is not allowed.
	// Lets imagine that we have an issue with Send features and we need to turn it off temporary.
	// 404 will be received by API client.
	// EnableSendTransactions()

	// Setup Graceful Shutdown
	// ref: https://pkg.go.dev/net/http#Server
	log.Printf("listening on %s:%s", host, port)
	srv := &http.Server{
		Addr:         host + ":" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		// reuse connection
		// ref: https://www.sobyte.net/post/2022-03/go-http-keep-alive/
		IdleTimeout: time.Second * 60,
		Handler:     v1,
	}

	go func() {
		// Run our server in a goroutine so that it doesn't block.
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM will be caught if server run inside docker.
	signal.Notify(c, os.Kill, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
