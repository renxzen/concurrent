package info

type OutgoingInfo struct {
	Age           int
	Gender        int
	FirstVaccine  int
	SecondVaccine int
	Ticket        string
}

type Info struct {
	Age           int
	Gender        string
	FirstVaccine  string
	SecondVaccine string
}

type IncomingInfo struct {
	Age           int
	Gender        string
	FirstVaccine  string
	SecondVaccine string
	Prediction    int
	Ticket        string
}
