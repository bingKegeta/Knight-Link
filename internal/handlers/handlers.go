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
	"github.com/lib/pq"
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
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	UserName   string `json:"username"`
	Password   string `json:"password"`
	University string `json:"university"`
	Email      string `json:"email"`
	UserType   string `json:"user_type"`
	Uid        int
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
	Name string `json:"event_name"`
	// Tags           []string `json:"tags"`
	Description    string `json:"event_description"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Location       string `json:"loc_name"`
	Visibility     string `json:"visibility"`
	UniversityName string `json:"uni_name"`
	RsoName        string `json:"rso_name"`
	UniId          int    `json:"uni_id"`
	RsoId          int    `json:"rso_id"`
	LocId          int
}

type University struct {
	Name        string `json:"uni_name"`
	Description string `json:"uni_description"`
	StudentNo   int    `json:"student_no"`
}

type Location struct {
	Address   string `json:"address"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
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

	check := 0
	checkStdNo := `SELECT EXISTS (
		SELECT 1
		FROM students
		WHERE student_no = 0
		AND university_name = $1
	  ); `
	err = db.QueryRow(checkStdNo, user.University).Scan((&check))

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"Error":   "Error",
			"message": "Error checking university" + " " + err.Error(),
		})
		return
	}

	if check == 1 {
		// do the assignment of user_type here
	}

	// Every user is by default an student, so assign it here
	user.UserType = "admin"

	// Query the DB to get the Uid
	checkUid := `SELECT u.uni_id FROM public."Universities" u WHERE name = $1`
	err = db.QueryRow(checkUid, user.University).Scan(&user.Uid)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"Error":   "Error",
			"message": "Error checking university" + " " + err.Error(),
		})
		return
	}

	// Check if user exists before creating
	var userCount int
	checkUserQuery := `SELECT COUNT(*) FROM public."Users" WHERE username = $1`

	err = db.QueryRow(checkUserQuery, user.UserName).Scan(&userCount)

	switch {
	case err != nil:
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"Error":   "Error",
			"message": "Error checking for username" + err.Error(),
		})
		return

	case userCount > 0:
		render.Status(r, http.StatusInternalServerError)

		render.JSON(w, r, map[string]interface{}{
			"Error":   "Error",
			"message": "Username already taken",
		})
		return
	}

	query := `INSERT INTO public."Users" (first_name, last_name, username, "password",
        									uni_id,
											email,
											user_type)
											VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err = db.Exec(query, user.FirstName, user.LastName, user.UserName,
		user.Password, user.Uid, user.Email, user.UserType)

	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, err.Error())
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"status":  "Success",
		"message": "User Created",
	})
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
		render.JSON(w, r, map[string]interface{}{
			"status":  "Error",
			"message": "Error checking for username" + err.Error(),
		})
		return
	}

	// This is gonna be what's in the DB, to test against the info user sent
	var DbUser LoginForm
	err = db.QueryRow(`SELECT username, password FROM public."Users" WHERE username=$1`, userLogin.UserName).Scan(&DbUser.UserName, &DbUser.Password)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		if err.Error() == "sql: no rows in result set" {
			render.JSON(w, r, map[string]interface{}{
				"status":  "Error",
				"message": "Incorrect Credentials",
			})
			return
		}
		return
	}

	// Comparing the hashed from the DB to the user sent
	err = bcrypt.CompareHashAndPassword([]byte(DbUser.Password), []byte(userLogin.Password))

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "Error",
			"message": "Incorrect Credentials",
		})
		return
	}

	// If we got here means that the password is correct. We can create the token
	tokenString, err := createToken(tokenAuth, userLogin.UserName)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "Error",
			"message": "Catastrophic failure, try again",
		})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Sending the token
	render.JSON(w, r, map[string]interface{}{
		"status":  "success",
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
	db, err := connectToDB()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT e.name, e.tags, e.description, e.start_time, e.end_time, e.uni_id, e.rso_id, e.visibility FROM public."Events" e`)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error getting Events",
		})
		return
	}

	defer rows.Close()

	var events []EventForm

	for rows.Next() {
		var event EventForm
		err = rows.Scan(&event.Name, &event.Tags, &event.Description, &event.StartTime, &event.EndTime, &event.UniId, &event.RsoId, &event.Visibility)

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]interface{}{
				"status":  "warning",
				"message": "Error getting Events array",
			})
			return
		}
		events = append(events, event)
	}

	// If there was an error in the for, it should get here. But I think the
	// first return should honestly take care of it in any weird case..
	if err = rows.Err(); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error iterating over rows",
		})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"status": "success",
		"data":   events,
	})
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
			"status":  "warning",
			"message": "There was an error on the submitted form",
		})
		return
	}

	// In case there is RSO
	query := `SELECT r.rso_id from public."RSOs" r WHERE r.name = $1`

	// Get rso_id
	if event.RsoName != "" {
		// Get rso_id
		err = db.QueryRow(query, event.RsoName).Scan(&event.RsoId)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]interface{}{
				"status":  "warning",
				"message": "There was an error getting RSO ID",
			})
			return
		}
	}

	// Now the uni_id
	query = `SELECT u.uni_id FROM public."Universities" u WHERE u.name = $1`
	err = db.QueryRow(query, event.UniversityName).Scan(&event.UniId)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "There was an error getting University ID " + err.Error(),
		})
		return
	}

	// Last, loc_id
	query = `SELECT l.loc_id FROM public."Locations" l WHERE l.address = $1`

	err = db.QueryRow(query, event.Location).Scan(&event.LocId)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "There was an error getting Location ID " + err.Error(),
		})
		return
	}

	var insertQuery string
	var args []interface{}
	if event.RsoId != 0 {
		insertQuery = `INSERT INTO public."Events" (name, description, start_time, end_time, loc_id, uni_id, rso_id, visibility) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		args = append(args, event.Name, event.Description, event.StartTime, event.EndTime, event.LocId, event.UniId, event.RsoId, event.Visibility)
	} else {
		insertQuery = `INSERT INTO public."Events" (name, description, start_time, end_time, loc_id, uni_id, visibility) VALUES ($1, $2, $3, $4, $5, $6, $7)`
		args = append(args, event.Name, event.Description, event.StartTime, event.EndTime, event.LocId, event.UniId, event.Visibility)
	}

	_, err = db.Exec(insertQuery, args...)

	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok { // Check if it is a pq.Error
			switch pgerr.Code {
			case "23505": // Unique violation
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, map[string]interface{}{
					"status":  "error",
					"message": "An event with the same name already exists.",
				})
			case "P0001": // Custom error code from your PostgreSQL function
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, map[string]interface{}{
					"status":  "error",
					"message": pgerr.Message, // or a custom message based on your error handling
				})
			default:
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, map[string]interface{}{
					"status":  "error",
					"message": "Internal Server Error: " + pgerr.Message,
				})
			}
			return
		}
		// If it's not a pq.Error, handle it as an unknown error
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Unknown error occurred: " + err.Error(),
		})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"status":  "Success",
		"message": "Event Created",
	})
}

func GetAllRSOs(w http.ResponseWriter, r *http.Request) {
	db, err := connectToDB()

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT r.name FROM public."RSOs" r`)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error getting RSOs",
		})
		return
	}

	defer rows.Close()

	var rsos []string

	for rows.Next() {
		var rso string
		err = rows.Scan(&rso)

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]interface{}{
				"status":  "warning",
				"message": "Error getting Locations array",
			})
			return
		}
		rsos = append(rsos, rso)
	}

	// If there was an error in the for, it should get here. But I think the
	// first return should honestly take care of it in any weird case..
	if err = rows.Err(); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error iterating over rows",
		})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"status": "success",
		"data":   rsos,
	})
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

	db, err := connectToDB()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT u.name, u.description, u.student_no FROM public."Universities" u`)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error getting Universities",
		})
		return
	}

	defer rows.Close()

	var universities []University

	for rows.Next() {
		var uni University
		err = rows.Scan(&uni.Name, &uni.Description, &uni.StudentNo)

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]interface{}{
				"status":  "warning",
				"message": "Error getting Universities array",
			})
			return
		}
		universities = append(universities, uni)
	}

	// If there was an error in the for, it should get here. But I think the
	// first return should honestly take care of it in any weird case..
	if err = rows.Err(); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error iterating over rows",
		})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"status": "success",
		"data":   universities,
	})
}

func UpdateUniDetails(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to update the details of a specific university
	render.JSON(w, r, "UpdateUniDetails endpoint")
}

func GetAllLocations(w http.ResponseWriter, r *http.Request) {
	db, err := connectToDB()

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT l.address, l.latitude, l.longitude FROM public."Locations" l`)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error getting Locations",
		})
		return
	}

	defer rows.Close()

	var locations []Location

	for rows.Next() {
		var location Location
		err = rows.Scan(&location.Address, &location.Latitude, &location.Longitude)

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]interface{}{
				"status":  "warning",
				"message": "Error getting Locations array",
			})
			return
		}
		locations = append(locations, location)
	}

	// If there was an error in the for, it should get here. But I think the
	// first return should honestly take care of it in any weird case..
	if err = rows.Err(); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error iterating over rows",
		})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"status": "success",
		"data":   locations,
	})
}

func CreateLocation(w http.ResponseWriter, r *http.Request) {
	db, err := connectToDB()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	defer db.Close()

	var location Location

	err = json.NewDecoder(r.Body).Decode(&location)

	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}

	// Check if location exists before creating
	var locationCount int
	checkUserQuery := `SELECT COUNT(*) FROM public."Locations" WHERE address = $1`

	err = db.QueryRow(checkUserQuery, location.Address).Scan(&locationCount)

	switch {
	case err != nil:
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"Error":   "Error",
			"message": "Error checking for username" + err.Error(),
		})
		return

	case locationCount > 0:
		render.Status(r, http.StatusInternalServerError)

		render.JSON(w, r, map[string]interface{}{
			"Error":   "Error",
			"message": "Location already exists",
		})
		return
	}

	query := `INSERT INTO public."Locations" (address, latitude, longitude)
											VALUES ($1, $2, $3);`

	_, err = db.Exec(query, location.Address, location.Latitude, location.Longitude)

	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]interface{}{
			"Error":   "Error",
			"message": "Error creating event " + err.Error(),
		})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"status":  "Success",
		"message": "Event Created",
	})

}
func GetLocations(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic to get all the locations of a user in a given radius
	render.JSON(w, r, "GetLocations endpoint")
}
