package user

import (
	"time"
)

type User struct {
	ID            uint       `gorm:"primaryKey" json:"id,omitempty"`
	Email         string     `gorm:"column:email;not null;typevarchar(100);unique" json:"email,omitempty"`
	Password      string     `gorm:"column:password;not null;typevarchar(100)" json:"password,omitempty"`
	Names         string     `gorm:"column:names;not null;typevarchar(100)" json:"names,omitempty"`
	Birthdate     time.Time  `gorm:"column:birthdate;not null" json:"birthdate,omitempty"`
	Age           int        `gorm:"column:age;not null" json:"age,omitempty"`
	Gender        string     `gorm:"column:gender;not null;typevarchar(10)" json:"gender,omitempty"`
	FirstVaccine  string     `gorm:"column:first_vaccine;not null;typevarchar(100)" json:"firstVaccine,omitempty"`
	SecondVaccine string     `gorm:"column:second_vaccine;not null;typevarchar(100)" json:"secondVaccine,omitempty"`
	CreatedAt     *time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt,omitempty"`
}

type Login struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type Token struct {
	Token string `json:"token,omitempty"`
}

type Info struct {
	Age           int    `json:"age,omitempty"`
	Gender        string `json:"gender,omitempty"`
	FirstVaccine  string `json:"firstVaccine,omitempty"`
	SecondVaccine string `json:"secondVaccine,omitempty"`
}
