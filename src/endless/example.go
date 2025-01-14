package main

import (
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Actual pid is %d", syscall.Getpid())
	w.Write([]byte("WORLD!"))
}

func main() {
	mux1 := mux.NewRouter()
	mux1.HandleFunc("/hello", handler).
		Methods("GET")
	log.Printf("Actual pid is %d", syscall.Getpid())
	err := endless.ListenAndServe("localhost:4242", mux1)
	if err != nil {
		log.Println(err)
	}

	log.Println("Server on 4242 stopped")

	os.Exit(0)
}

// kill -1 pid
// 平滑重启server  1 是SIGHUP信号
