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
	Description    string        `json:"event_description"`
	StartTime      string        `json:"start_time"`
	EndTime        string        `json:"end_time"`
	Location       string        `json:"loc_name"`
	Visibility     string        `json:"visibility"`
	UniversityName string        `json:"uni_name"`
	RsoName        string        `json:"rso_name"`
	UniId          int           `json:"uni_id"`
	RsoId          sql.NullInt32 `json:"rso_id"`
	LocId          sql.NullInt32
}

type UserEventForm struct {
	Name        string         `json:"event_name"`
	Description sql.NullString `json:"event_description"`
	StartTime   string         `json:"start_time"`
	EndTime     string         `json:"end_time"`
	UniId       int            `json:"uni_id"`
}

type DbEventForm struct {
	Name string `json:"event_name"`
	// Tags           []string `json:"tags"`
	Description    sql.NullString `json:"event_description"`
	StartTime      string         `json:"start_time"`
	EndTime        string         `json:"end_time"`
	Location       string         `json:"loc_name"`
	Visibility     string         `json:"visibility"`
	UniversityName string         `json:"uni_name"`
	RsoName        string         `json:"rso_name"`
	UniId          int            `json:"uni_id"`
	RsoId          sql.NullInt32  `json:"rso_id"`
	LocId          sql.NullInt32
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

type FeedbackForm struct {
	Username  string `json:"username"`
	Eventname string `json:"event_name"`
	Type      string `json:"type"`
	Feedback  string `json:"feedback"`
}

type FeedbackReturn struct {
	Username  string `json:"username"`
	Eventname string `json:"event_name"`
	Type      string `json:"type"`
	Feedback  string `json:"feedback"`
	Timestamp string `json:"timestamp"`
}

type RsoForm struct {
	Name          string `json:"rso_name"`
	Description   string `json:"description"`
	AdminName     string `json:"am_name"`
	PromotionUser string `json:"promotion_value"`
	Sone          string `json:"s1_name"`
	Stwo          string `json:"s2_name"`
	Sthree        string `json:"s3_name"`
	Sfour         string `json:"s4_name"`

	AdminId      int
	DateCreated  string `json:"date_created"`
	RsoId        int
	UniId        int
	StudentEmail string
}

type RsoJoin struct {
	Username string `json:"username"`
	RsoName  string `json:"rso_name"`
}

type EventJoin struct {
	Username  string `json:"username"`
	Eventname string `json:"event_name"`
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
	// Actually.
	checkStdNo := `SELECT u.student_no FROM public."Universities" u WHERE name = $1`

	err = db.QueryRow(checkStdNo, user.University).Scan((&check))

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"Error":   "Error",
			"message": "Error checking university" + " " + err.Error(),
		})
		return
	}

	if check == 0 {
		user.UserType = "admin"
	} else {
		// Every user is by default an student, so assign it here
		user.UserType = "student"
	}

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

