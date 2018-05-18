package main

//UniqueEvent ... An Event with uuid
type UniqueEvent struct {
	ID string `json:"ID"`
	Event
}

//Event ...
type Event struct {
	Titel       string `json:"titel"`
	Presentor   string `json:"presentor"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Live        bool   `json:"live"`
	Featured    bool   `json:"featured"`
}

//Msg ...
type Msg struct {
	Message string `json:"message"`
}
