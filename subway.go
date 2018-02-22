package main

import (
	"log"
	"net/http"

	"github.com/carillonator/subway/stationinfo"
)

var cis = stationinfo.ComplexInfoSet{}

func init() {
	var err error
	cis, err = stationinfo.NewComplexInfoSet()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	http.HandleFunc("/", requestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	/*
		stops := []int{237, 167, 611, 636}
		ss, err := NewStationSet(stops, cis)
		if err != nil {
			log.Fatal(err)
		}

		printText(ss, os.Stdout)
	*/
}
