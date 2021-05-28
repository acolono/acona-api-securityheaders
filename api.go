package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
)

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	url, err := getParameter(r.URL, "url")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidURL(url) {
		err := fmt.Errorf("Invalid url: '%v'", url)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := exec.Command("python", "securityheaders.py", "--formatter", "json", "--", url)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out.Bytes())
}

func getParameter(u *url.URL, name string) (string, error) {
	values, ok := u.Query()[name]
	if !ok || len(values[0]) < 1 {
		return "", fmt.Errorf("Parameter '%v' not found", name)
	}
	return values[0], nil
}

func isValidURL(url string) bool {
	re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)
	if re.MatchString(url) {
		return true
	}
	return false
}
