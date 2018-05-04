package main

//LiveEvent Struct
type LiveEvent struct {
	ID string`json:"ID"`
	Titel string`json:"Titel"`
	Presentor string`json:"Presentor"`
	Description string`json:"Description"`
    DateBegin string`json:"DateBegin"`
	DateEnd string`json:"DateEnd"`
	Live bool`json:"Live"`
	Featured bool`json:"Featured"`
}

//Error Struct 
type Error struct {
	Msg string
	Err string
}