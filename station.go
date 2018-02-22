package main

import (
	"strings"

	"github.com/carillonator/subway/stationinfo"
	"github.com/google/gtfs-realtime-bindings/golang/gtfs"
)

const stationFile string = "Stations.csv"

type Station struct {
	Id   int
	Name string
	// map of lines to arrival times
	NorthTimes map[string][]int64
	SouthTimes map[string][]int64
}

type StationSet struct {
	Timestamp uint64
	Stations  map[int]*Station
	Feeds     FeedSet
}

func NewStationSet(complexes []int, cis stationinfo.ComplexInfoSet) (*StationSet, error) {
	// TODO what happens if the set of stations changes? can this be prevented?
	// TODO probably load the cis here so it doesn't have to be done in main and passed all around here
	ss := StationSet{}
	ss.Stations = make(map[int]*Station)

	for _, s := range complexes {
		ss.Stations[s] = &Station{Id: s}
		ss.Stations[s].Name = cis[s].Name
		ss.Stations[s].NorthTimes = make(map[string][]int64)
		ss.Stations[s].SouthTimes = make(map[string][]int64)
	}

	feeds, err := NewFeedSetFromComplexes(complexes, cis)
	if err != nil {
		return nil, err
	}
	ss.Feeds = feeds

	err = ss.update()
	if err != nil {
		return nil, err
	}

	return &ss, nil
}

func (ss *StationSet) update() error {

	// TODO could this overwrite anything?
	for _, feed := range ss.Feeds {
		err := stationsFromFeed(ss, feed)
		if err != nil {
			return err
		}
	}

	return nil
}

func stationsFromFeed(ss *StationSet, feed *Feed) error {
	ts := feed.Timestamp
	// TODO this re-updates the ss with every feed, which is inaccurate
	ss.Timestamp = ts

	for _, tu := range feed.TripUpdates {
		routeId, times := stopTimesForTrip(tu)

		for stopId, time := range times { //gtfs ids
			complexId, dir := parseStopId(stopId)
			// only keep the stops in the StationSet
			if _, ok := ss.Stations[complexId]; ok {
				if dir == "N" {
					ss.Stations[complexId].NorthTimes[routeId] = append(ss.Stations[complexId].NorthTimes[routeId], time)
				} else {
					ss.Stations[complexId].SouthTimes[routeId] = append(ss.Stations[complexId].SouthTimes[routeId], time)
				}
			}
			// TODO if this is an update, remove stale ones??
		}
	}

	return nil
}

// return complexId and direction for given stop
func parseStopId(id string) (int, string) {
	var dir string
	if strings.HasSuffix(id, "N") {
		dir = "N"
	} else if strings.HasSuffix(id, "S") {
		dir = "S"
	} else {
		// TODO handle this better
		panic("unknown stopId: " + id)
	}

	complexId := stationinfo.ComplexFromGtfsId(id)

	return complexId, dir
}

func stopTimesForTrip(tu *gtfs.TripUpdate) (string, map[string]int64) {
	routeId := *tu.Trip.RouteId

	times := make(map[string]int64)
	for _, u := range tu.StopTimeUpdate {
		// someimes a stu doesn't have a departure field
		if u.Departure != nil {
			times[*u.StopId] = *u.Departure.Time
		}
	}
	return routeId, times
}
