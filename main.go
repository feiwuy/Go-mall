package main

import (
	"context"
	"gomall/config"
	"gomall/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.InitConfig()
	config.InitDB()

	// 注册一个Router
	r := routers.SetUpRouter()

	port := config.AppConfig.App.Port
	if port == "" {
		log.Fatal("Failed reading port, the server will run in http://localhost:8080")
		port = "8080"
	}

	// if err := r.Run(port); err != nil {
	// 	log.Fatalf("Failed to run this program, error: %v", err)
	// }

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown, %v", err)
	}
	log.Println("Server exiting")
}
