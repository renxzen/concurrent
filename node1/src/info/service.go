package info

import (
	"errors"
	"log"
	"net/http"
	res "nodo1/src/response"
	"strings"
)

type infoService struct {
	NodeService NodeService
}

type InfoService interface {
	Predict(info Info) (res.Response, error)
}

func NewInfoService() InfoService {
	return &infoService{
		NodeService: NewNodeService(),
	}
}

func (is infoService) Predict(info Info) (res.Response, error) {
	outgoingInfo := OutgoingInfo{
		Age:           info.Age,
		Gender:        is.encodeGender(info.Gender),
		FirstVaccine:  is.encodeVaccine(info.FirstVaccine),
		SecondVaccine: is.encodeVaccine(info.SecondVaccine),
		Ticket:        RandomString(10),
	}

	if outgoingInfo.FirstVaccine == -1 || outgoingInfo.SecondVaccine == -1 || outgoingInfo.Gender == -1 {
		return res.Response{}, errors.New("invalid vaccine")
	}

	// send to node 2
	log.Println("Sending data of ticket:", outgoingInfo.Ticket)
	is.NodeService.PrepareOutgoingInfo(outgoingInfo)

	// get from node 3
	var incomingInfo IncomingInfo
	for {
		incomingInfo = is.NodeService.GetIncomingInfoFromChan()
		if incomingInfo.Ticket == outgoingInfo.Ticket {
			break
		}

		go is.NodeService.AddToIncomingInfoChan(incomingInfo)
	}

	log.Println("Got data of ticket:", incomingInfo.Ticket)

	// 0 - Alive, 1 - Dead
	code := "Alive"
	if incomingInfo.Prediction == 1 {
		code = "Dead"
	}

	response := res.NewResponse(http.StatusOK, code)
	return response, nil
}

func (is infoService) encodeVaccine(vaccine string) int {
	vaccines := []string{"ninguna", "sinopharm", "pfizer", "astrazeneca"}
	lowVaccine := strings.ToLower(vaccine)

	for i, v := range vaccines {
		if v == lowVaccine {
			return i
		}
	}

	return -1
}

func (is infoService) encodeGender(gender string) int {
	genders := []string{"female", "male"}
	lowGender := strings.ToLower(gender)

	for i, v := range genders {
		if v == lowGender {
			return i
		}
	}

	return -1
}
