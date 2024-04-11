package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// * Probably need to add more structs to make the JSONs easier to make(?)
type User struct {
	UserID         int    `json:"user_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	UserName       string `json:"username"`
	Email          string `json:"email"`
	Auth           string `json:"auth"`
	RSOAffiliation bool   `json:"is_affiliated_with_rso"`
}

type LoginForm struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserNoId struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Uid       int    `json:"uni_id"`
	Email     string `json:"email"`
	UserType  string `json:"user_type"`
}

type SQLParseToType struct {
	query        string
	CustomStruct interface{}
}

type AuthMessage struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

type EventForm struct {
	Name        string `json:"event_name"`
	Description string `json:"event_description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Location    int    `json:"loc_id"`
	Visibility  string `json:"visibility"`
	UniId       int    `json:"uni_id"`
	RsoId       int    `json:"rso_id"`
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

// Function to parse SQL script into string
func parseSQL(fp string, class interface{}) (string, error) {
	query, err := os.ReadFile(fp)
	if err != nil {
		return "", err
	}

	q := string(query)

	var result = SQLParseToType{query: q, CustomStruct: class}

	valueOfStruct := reflect.ValueOf(result.CustomStruct)
	var values []interface{}

	for i := 0; i < valueOfStruct.NumField(); i++ {
		fieldValue := valueOfStruct.Field(i)

		values = append(values, fieldValue.Interface())
	}

	query_string := fmt.Sprintf(result.query, values...)

	return query_string, nil
}

//! Remember to set the status codes

func GetUser(w http.ResponseWriter, r *http.Request) {
	db, err := connectToDB()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	defer db.Close()

	user_id := chi.URLParam(r, "userId")

	uid := struct {
		User_id string `json:"user_id"`
	}{
		User_id: user_id,
	}

	query_string, err := parseSQL("./SQL/api/user/GetById.sql", uid)

	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, err.Error())
		return
	}

	rows, err := db.Query(query_string)

	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, err.Error())
		return
	}

	defer rows.Close()

	var scannedUser User

	for rows.Next() {
		err = rows.Scan(&scannedUser.UserID, &scannedUser.FirstName, &scannedUser.LastName, &scannedUser.UserName, &scannedUser.Email)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			render.PlainText(w, r, "Error scanning row")
			return
		}
		fmt.Println("User Details:", scannedUser)
	}

	render.Status(r, http.StatusFound)
	render.JSON(w, r, scannedUser)
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
	db, err := connectToDB()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	defer db.Close()

	var user UserNoId
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		// Handle error during hashing
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}

	// Set the password to the newly hashed password
	user.Password = string(hashedPassword)

	// Every user is by default an student, so assign it here
	user.UserType = "student"
	// This should be selected in the creation logic or something. This selects
	// the school
	user.Uid = 1

	query, err := parseSQL("./SQL/api/user/CreateUser.sql", user)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}

	_, err = db.Exec(query)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, err.Error())
		return
	}

	render.JSON(w, r, "User Created")
}

func Login(tokenAuth *jwtauth.JWTAuth, w http.ResponseWriter, r *http.Request) {

	db, err := connectToDB()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	defer db.Close()

	decoder := json.NewDecoder(r.Body)
	var userLogin LoginForm

	err = decoder.Decode(&userLogin)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}

	// This is gonna be what's in the DB, to test against the info user sent
	var DbUser LoginForm
	err = db.QueryRow(`SELECT username, password FROM public."Users" WHERE username=$1`, userLogin.UserName).Scan(&DbUser.UserName, &DbUser.Password)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, "The username is incorrect")
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}

	// Comparing the hashed from the DB to the user sent
	err = bcrypt.CompareHashAndPassword([]byte(DbUser.Password), []byte(userLogin.Password))

	if err != nil {
		render.JSON(w, r, DbUser)
		render.JSON(w, r, "The password is incorrect")
		return
	}

	// If we got here means that the password is correct. We can create the token
	tokenString, err := createToken(tokenAuth, userLogin.UserName)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "Failed to generate token")
		return
	}

	// Sending the token
	render.JSON(w, r, map[string]interface{}{
		"token":   tokenString,
		"message": "Login successful",
	})
}

func createToken(tokenAuth *jwtauth.JWTAuth, username string) (string, error) {
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	if err != nil {
		return "", err
	}
	return tokenString, nil
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

// Auth token required...
func CreateEvent(w http.ResponseWriter, r *http.Request) {

	db, err := connectToDB()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	defer db.Close()

	//
	var event EventForm

	err = json.NewDecoder(r.Body).Decode(&event)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"err":     err.Error(),
			"message": "There was an error on the submitted form",
		})
		return
	}

	// No need to add uni_id or rso_id, since that can get added later or not at all

	query := `INSERT INTO public."Events" (name, event_description, start_time, end_time, loc_id, visibility) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = db.Exec(query, event.Name, event.Description)

	render.JSON(w, r, event)
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
