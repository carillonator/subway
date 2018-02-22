package main

import (
	"net/http"
	"strconv"
	"strings"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	ss := strings.Split(q.Get("s"), ",")
	var si []int
	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			return
		}
		si = append(si, i)
	}

	times, err := NewStationSet(si, cis)
	if err != nil {
		return
	}

	printText(times, w)
}
