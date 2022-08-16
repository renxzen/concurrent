package user

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	res "backend/src/response"
)

type userController struct {
	UserService UserService
}

type UserController interface {
	Default(http.ResponseWriter, *http.Request)
	Registration(http.ResponseWriter, *http.Request)
	LoginAttempt(http.ResponseWriter, *http.Request)
	FindUsername(http.ResponseWriter, *http.Request)
	GetPrediction(http.ResponseWriter, *http.Request)
}

func NewUserController() UserController {
	return &userController{
		UserService: NewUserService(),
	}
}

func EnableUserController(handler *http.ServeMux) {
	uc := NewUserController()

	handler.HandleFunc("/", uc.Default)

	handler.HandleFunc("/register", res.POST(uc.Registration))
	handler.HandleFunc("/login", res.GET(uc.LoginAttempt))

	handler.HandleFunc("/find", res.GET(CheckJWT(uc.FindUsername)))
	handler.HandleFunc("/predict", res.POST(uc.GetPrediction))
}

func (uc *userController) Default(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status  string `json:"status"`
		Version string `json:"version"`
	}{
		Status:  "OK",
		Version: "0.6.9",
	}

	res.SendData(w, http.StatusOK, response)
}

func (uc *userController) Registration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := User{}
	reqBody, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(reqBody, &user); err != nil {
		response := res.NewResponse(http.StatusBadRequest, err.Error())
		res.SendResponse(w, response)
		return
	}

	_, err := uc.UserService.Register(user)
	if err != nil {
		log.Println("Error in registration:", user.Email)
		response := res.NewResponse(http.StatusBadRequest, err.Error())
		res.SendResponse(w, response)
		return
	}

	log.Println("User registered: ", user.Email)
	response := res.NewResponse(http.StatusOK, "user registered")
	res.SendResponse(w, response)
}

func (uc *userController) LoginAttempt(w http.ResponseWriter, r *http.Request) {
	login := Login{}
	reqBody, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(reqBody, &login); err != nil {
		response := res.NewResponse(http.StatusBadRequest, err.Error())
		res.SendResponse(w, response)
		return
	}

	token, err := uc.UserService.LoginUser(login)
	if err != nil {
		log.Println("Error in login:", login.Email)
		response := res.NewResponse(http.StatusBadRequest, err.Error())
		res.SendResponse(w, response)
		return
	}

	log.Println("User logged in:", login.Email)
	res.SendData(w, http.StatusOK, token)
}

func (uc *userController) FindUsername(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		response := res.NewResponse(http.StatusBadRequest, "email required")
		res.SendResponse(w, response)
		return
	}

	user, err := uc.UserService.FindByUsername(email)
	if err != nil {
		log.Println("Error in finding user:", user.Email)
		response := res.NewResponse(http.StatusBadRequest, err.Error())
		res.SendResponse(w, response)
		return
	}

	log.Println("User found:", user.Email)
	res.SendData(w, http.StatusOK, user)
}

func (uc *userController) GetPrediction(w http.ResponseWriter, r *http.Request) {
	info := Info{}
	reqBody, _ := ioutil.ReadAll(r.Body)

	if err := json.Unmarshal(reqBody, &info); err != nil {
		response := res.NewResponse(http.StatusBadRequest, err.Error())
		res.SendResponse(w, response)
		return
	}

	prediction, err := uc.UserService.GetPrediction(info)
	if err != nil {
		log.Println("Error in prediction:", info)
		response := res.NewResponse(http.StatusBadRequest, err.Error())
		res.SendResponse(w, response)
		return
	}

	log.Println("Prediction sent:", prediction.Message)
	res.SendData(w, http.StatusOK, prediction)
}
