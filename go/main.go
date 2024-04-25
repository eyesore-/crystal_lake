package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Message struct {
	Text string `json:"text"`
}

func relativePath(s string) (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("Unable to get current file path")
	}

	path := filepath.Join(filepath.Dir(filename), "../data.json")

	return path, nil
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s %s - %s %s %v",
			r.UserAgent(),
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	path, err := relativePath("../data.json")
	if err != nil {
		fmt.Println("Error: Unable to get data filepath", err)
		return
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		http.Error(w, "Error reading JSON", http.StatusInternalServerError)
		return
	}

	w.Write(byteValue)
}

func main() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal("ERROR: Unable to retrieve network interfaces", err)
	}

	var localIP string

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
			}
		}
	}

	http.Handle("/", logMiddleware(http.HandlerFunc(jsonHandler)))

	port := ":8080"

	fmt.Printf("\n\n   Local:   http://localhost%s", port)
	fmt.Printf("\n   Network: http://%s:8080\n\n", localIP)
	log.Fatal(http.ListenAndServe(port, nil))
}
