package main

import (
	"encoding/json"

	"io"

	"log"

	"net/http"

	"os"

	"strconv"

	"github.com/gorilla/mux"
)

type Image struct {
	ID int `json:"id"`

	Path string `json:"path"`
}

var images []Image

func UploadImage(w http.ResponseWriter, r *http.Request) {

	imagePath := "/home/username/go/src/imageuploader/"

	imageFile, header, err := r.FormFile("image")

	if err != nil {

		log.Fatal("Error retrieving the image from the form data: ", err)

	}

	defer imageFile.Close()

	fileName := header.Filename

	filePath := imagePath + fileName

	newImage := Image{ID: len(images) + 1, Path: filePath}

	images = append(images, newImage)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {

		log.Fatal("Error opening the file: ", err)

	}

	defer file.Close()

	io.Copy(file, imageFile)

	json.NewEncoder(w).Encode(newImage)

}

func GetImage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	for _, image := range images {

		if image.ID == id {

			http.ServeFile(w, r, image.Path)

			return

		}

	}

	json.NewEncoder(w).Encode("Image not found")

}

func DeleteImage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	for index, image := range images {

		if image.ID == id {

			err := os.Remove(image.Path)

			if err != nil {

				log.Fatal("Error deleting the file: ", err)

			}

			images = append(images[:index], images[index+1:]...)

			json.NewEncoder(w).Encode("Image successfully deleted")

			return

		}

	}

	json.NewEncoder(w).Encode("Image not found")

}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/upload", UploadImage).Methods("POST")

	router.HandleFunc("/get/{id}", GetImage).Methods("GET")

	router.HandleFunc("/delete/{id}", DeleteImage).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8081", router))

}
