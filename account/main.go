package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting server ...")

	router := gin.Default()
	router.GET("/api/account", func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, gin.H{
			"hello": "suckaz",
		})
	})
	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server %v\n", err)
		}
	}()
	log.Printf("Listening on port %v\n", srv.Addr)

	// Wait for kill signal of channel
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// bloc until a signall
	<-quit
	// the contect is used to inform the server it has 5 second to finish
	// the request it is currently handling
	ctx, cancell := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancell()

	// shutdown server
	log.Println("Shutting down server ... ")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown: %v\n", err)
	}
}
