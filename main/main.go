package main

import (
	"URL_Shortener"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	mux := defaultMux()

	filePtr := flag.String("filename", "source", "a string")
	filePath := "C:\\GolandProjects\\URL_Shortener\\" + *filePtr + ".yml"
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := URL_Shortener.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback

	yamlHandler, err := URL_Shortener.YAMLHandler(file, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8000")
	err = http.ListenAndServe(":8000", yamlHandler)
	if err != nil {
		panic(err)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "Hello, world!")
	if err != nil {
		return
	}
}