func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	db, err := connectToDB()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT u.username FROM public."Users" u`)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error getting Students",
		})
		return
	}

	defer rows.Close()
	var students []string

	for rows.Next() {
		var user string
		err = rows.Scan(&user)

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]interface{}{
				"status":  "warning",
				"message": "Error getting Locations array",
			})
			return
		}
		students = append(students, user)
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
		"data":   students,
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

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    userLogin.UserName,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		Secure:   true,                    // Only send over HTTPS
		SameSite: http.SameSiteStrictMode, // Prevent CSRF attacks
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
	username := r.URL.Query().Get("username")
	if username != "" {
		GetUserEvents(w, r, username)
		return
	}
	db, err := connectToDB()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT e.name, e.description, e.start_time, e.end_time, e.uni_id, e.rso_id, e.visibility FROM public."Events" e`)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error getting Events",
		})
		return
	}

	defer rows.Close()

	var events []DbEventForm

	for rows.Next() {
		var event DbEventForm
		err = rows.Scan(&event.Name, &event.Description, &event.StartTime, &event.EndTime, &event.UniId, &event.RsoId, &event.Visibility)

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

func GetUserEvents(w http.ResponseWriter, r *http.Request, username string) {
	db, err := connectToDB()

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "There was an error connecting to the DB",
		})
		return
	}

	// Get the current user based on the token (in this case username cookie).
	// var eventJoin EventJoin

	// err = json.NewDecoder(r.Body).Decode(&eventJoin)

	// if err != nil {
	// 	render.Status(r, http.StatusInternalServerError)
	// 	render.JSON(w, r, map[string]interface{}{
	// 		"status":  "warning",
	// 		"message": "There was an error parsing the data",
	// 	})
	// 	return
	// }

	// get user_id
	var userId int
	query := `SELECT user_id FROM public."Users" WHERE username = $1`
	err = db.QueryRow(query, username).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]interface{}{
				// This should never ever happen. But added just as precaution
				"status":  "warning",
				"message": "User not found",
			})
			return
		} else {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]interface{}{
				"status":  "error",
				"message": "Database error: " + err.Error(),
			})
			return
		}
	}

	query = `select e."name", e."description", e."start_time", e."end_time", e."uni_id" 
			from public."Events" e 
			left join user_event_membership uem 
			on 
			e.event_id = uem.event_id where uem.user_id = $1`

	rows, err := db.Query(query, userId)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error getting Events.",
		})
		return
	}

	defer rows.Close()

	var events []UserEventForm

	for rows.Next() {
		var event UserEventForm
		err = rows.Scan(&event.Name, &event.Description, &event.StartTime, &event.EndTime, &event.UniId)

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
	if event.RsoId.Int32 != 0 {
		insertQuery = `INSERT INTO public."Events" (name, description, start_time, end_time, loc_id, uni_id, rso_id, visibility) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		args = append(args, event.Name, event.Description, event.StartTime, event.EndTime, event.LocId, event.UniId, event.RsoId, event.Visibility)
	} else {
		insertQuery = `INSERT INTO public."Events" (name, description, start_time, end_time, loc_id, uni_id, visibility) VALUES ($1, $2, $3, $4, $5, $6, $7)`
		args = append(args, event.Name, event.Description, event.StartTime, event.EndTime, event.LocId, event.UniId, event.Visibility)
	}

	_, err = db.Exec(insertQuery, args...)

	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			switch pgerr.Code {
			case "23505": // Unique violation
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, map[string]interface{}{
					"status":  "error",
					"message": "An event with the same name already exists.",
				})
			case "P0001":
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, map[string]interface{}{
					"status":  "error",
					"message": pgerr.Message,
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

	rows, err := db.Query(`SELECT r.name, r.description, r.date_created FROM public."RSOs" r`)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error getting RSOs",
		})
		return
	}

	defer rows.Close()

	var rsos []RsoForm

	for rows.Next() {
		var rso RsoForm
		err = rows.Scan(&rso.Name, &rso.Description, &rso.DateCreated)

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
	db, err := connectToDB()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}
	defer db.Close()

	decoder := json.NewDecoder(r.Body)
	var rsoForm RsoForm
	err = decoder.Decode(&rsoForm)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "Error",
			"message": "Error parsing form data " + err.Error(),
		})
		return
	}

	// Get the user details based on promotion value
	usernames := map[string]string{"1": rsoForm.Sone, "2": rsoForm.Stwo, "3": rsoForm.Sthree, "4": rsoForm.Sfour}
	username, exists := usernames[rsoForm.PromotionUser]
	if !exists {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]interface{}{
			"status":  "Error",
			"message": "Invalid promotion user selection",
		})
		return
	}

	query := `SELECT user_id, uni_id, email FROM public."Users" WHERE username = $1`
	err = db.QueryRow(query, username).Scan(&rsoForm.AdminId, &rsoForm.UniId, &rsoForm.StudentEmail)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "Error",
			"message": "Error fetching user details " + err.Error(),
		})
		return
	}

	// Create RSO record
	query = `INSERT INTO public."RSOs" (name, description, uni_id, admin_id, date_created)
			 VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) RETURNING rso_id;`
	err = db.QueryRow(query, rsoForm.Name, rsoForm.Description, rsoForm.UniId, rsoForm.AdminId).Scan(&rsoForm.RsoId)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "Error",
			"message": "Failed to create RSO " + err.Error(),
		})
		return
	}

	// Add all users to the membership table
	usernamesToInsert := []string{rsoForm.Sone, rsoForm.Stwo, rsoForm.Sthree, rsoForm.Sfour}
	for _, uname := range usernamesToInsert {
		query = `INSERT INTO public."User_RSO_Membership" (user_id, rso_id)
				 SELECT user_id, $1 FROM public."Users" WHERE username = $2;`
		_, err = db.Exec(query, rsoForm.RsoId, uname)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, "Error adding user to RSO: "+err.Error())
			return
		}
	}

	render.JSON(w, r, map[string]interface{}{
		"status":  "Success",
		"message": "RSO created successfully and users added.",
	})
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
	db, err := connectToDB()

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "There was an error connecting to the DB " + err.Error(),
		})
		return
	}

	defer db.Close()

	var feedback FeedbackForm

	err = json.NewDecoder(r.Body).Decode(&feedback)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "There was an error on the submitted comment",
		})
		return
	}

	render.JSON(w, r, "CreateFeedback endpoint")
	var user_id int
	query := `SELECT user_id FROM "Users" WHERE username = $1`
	err = db.QueryRow(query, feedback.Username).Scan(&user_id)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": ("There is an error with the username" + err.Error()),
		})
		return
	}

	var event_id int
	query = `SELECT event_id FROM "Events" WHERE name = $1`
	err = db.QueryRow(query, feedback.Eventname).Scan(&event_id)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "There is an error with the event name",
		})
		return
	}

	if feedback.Type == "comment" {
		query := `INSERT INTO public."Event_Feedback" (user_id, event_id, "content", "feedback_type", "timestamp") 
				VALUES ($1, $2, $3, 'comment', CURRENT_TIMESTAMP)`

		_, err := db.Exec(query, user_id, event_id, feedback.Feedback)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, err.Error())
			return
		}

		render.JSON(w, r, map[string]interface{}{
			"status":  "Success",
			"message": "Comment Added",
		})
	} else {
		query := `INSERT INTO "Event_Feedback" (user_id, event_id, "rating", "feedback_type", "timestamp")
					VALUES ($1, $2, $3, 'rating', CURRENT_TIMESTAMP)`
		_, err := db.Exec(query, user_id, event_id, feedback.Feedback)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, err.Error())
			return
		}

		render.JSON(w, r, map[string]interface{}{
			"status":  "Success",
			"message": "Rating Added",
		})
	}

	// TODO: Implement the logic to create comment and rating from a user to an event
	// render.JSON(w, r, "CreateFeedback endpoint")
}

func GetFeedback(w http.ResponseWriter, r *http.Request) {
	db, err := connectToDB()

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "There was an error connecting to the DB " + err.Error(),
		})
		return
	}

	defer db.Close()

	event_name := r.Header.Get("event_name")
	var event_id int

	err = db.QueryRow(`SELECT event_id FROM "Events" WHERE name = $1;`, event_name).Scan(&event_id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error getting Event Id",
		})
		return
	}

	query := `SELECT user_id, feedback_type, content, rating, timestamp FROM "Event_Feedback" WHERE event_id = $1;`
	rows, err := db.Query(query, event_id)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error getting Feedback",
		})
		return
	}

	var feedback []FeedbackReturn

	for rows.Next() {
		var user_id int
		var fb_type string
		var comment sql.NullString
		var rating sql.NullString
		var timestamp string
		var username string

		err = rows.Scan(&user_id, &fb_type, &comment, &rating, &timestamp)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]interface{}{
				"status":  "warning",
				"message": "Error getting Feedback array",
			})
			return
		}

		err = db.QueryRow(`SELECT username FROM "Users" WHERE user_id = $1;`, user_id).Scan(&username)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]interface{}{
				"status":  "warning",
				"message": "Error getting username",
			})
			return
		}

		var fb FeedbackReturn
		fb.Username = username
		fb.Eventname = event_name
		fb.Type = fb_type
		fb.Timestamp = timestamp
		if fb_type == "rating" {
			fb.Feedback = rating.String
		} else {
			fb.Feedback = comment.String
		}

		feedback = append(feedback, fb)
	}

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
		"data":   feedback,
	})
}

func JoinRSO(w http.ResponseWriter, r *http.Request) {
	db, err := connectToDB()

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error connecting to Database",
		})
		return
	}
	defer db.Close()

	var rsoJoin RsoJoin

	err = json.NewDecoder(r.Body).Decode(&rsoJoin)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "There was an error parsing the data",
		})
		return
	}

	// get rso_id
	var rso_id string
	query := `SELECT rso_id FROM public."RSOs" WHERE name = $1`
	err = db.QueryRow(query, rsoJoin.RsoName).Scan(&rso_id)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "There was an error getting the RSO ID",
		})
		return
	}

	var userId int

	// Get the userid off the username
	query = `SELECT user_id FROM public."Users" WHERE username = $1`
	err = db.QueryRow(query, rsoJoin.Username).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]interface{}{
				// This should never ever happen. But added just as precaution
				"status":  "warning",
				"message": "User not found",
			})
			return
		} else {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]interface{}{
				"status":  "error",
				"message": "Database error: " + err.Error(),
			})
			return
		}
	}

	// Check if the user is already a member of the RSO
	var exists bool
	query = `SELECT EXISTS(SELECT 1 FROM public."User_RSO_Membership" WHERE user_id = $1 AND rso_id = $2)`
	err = db.QueryRow(query, userId, rso_id).Scan(&exists)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "error",
			"message": "Database error: " + err.Error(),
		})
		return
	}

	if exists {
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "You are already in this RSO",
		})
		return
	}

	// If not a member, insert the user into the RSO membership table
	query = `INSERT INTO public."User_RSO_Membership" (user_id, rso_id) VALUES ($1, $2)`
	_, err = db.Exec(query, userId, rso_id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"status":  "warning",
			"message": "Error adding user to RSO: " + err.Error(),
		})
		return
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, map[string]interface{}{
		"status": "success",
		"data":   "User added to RSO group",
	})
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

// func GetAllComments(w http.ResponseWriter, r *http.Request) {
// 	db, err := connectToDB()
// 	if err != nil {
// 		render.Status(r, http.StatusInternalServerError)
// 		render.PlainText(w, r, err.Error())
// 		return
// 	}
// 	defer db.Close()

// 	rows, err := db.Query(`SELECT u.name, u.description, u.student_no FROM public."Universities" u`)

// }

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
