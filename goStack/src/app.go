package src

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

var myStack []int

func push(w http.ResponseWriter, r *http.Request) {
	elem := r.URL.Query().Get("num")
	num, _ := strconv.Atoi(elem)
	myStack = append(myStack, num)
	fmt.Fprintf(w, "Stack Contents: %v after pushing %v", myStack, elem)
}

func pop(w http.ResponseWriter, r *http.Request) {
	elem := myStack[len(myStack)-1]
	myStack = myStack[:len(myStack)-1]
	fmt.Fprintf(w, "Stack Contents: %v after popping %v", myStack, elem)
}

func printStack(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Stack Contents: %v", myStack)
}

func Start() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", printStack).Methods(http.MethodGet)
	myRouter.HandleFunc("/push", push).Methods(http.MethodGet)
	myRouter.HandleFunc("/pop", pop).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      myRouter,
		Addr:         ":5151",
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	go func() {
		log.Println("Starting server at port 5151")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	intChan := make(chan os.Signal, 1)
	signal.Notify(intChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-intChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	os.Exit(0)
}
