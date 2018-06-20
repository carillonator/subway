package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	gtfs "github.com/carillonator/subway/gtfs-realtime"
	_ "github.com/carillonator/subway/nyct-subway"
	"github.com/golang/protobuf/proto"
)

const urlBase = "http://datamine.mta.info/mta_esi.php"

type Feed struct {
	Id          uint8
	Timestamp   uint64
	TripUpdates []*gtfs.TripUpdate
	Vehicles    []*gtfs.VehiclePosition
	Alerts      []*gtfs.Alert
}

func NewFeed(id uint8) (*Feed, error) {
	feed, err := fetch(id)
	if err != nil {
		fmt.Sprintf("Failed to fetch feed %d: %s", id, err)
		return nil, err
	}

	tu, veh, al := getEntities(feed)

	return &Feed{
		Id:          id,
		Timestamp:   *feed.Header.Timestamp,
		TripUpdates: tu,
		Vehicles:    veh,
		Alerts:      al,
	}, nil
}

// TODO make interface
type FeedSet map[uint8]*Feed

func NewFeedSetFromComplexes(complexes []int, cis ComplexInfoSet) (FeedSet, error) {
	ids := make(map[uint8]bool)
	for _, c := range complexes {
		for _, s := range cis[c].Stations {
			for _, id := range s.Feeds {
				ids[id] = true
			}
		}
	}

	var feeds []uint8
	for id, _ := range ids {
		feeds = append(feeds, id)
	}

	fs, err := NewFeedSetFromIds(feeds)
	if err != nil {
		return nil, err
	}

	return fs, nil
}

func NewFeedSetFromIds(ids []uint8) (FeedSet, error) {

	set := make(map[uint8]*Feed)
	// TODO make parallel
	for _, id := range ids {
		feed, err := NewFeed(id)
		if err != nil {
			return nil, err
		}
		set[id] = feed
	}

	return set, nil
}

func (f *Feed) Refresh() error {
	return nil
}

func getEntities(feed *gtfs.FeedMessage) ([]*gtfs.TripUpdate, []*gtfs.VehiclePosition, []*gtfs.Alert) {
	var (
		tu  []*gtfs.TripUpdate
		veh []*gtfs.VehiclePosition
		al  []*gtfs.Alert
	)

	for _, entity := range feed.Entity {

		if entity.Vehicle != nil {
			veh = append(veh, entity.Vehicle)
		} else if entity.Alert != nil {
			al = append(al, entity.Alert)
		} else if entity.TripUpdate != nil {
			tu = append(tu, entity.TripUpdate)
		} else {
			// unknown entity type
		}
	}

	return tu, veh, al
}

func fetch(id uint8) (*gtfs.FeedMessage, error) {
	url := fmt.Sprintf("%s?key=%s&feed_id=%d", urlBase, mtaKey, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	feed := gtfs.FeedMessage{}
	err = proto.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	// TODO better handle flaky feeds
	return &feed, nil
}
