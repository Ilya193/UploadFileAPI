package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("upload")
	r.ParseMultipartForm(32 << 20)
	id := r.Form.Get("id")
	file, handler, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	f, err := os.OpenFile("upload-images/"+strconv.Itoa(randInt(1, 1000))+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	io.Copy(f, file)
	w.Write([]byte(id))
}

func images(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Open("D:\\code\\golang\\practice\\upload-images")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return
	}

	list := make([]string, len(files))

	for index, file := range files {
		list[index] = "http://192.168.1.7:8080/upload-images/" + file.Name()
	}

	jsonBytes, err := json.Marshal(list)
	w.Write(jsonBytes)
}

func main() {
	fs := http.FileServer(http.Dir("D:\\code\\golang\\practice\\upload-images"))
	http.HandleFunc("/upload", upload)
	http.Handle("/upload-images/", http.StripPrefix("/upload-images/", fs))
	http.HandleFunc("/images/", images)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
