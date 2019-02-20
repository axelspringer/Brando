package main

//UniqueEvent ... An Event with uuid
type UniqueEvent struct {
	ID string `json:"ID"`
	Event
}

//Event ...
type Event struct {
	Title       string `json:"Title"`
	Presentor   string `json:"Presentor"`
	Description string `json:"Description"`
	StartDate   string `json:"StartDate"`
	EndDate     string `json:"EndDate"`
	Live        bool   `json:"Live"`
	Featured    bool   `json:"Featured"`
	Teaser      string `json:"Teaser"`
	Source      string `json:"Source"`
	Hidden      bool   `json:"Hidden"`
	Password    string `json:"Password,omitempty"`
}

//Msg ...
type Msg struct {
	Message string `json:"message"`
}
