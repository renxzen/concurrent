package info

type IncomingInfo struct {
	Age           int
	Gender        int
	FirstVaccine  int
	SecondVaccine int
	Prediction    int
	Ticket        string
}

type OutgoingInfo struct {
	Age           int
	Gender        string
	FirstVaccine  string
	SecondVaccine string
	Prediction    int
	Ticket        string
}
