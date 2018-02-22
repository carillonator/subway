package stationinfo

import (
	"encoding/gob"
	"os"
)

func saveGob(path string, object interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(object)
	if err != nil {
		return err
	}

	return nil
}

func loadGob(path string, object interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(object)
	if err != nil {
		return err
	}

	return nil
}
