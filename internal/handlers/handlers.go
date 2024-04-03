package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// * Probably need to add more structs to make the JSONs easier to make(?)
type User struct {
	User_id         int    `json:"user_id"`
	Fname           string `json:"first_name"`
	Lname           string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	Authority       string `json:"auth"`
	RSO_affiliation bool   `json:"is_affiliated_with_rso"`
}

// Function to establish a connection to the database
func connectToDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, os.Getenv("PG_USER"), os.Getenv("PG_PW"), os.Getenv("PG_DB"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}

//! Remember to set the status codes

func GetUser(w http.ResponseWriter, r *http.Request) {
	db, err := connectToDB()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "500 Internal Server Error")
		return
	}
	defer db.Close()

	user_id := chi.URLParam(r, "userId")
	query, err := os.ReadFile("./SQL/api/user/GetById.sql")
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "500 Error Getting the Query")
		return
	}

	q := string(query)
	query_string := fmt.Sprintf(q, user_id)

	rows, err := db.Query(query_string)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, "404 Entry Not Found")
		return
	}
	defer rows.Close()

	var user User
	rows.Next()
	rows.Scan(&user.User_id, &user.Fname, &user.Lname, &user.Email, &user.Password, &user.Authority, &user.RSO_affiliation)

	render.Status(r, http.StatusFound)
	render.JSON(w, r, user)
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
