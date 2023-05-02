package main

import (
	"encoding/json"
	"fmt"
)

type Profile struct {
	Pk        string          `json:"pk"`
	Following bool            `json:"following"`
	Meta      ProfileMetadata `json:"meta"`
	Npub      string          `json:"npub"`
	Relays    []string        `json:"relays"`
}

func NewProfile() Profile {
	return Profile{
		Pk:        "",
		Npub:      "",
		Following: false,
		Meta:      ProfileMetadata{},
		Relays:    []string{},
	}
}

func NewProfileFromJson(j string) Profile {
	profile := NewProfile()
	err := json.Unmarshal([]byte(j), &profile)
	fmt.Errorf("NewProfileFromJson:", err.Error())
	return profile
}

func (p *Profile) ToJson() string {
	j, _ := json.Marshal(p)
	return string(j)
}
