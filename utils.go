package main

import (
	"encoding/json"
	"github.com/nbd-wtf/go-nostr"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func containsEvent(events []*nostr.Event, id string) bool {
	for _, v := range events {
		if v.ID == id {
			return true
		}
	}
	return false
}

func getContentMeta(event *nostr.Event) (*ProfileMetadata, error) {
	var metadata *ProfileMetadata
	err := json.Unmarshal([]byte(event.Content), &metadata)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}
func setContentMeta(meta *ProfileMetadata) (*nostr.Event, error) {
	var ev = nostr.Event{}
	m, err := json.Marshal(&meta)
	if err != nil {
		return nil, err
	}
	ev.Content = string(m)
	return &ev, nil
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for {
		if len(slice) == 0 {
			break
		}

		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
