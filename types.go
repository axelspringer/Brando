package main

//UniqueEvent ... An Event with uuid
type UniqueEvent struct {
	ID string `json:"ID"`
	Event
}

//Event ...
type Event struct {
	Titel       string `json:"Titel"`
	Presentor   string `json:"Presentor"`
	Description string `json:"Description"`
	StartDate   string `json:"StartDate"`
	EndDate     string `json:"EndDate"`
	Live        bool   `json:"Live"`
	Featured    bool   `json:"Featured"`
}

//Msg ...
type Msg struct {
	Message string `json:"message"`
}
