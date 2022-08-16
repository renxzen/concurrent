package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/mail"
	"strings"
	"time"

	res "backend/src/response"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	UserRepository UserRepository
}

type UserService interface {
	Register(User) (User, error)
	LoginUser(Login) (Token, error)
	FindByUsername(string) (User, error)
	GetPrediction(Info) (res.Response, error)
}

func NewUserService() UserService {
	return &userService{
		UserRepository: NewUserRepository(),
	}
}

func (us *userService) Register(user User) (User, error) {
	user.Email = strings.ToLower(user.Email)

	// Check if email is a valid address
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return User{}, errors.New("invalid email")
	}

	// Check if email is already registered
	if _, err := us.UserRepository.GetByEmail(user.Email); err == nil {
		return User{}, errors.New("email already in use")
	}

	// Correct the casing of the whole name
	user.Names = strings.Title(strings.ToLower(user.Names))

	// Get age from birthdate
	user.Age = int(time.Since(user.Birthdate).Hours() / 8760)

	// Encrypt password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return User{}, errors.New("encrypting failed")
	}

	user.Password = string(bytes)

	// Register user in the database
	user, err = us.UserRepository.Create(user)
	if err != nil {
		return User{}, errors.New("user registration failed")
	}

	return User{}, nil
}

func (us *userService) LoginUser(login Login) (Token, error) {
	// Get user from database
	user, err := us.UserRepository.GetByEmail(login.Email)
	if err != nil {
		return Token{}, errors.New("email not found")
	}

	// Convert to byte array the login and user hashed password
	userPassword := []byte(user.Password)
	loginPassword := []byte(login.Password)

	// Encrypt the login password and compare
	err = bcrypt.CompareHashAndPassword(userPassword, loginPassword)
	if err != nil {
		return Token{}, errors.New("rmail/password don't match")
	}

	// Get a JWT from Provider
	token, err := JWT.GetToken(user)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func (us *userService) FindByUsername(email string) (User, error) {
	// Search user in the database
	user, err := us.UserRepository.GetByEmail(email)
	if err != nil {
		return User{}, errors.New("email not found")
	}

	// Exclude values from the response
	user.ID = 0
	user.Password = ""

	return user, nil
}

func (us *userService) GetPrediction(info Info) (res.Response, error) {
	nodo1Host := "http://localhost:8081"

	infoBytes, _ := json.Marshal(info)
	infoJson := bytes.NewBuffer(infoBytes)

	resp, err := http.Post(nodo1Host+"/predict", "application/json", infoJson)
	if err != nil {
		return res.Response{}, errors.New("could not make request")
	}

	response := res.Response{}

	reqBody, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(reqBody, &response); err != nil {
		return res.Response{}, err
	}

	return response, nil
}
