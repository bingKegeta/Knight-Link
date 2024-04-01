package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

type User struct {
	UID         int    `json:"uid"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	DateCreated string `json:"dateCreated"`
	Desc        string `json:"desc"`
	Authority   string `json:"authority"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to get a user
	render.JSON(w, r, "GetUser endpoint")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to delete a user
	render.JSON(w, r, "DeleteUser endpoint")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to update a user
	render.JSON(w, r, "UpdateUser endpoint")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to create a user
	render.JSON(w, r, "CreateUser endpoint")
}

func Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to login a user
	render.JSON(w, r, "Login endpoint")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to logout a user
	render.JSON(w, r, "Logout endpoint")
}

func GetAllEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to get all events
	render.JSON(w, r, "GetAllEvents endpoint")
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to get an event
	render.JSON(w, r, "GetEvent endpoint")
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to delete an event
	render.JSON(w, r, "DeleteEvent endpoint")
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to update an event
	render.JSON(w, r, "UpdateEvent endpoint")
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to create an event
	render.JSON(w, r, "CreateEvent endpoint")
}

func GetAllRSOs(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to get all RSOs
	render.JSON(w, r, "GetAllRSOs endpoint")
}

func GetRSO(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to get an RSO
	render.JSON(w, r, "GetRSO endpoint")
}

func DeleteRSO(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to delete an RSO
	render.JSON(w, r, "DeleteRSO endpoint")
}

func UpdateRSO(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to update an RSO
	render.JSON(w, r, "UpdateRSO endpoint")
}

func CreateRSO(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to create an RSO
	render.JSON(w, r, "CreateRSO endpoint")
}
