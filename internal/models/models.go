package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ParamsUser struct {
	UserID    string    `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Lastname  string    `json:"last_name" db:"last_name"`
	Password  string    `json:"password,omitempty"`
	Email     string    `json:"email"`
	Role      string    `json:"role" db:"role"`
	Verified  string    `json:"verified" db:"verified"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type DataSearchedUser struct {
	Total int           `json:"total"`
	Users []*ParamsUser `json:"users"`
}

func (p *ParamsUser) ClearPass() {
	p.Password = ""
}

type ParamsCreateAdmin struct {
	Id        string
	Name      string
	Lastname  string
	Password  string
	Email     string
	Role      string
	Verified  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrRoleUndefined = errors.New("role undefined")
)

func (p *ParamsCreateAdmin) Validate() error {

	p.HashPass()

	switch p.Role {
	case "admin":
	case "client":
	default:
		return ErrRoleUndefined
	}

	return nil
}

func (p *ParamsCreateAdmin) HashPass() {
	bc, _ := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	p.Password = string(bc)
}

type ParamsFind struct{}
type ParamsFindAll struct{}
type ParamsDeleteUser struct {
	UserId string
}
