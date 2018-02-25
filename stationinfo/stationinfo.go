package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// http://web.mta.info/developers/developer-data-terms.html#data
const stationsUrl = "http://web.mta.info/developers/data/nyct/subway/Stations.csv"
const complexUrl = "http://web.mta.info/developers/data/nyct/subway/StationComplexes.csv"

func main() {
	// TODO maybe write raw data to file also to track changes, since
	// the gob will change every time even with no csv changes
	data, err := csvFromUrl(stationsUrl)
	if err != nil {
		log.Fatal(err)
	}

	stations := infoByGtfsIds(data)
	complexes := complexInfoFromStationSet(stations)

	gob, err := makeGob(complexes)
	if err != nil {
		log.Fatal(err)
	}

	err = printCode(gob)
	if err != nil {
		log.Fatal(err)
	}
}

func printCode(gob []byte) error {
	hexGob := hex.EncodeToString(gob)
	fmt.Printf(`// DO NOT EDIT
// This file generated by the stationinfo package

package main

const cisGobHex = "%s"
`, hexGob)

	return nil
}

func makeGob(object interface{}) ([]byte, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(object)
	if err != nil {
		return nil, err
	}

	gob := buff.Bytes()

	return gob, nil
}

// It's easier to work with only complexes
// even if a given complex ID only has one station
// Some "complexes" are two lines that stop at the same platform
//   which are separate lines in Stations.csv, but with identical station and complex ids
func complexInfoFromStationSet(ss StationInfoSet) ComplexInfoSet {
	// TODO call infoByGtfsId() from here

	cs := ComplexInfoSet{}

	for _, s := range ss { // iterate all stations
		complexId := s.ComplexId

		var ci *ComplexInfo
		// initialize the complex if this is the first
		if _, ok := cs[complexId]; !ok {
			ci = &ComplexInfo{}

			// for now just use the first station found's name for the complex
			// TODO read StationComplexes.csv
			// Complex ID,Complex Name
			ci.Name = s.Name
			ci.Id = s.ComplexId
		} else {
			ci = cs[complexId]
		}

		ci.Stations = append(ci.Stations, s)

		cs[complexId] = ci
	}

	return cs
}

// GTFS Stop ID is the only unique station identifier
// e.g. 167 is duplicate for both Station and Complex IDs, both IND
// Not all multi-stations like 167 are in Complexes.csv
func infoByGtfsIds(data [][]string) StationInfoSet {

	set := StationInfoSet{}

	for i, line := range data {
		// skip the header
		if i == 0 {
			continue
		}

		s := StationInfo{}

		// Station ID,Complex ID,GTFS Stop ID,Division,Line,Stop Name,Borough,Daytime Routes,Structure,GTFS Latitude,GTFS Longitude
		s.Id, _ = strconv.Atoi(line[0])
		s.ComplexId, _ = strconv.Atoi(line[1])
		s.GtfsId = line[2]
		s.Division = line[3]
		s.Line = line[4]
		s.Name = line[5]
		s.Borough = line[6]
		s.DayRoutes = strings.Split(line[7], " ")
		s.Structure = line[8]
		s.Lat, _ = strconv.ParseFloat(line[9], 64)
		s.Long, _ = strconv.ParseFloat(line[10], 64)

		feeds := make(map[uint8]bool)
		for _, r := range s.DayRoutes {
			feeds[feedByRoute[r]] = true
		}
		for f, _ := range feeds {
			s.Feeds = append(s.Feeds, f)
		}

		set[s.GtfsId] = &s
	}
	return set
}

func csvFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// http://datamine.mta.info/list-of-feeds
var feedByRoute = map[string]uint8{
	"1": 1,
	"2": 1,
	"3": 1,
	"4": 1,
	"5": 1,
	"6": 1,
	"7": 51,
	"A": 26,
	"C": 26,
	"E": 26,
	"N": 16,
	"Q": 16,
	"R": 16,
	"W": 16,
	"B": 21,
	"D": 21,
	"F": 21,
	"M": 21,
	"L": 2,
	"G": 31,
	"J": 36,
	"Z": 36,
}
