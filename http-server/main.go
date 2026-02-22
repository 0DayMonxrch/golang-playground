//Manual File Handling version
package main
import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func main(){
	http.HandleFunc("/",fileHandler)
	log.Println("Server running on port 3000")
	log.Fatal(http.ListenAndServe(":3000",nil))
}

func fileHandler(w http.ResponseWriter, r *http.Request){
	//Resolving file path
	path:=r.URL.Path
	if path == "/"{
		path="/index.html"
	}

	//preventing directory traversal attacks (Security Check)
	cleanPath := filepath.Clean("."+path)

	//read file
	data, err := os.ReadFile(cleanPath)
	if err!=nil{
		http.Error(w,"404 - File not found",http.StatusNotFound)
		return
	}

	// detecting MIME type
	ext := filepath.Ext(cleanPath)
	contentType := mime.TypeByExtension(ext)
	if contentType ==""{
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type",contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}