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

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

func style(c string, s string) string {
	return fmt.Sprintf("%s%s%s", c, s, Reset)
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
			style(Magenta, r.UserAgent()),
			r.RemoteAddr,
			style(Cyan, r.Method),
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
	localhost := "http://localhost" + port
	localNetwork := "http://" + localIP + ":" + port

	fmt.Printf("\n\n   %s   %s", style(Bold, "Local:"), style(Blue, localhost))
	fmt.Printf("\n   %s %s\n\n", style(Bold, "Network:"), style(Blue, localNetwork))
	log.Fatal(http.ListenAndServe(port, nil))
}
