package models

type APIResponse struct {
	Kind string  `json:"kind"`
	Data APIData `json:"data"`
}

type APIData struct {
	After    string        `json:"after"`
	Children []APIChildren `json:"children"`
}

type APIChildren struct {
	Kind string  `json:"kind"`
	Data Listing `json:"data"`
}

type Listing struct {
	Name           string `json:"name"`
	Ups            uint   `json:"ups"`
	AuthorFullName string `json:"author_fullname"`
	Author         string `json:"author"`
	Title          string `json:"title"`
}
