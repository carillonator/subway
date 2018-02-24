package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
)

func startServer() {
	r := http.NewServeMux()

	r.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(requestHandler)))

	http.ListenAndServe(":8080", handlers.CompressHandler(r))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ss := strings.Split(q.Get("s"), ",")
	var si []int
	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			http.Error(w, "Invalid Complex ID", 400)
			fmt.Println(err)
			return
		}
		si = append(si, i)
	}

	times, err := NewStationSet(si, cis)
	if err != nil {
		if strings.Contains(err.Error(), "not an existing complex id") {
			http.Error(w, "Invalid Complex ID", 400)
		} else {
			http.Error(w, "Server Error", 500)
		}
		fmt.Fprintln(os.Stderr, err)
		return
	}

	printText(times, w)
}
