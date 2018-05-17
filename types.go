package main

//ULiveEvent ... a unique LiveEvent with an uuid
type ULiveEvent struct {
	ID string`json:"ID"`
	LiveEvent
}

//LiveEvent struct
type LiveEvent struct {
	Titel string`json:"titel"`
	Presentor string`json:"presentor"`
	Description string`json:"description"`
    DateBegin string`json:"dateBegin"`
	DateEnd string`json:"dateEnd"`
	Live bool`json:"live"`
	Featured bool`json:"featured"`
}

//LiveEventID struct
type LiveEventID struct {
	ID string`json:"ID"`
}


