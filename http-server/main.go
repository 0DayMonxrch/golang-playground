//using built-in file server
package main

import (
	"log"
	"net/http"
)

func main(){
	fileServer := http.FileServer(http.Dir("./"))
	http.Handle("/", fileServer)
	log.Println("Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000",nil))
}