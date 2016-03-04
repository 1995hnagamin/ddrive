package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"./pref"
)

func isExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return err != nil
}

func isResponsible(hash string) bool {
	return true
}

func getSuccessorAddress() (string, error) {
	self, err := pref.GetSelf()
	if err != nil {
		return "", err
	}
	return self.String(), nil
}

func getFilepath(hash string) string {
	return fmt.Sprintf("objects/%s/%s", hash[0:2], hash[2:])
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[len("/get/"):]
	filepath := getFilepath(hash)
	if isExist(filepath) {
		http.ServeFile(w, r, filepath)
		return
	} else if isResponsible(hash) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	address, err := getSuccessorAddress()
	if err != nil {
		return
	}
	redirectURL := fmt.Sprintf(
		"http://%s/get/%s",
		address,
		hash)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, http.StatusText(http.StatusBadRequest))
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hash := sha256.Sum256(body)
	path := getFilepath(string(hash[:]))
	if isExist(path) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Already uploaded")
		return
	}
}

func main() {
	http.HandleFunc("/get/", getHandler)
	http.HandleFunc("/upload/", uploadHandler)
	err := http.ListenAndServe(":8180", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
}
