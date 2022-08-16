package model

type Info struct {
	Age           int
	Gender        int
	FirstVaccine  int
	SecondVaccine int
}

type IncomingInfo struct {
	Age           int
	Gender        int
	FirstVaccine  int
	SecondVaccine int
	Ticket        string
}

type OutgoingInfo struct {
	Age           int
	Gender        int
	FirstVaccine  int
	SecondVaccine int
	Prediction    int
	Ticket        string
}
