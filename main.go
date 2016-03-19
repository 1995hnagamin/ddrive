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

var (
	preference *pref.Preference
)

func isExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return os.IsExist(err)
}

func isResponsible(hash string) bool {
	return true
}

func getSuccessorAddress(preference *pref.Preference) (string, error) {
	successor, err := preference.GetSuccessor()
	if err != nil {
		return "", err
	}
	return successor.String(), nil
}

func getFilepath(hash string) string {
	return fmt.Sprintf("objects/%s/%s", hash[0:2], hash[2:])
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[len("/get/"):]
	filepath := getFilepath(hash)
	if isExist(filepath) {
		log.Printf("[get/serve]%s", hash)
		http.ServeFile(w, r, filepath)
		return
	} else if isResponsible(hash) {
		log.Printf("[get/not found]%s", filepath)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found")
		return
	}
	address, err := getSuccessorAddress(preference)
	if err != nil {
		log.Printf("[get/error]%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error")
		return
	}
	redirectURL := fmt.Sprintf(
		"http://%s/get/%s",
		address,
		hash)
	log.Printf("[get/redirect]%s", hash)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, http.StatusText(http.StatusBadRequest))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("[upload/error]%s", err.Error())
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
	var err error
	preference, err = pref.NewPreference("database")
	defer preference.Close()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/get/", getHandler)
	http.HandleFunc("/upload/", uploadHandler)
	err = http.ListenAndServe(":8180", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
}
