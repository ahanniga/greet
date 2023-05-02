package main

type ProfileMetadata struct {
	Name        string `json:"name,omitempty"`
	About       string `json:"about,omitempty"`
	Picture     string `json:"picture,omitempty"`
	NIP05       string `json:"nip05,omitempty"`
	DisplayName string `json:"display_name"`
	Lud06       string `json:"lud06"`
	Banner      string `json:"banner"`
	Website     string `json:"website"`
}
