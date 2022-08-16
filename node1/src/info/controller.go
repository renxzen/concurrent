package info

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	res "nodo1/src/response"
)

type infoController struct {
	InfoService InfoService
}

type InfoController interface {
	IncomingRequest(http.ResponseWriter, *http.Request)
}

func NewRequestController() InfoController {
	return &infoController{
		InfoService: NewInfoService(),
	}
}

func EnableInfoController(handler *http.ServeMux) {
	rc := NewRequestController()

	handler.HandleFunc("/predict", res.POST(rc.IncomingRequest))
}

func (rc infoController) IncomingRequest(w http.ResponseWriter, r *http.Request) {
	info := Info{}
	reqBody, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(reqBody, &info); err != nil {
		response := res.NewResponse(http.StatusBadRequest, err.Error())
		res.SendResponse(w, response)
		return
	}

	response, err := rc.InfoService.Predict(info)
	if err != nil {
		response := res.NewResponse(http.StatusBadRequest, err.Error())
		res.SendResponse(w, response)
		return
	}

	res.SendData(w, http.StatusOK, response)
}
