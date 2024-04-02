package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

// * Probably need to add more structs to make the JSONs easier to make(?)
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

//! Remember to set the status codes

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

func AttendEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to add user to event list
	render.JSON(w, r, "AttendEvent endpoint")
}

func UnattendEvent(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to remove user from event list
	render.JSON(w, r, "UnattendEvent endpoint")
}

func CreateFeedback(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to create comment and rating from a user to an event
	render.JSON(w, r, "CreateFeedback endpoint")
}

func JoinRSO(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to allow user to join an RSO
	render.JSON(w, r, "JoinRSO endpoint")
}

func LeaveRSO(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to allow user to leave an RSO
	render.JSON(w, r, "LeaveRSO endpoint")
}

func GetAllUnis(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to get all Universities
	render.JSON(w, r, "GetAllUnis endpoint")
}

func GetUni(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to get a specific university by id(?)
	render.JSON(w, r, "GetUni endpoint")
}

func DeleteUni(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to delete a University
	render.JSON(w, r, "DeleteUni endpoint")
}

func UpdateUniDetails(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to update the details of a specific university
	render.JSON(w, r, "UpdateUniDetails endpoint")
}

func CreateUni(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to create a university
	render.JSON(w, r, "CreateUni endpoint")
}

func JoinUni(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to allow a user to join a university
	render.JSON(w, r, "JoinUni endpoint")
}

func LeaveUni(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to allow a user to leave a university
	render.JSON(w, r, "LeaveUni endpoint")
}

func GetLocations(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to get all the locations of a user in a given radius
	render.JSON(w, r, "GetLocations endpoint")
}
