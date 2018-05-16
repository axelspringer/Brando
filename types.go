package main

//ULiveEvent ... a unique LiveEvent with an uuid
type ULiveEvent struct {
	ID string`json:"ID"`
	LiveEvent
}

//LiveEvent struct
type LiveEvent struct {
	titel string`json:"titel"`
	presentor string`json:"presentor"`
	description string`json:"description"`
    	dateBegin string`json:"dateBegin"`
	dateEnd string`json:"dateEnd"`
	live bool`json:"live"`
	featured bool`json:"featured"`
}

//LiveEventID struct
type LiveEventID struct {
	ID string`json:"ID"`
}


