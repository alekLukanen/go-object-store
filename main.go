package main

import (
	"fmt"
	"net/http"
	"log"
	"io"
	"os"
	"errors"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
)

var PORT = 3000
var DATADIRECTORY = "./objectData"

func main() {
	fmt.Println("- Initializing the object storage")

	if !doesFileExist(DATADIRECTORY) {
		makeDirError := os.MkdirAll(DATADIRECTORY, os.ModePerm)
		if makeDirError != nil {
			panic(makeDirError)
		}
	}

	setupEndpoints()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}


func setupEndpoints() {
	r := mux.NewRouter()
	r.HandleFunc("/alive", alive)
	r.PathPrefix("/object/").Handler(objectHandler{})
	http.Handle("/", r)
}

func alive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Object Storage Alive")
}



type baseResponse struct {
    ObjectKey string
    FileName string
}

type objectHandler struct {}
func (obj objectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if (strings.Contains(r.URL.Path, " ") == true) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	objectKey := strings.Replace(r.URL.Path, "/object", "", 1)
	hashObjectKey := md5.Sum([]byte(objectKey))
	fileName := hex.EncodeToString(hashObjectKey[:])

	switch r.Method {
		case "GET": getObject(w, r, fileName)	
		case "PUT": putObject(w, r, objectKey, fileName)
		case "DELETE": deleteObject(w, r, fileName)
		default: {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}


func putObject(w http.ResponseWriter, r *http.Request, objectKey string, fileName string) {

	defer r.Body.Close()
	body, bodyReadError := io.ReadAll(r.Body)
	if bodyReadError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filePath := fmt.Sprintf("%s/%s", DATADIRECTORY, fileName)
	fileWriteError := os.WriteFile(filePath, body, 0644)
    if fileWriteError != nil {
        panic(fileWriteError)
    }

	response := &baseResponse{
		ObjectKey: objectKey,
        FileName: fileName,
	}
    resonseData, _ := json.Marshal(response)

	fmt.Fprintf(w, string(resonseData))

}

func getObject(w http.ResponseWriter, r *http.Request, fileName string) {

	filePath := fmt.Sprintf("%s/%s", DATADIRECTORY, fileName)

	if !doesFileExist(filePath) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, fileReadError := os.ReadFile(filePath)
	if fileReadError != nil {
		log.Fatal(fileReadError)
	}
	w.Write(data);

}

func deleteObject(w http.ResponseWriter, r *http.Request, fileName string) {

	filePath := fmt.Sprintf("%s/%s", DATADIRECTORY, fileName)
	if !doesFileExist(filePath) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fileRemoveError := os.Remove(filePath)
	if fileRemoveError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func doesFileExist(filePath string) bool {

	if _, err := os.Stat(filePath); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return false
	}

}
