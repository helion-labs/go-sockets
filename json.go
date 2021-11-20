package main

import (
		"encoding/json"
		"log"
		)

func get_json(time, distance int) []byte{
	l := Location{time, distance}

	data, err := json.Marshal(&l)
	if err != nil {
		log.Fatal(err)
	}
	return data 
}
