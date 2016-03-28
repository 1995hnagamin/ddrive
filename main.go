package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"./id"
	"./node"
	"./pref"
)

var (
	preference *pref.Preference
)

func isBetween(hash string, p, q *node.Node) bool {
	return id.IsBetween(id.NewID([]byte(hash)), p.Id(), q.Id())
}

func isExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return os.IsExist(err)
}

func isResponsible(hash string) (bool, error) {
	self, err := preference.GetSelf()
	if err != nil {
		return false, err
	}
	successor, err := preference.GetSuccessor()
	if err != nil {
		return false, err
	}
	return isBetween(hash, self, successor), nil
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
	}
	is_responsible, err := isResponsible(hash)
	if err != nil {
		log.Printf("[get/error]%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error")
		return
	}
	if is_responsible {
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
	preference, err = pref.NewPreference("pref.db")
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
