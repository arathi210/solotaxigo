package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"

	"golang.org/x/crypto/bcrypt"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

type Person struct {
	Loc  string `json:"loc"`
	Date string `json:"date"`
}
type User struct {
	Id       int
	Name     string
	Email    string
	Mobile   string
	Password string
}
type Fare struct {
	Id       int
	FareName string
	BaseFare string
	MinFare  string
	Cost     string
	UserId   int
}
type RideHistory struct {
	Id int

	FromDate    string
	FromTime    string
	ToTime      string
	FromAddress string
	ToAddress   string
	Duration    string
	Distance    string
	Total       string
	BaseFare    string
	Cost        string
	Waiting     string
	UserSession int
}
type RideHistoryDetail struct {
	Id int

	FromDate    string
	FromTime    string
	ToTime      string
	FromAddress string
	ToAddress   string
	Duration    string
	Distance    string
	Total       string
	BaseFare    string
	Cost        string
	Waiting     string
	UserSession int
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// func dbConn() (db *sql.DB) {
// 	dbDriver := "mysql"
// 	dbUser := "solotaxi-go-test"
// 	dbPass := "kyOHnErbeKNpl9W"
// 	dbName := "solotaxi_go_test"
// 	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	//log.Println("connected")
// 	return db

// }

func dbConn() (db *sql.DB) {
	dbDriver := goDotEnvVariable("DB_DRIVER")
	dbUser := goDotEnvVariable("DB_USERNAME")
	dbPass := goDotEnvVariable("DB_PASSWORD")
	dbName := goDotEnvVariable("DB_DATABASE")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("godotenv : %s = %s \n", "STRONGEST_AVENGER", dbDriver)
	//log.Println("connected")
	return db

}

var tmpl = template.Must(template.ParseGlob("form/*"))

func LoginPage(w http.ResponseWriter, r *http.Request) {

	res := 0

	tmpl.ExecuteTemplate(w, "Login", res)

}
func Logincheck(res http.ResponseWriter, req *http.Request) {

	db := dbConn()
	//log.Println(db)
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

	email := req.FormValue("email")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT email, password FROM users WHERE email=?", email).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		//fmt.Println(err)
		//dialog.Alert("Username Incorrect")
		message := "Username Incorrect"
		http.Redirect(res, req, "/error?error_message="+message+"&page=login", 301)
		//http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		message := "Password Incorrect"
		http.Redirect(res, req, "/error?error_message="+message+"&page=login", 301)
		// dialog.Alert("Password incorrect")
		// http.Redirect(res, req, "/login", 301)
		return
	}

	session, _ := store.Get(req, "cookie-name")

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true

	selDB, err := db.Query("SELECT id, email  FROM users  WHERE email=?", email)
	if err != nil {
		panic(err.Error())
	}
	user := User{}
	for selDB.Next() {
		var id int
		var email string
		err = selDB.Scan(&id, &email)
		if err != nil {
			panic(err.Error())
		}
		user.Id = id
		user.Name = email
		session.Values["username"] = email
		session.Values["id"] = id

	}

	session.Save(req, res)

	defer db.Close()
	message := "Login Successful"
	http.Redirect(res, req, "/success?success_message="+message+"&page=startmap", 301)
	//res.Write([]byte("Hello" + databaseUsername))
	//http.Redirect(res, req, "/startmap", 301)

}
func Logout(w http.ResponseWriter, r *http.Request) {

	//res := 0
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Values["id"] = ""
	session.Save(r, w)

	http.Redirect(w, r, "/login", 301)

}
func Register(w http.ResponseWriter, r *http.Request) {

	res := 0
	tmpl.ExecuteTemplate(w, "Register", res)

}
func RegisterUser(res http.ResponseWriter, req *http.Request) {
	//session, _ := store.Get(req, "cookie-name")
	db := dbConn()
	if req.Method != "POST" {
		//http.ServeFile(res, req, "signup.html")
		tmpl.ExecuteTemplate(res, "Register", res)
		return
	}
	name := req.FormValue("name")
	mobile := req.FormValue("mobile")
	email := req.FormValue("email")

	password := req.FormValue("password")

	var databaseUsername string

	err := db.QueryRow("SELECT email FROM users WHERE email=?", email).Scan(&databaseUsername)
	fmt.Println(err)
	if err == nil {
		fmt.Println("Email already exists ! Plese register with different email")
		//dialog.Alert("Email Exists")
		message := "Email already exists ! Plese register with different email"
		http.Redirect(res, req, "/error?error_message="+message+"&page=register", 301)

		//tmpl.ExecuteTemplate(res, "Emailerror", message)
		return
	}

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			res.Write([]byte("Server error!"))
			//session.Values["error_message"] = "Server error, unable to create your account."
			http.Error(res, "Server error, unable to create your account.", 500)
			// message := beeep.Notify("Server error", "unable to create your account", "assets/information.png")
			// if message != nil {
			// 	panic(message)
			// }
			return
		}

		_, err = db.Exec("INSERT INTO users(name,mobile,email, password) VALUES(?, ?, ?, ?)", name, mobile, email, hashedPassword)
		if err != nil {
			res.Write([]byte("Server error!"))
			//session.Values["error_message"] = "Server error, unable to create your account."
			http.Error(res, "Server error, unable to create your account.", 500)
			// message := beeep.Notify("Server error", "unable to create your account", "assets/information.png")
			// if message != nil {
			// 	panic(message)
			// }
			return
		}

		// notify.Notify("app name", "alert", "some text", "public/asset/images/add.png")
		// Swal.fire("Any fool can use a computer")
		defer db.Close()
		//res.Write([]byte("User created!"))
		// message := beeep.Notify("User", "User created!", "assets/information.png")
		// if message != nil {
		// 	panic(message)
		// }
		// c echo.Context;
		// main.Set(c, "message", "Password is correct, you have been authenticated!")
		//res.Write([]byte("User Registered successfully!" + name))
		//dialog.Alert("User created!")
		//session.Values["success_message"] = "User created"

		// data := make(map[string]interface{})
		// data["error_message"] = session.Values["error_message"]
		// data["success_message1"] = "coming"

		// res = append(res, data)
		thank_message := "Registered Successfully! Thank you !"
		http.Redirect(res, req, "/success?success_message="+thank_message+"&page=login", 301)
		// http.Redirect(res, req, "/thank-you", 301)
		//tmpl.ExecuteTemplate(res, "Register", data)
		return

	case err != nil:
		res.Write([]byte("Server error!"))
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(res, req, "/register", 301)

	}

}
func Success(w http.ResponseWriter, r *http.Request) {

	res := 0
	tmpl.ExecuteTemplate(w, "Success", res)

}
func Error(w http.ResponseWriter, r *http.Request) {

	res := 0
	tmpl.ExecuteTemplate(w, "Error", res)

}

func Forgotpassword(w http.ResponseWriter, r *http.Request) {

	res := 0
	tmpl.ExecuteTemplate(w, "Forgotpassword", res)

}

func Startmap(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 301)
		return
	}
	uid := session.Values["id"]
	db := dbConn()
	nId := r.URL.Query().Get("id")
	var selDB *sql.Rows
	selDB, err := db.Query("SELECT * FROM custom_fare WHERE user_id=? order by id desc limit 1", uid)
	if err != nil {
		panic(err.Error())
	}
	if nId != "" {
		selDB1, err := db.Query("SELECT * FROM custom_fare WHERE id=?", nId)
		if err != nil {
			panic(err.Error())
		}
		selDB = selDB1

	}

	fare := Fare{}
	for selDB.Next() {
		var id, user_id int
		var fare_name, base_fare, min_fare, cost string
		err = selDB.Scan(&id, &fare_name, &base_fare, &min_fare, &cost, &user_id)
		if err != nil {
			panic(err.Error())
		}
		fare.Id = id
		fare.FareName = fare_name
		fare.BaseFare = base_fare
		fare.MinFare = min_fare
		fare.Cost = cost
		fare.UserId = user_id

	}
	tmpl.ExecuteTemplate(w, "Startmap", fare)
	defer db.Close()

	//tmpl.ExecuteTemplate(w, "Startmap", res)

}
func Stopmap(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 301)
		return
	}
	uid := session.Values["id"]
	db := dbConn()
	//var c *gin.Context
	now := time.Now()

	// fmt.Println("Year:", now.Year())
	// fmt.Println("Month:", now.Month())
	// fmt.Println("Day:", now.Day())
	// fmt.Println("Hour:", now.Hour())
	// fmt.Println("Minute:", now.Minute())
	// fmt.Println("Second:", now.Second())
	// fmt.Println("Nanosecond:", now.Nanosecond())

	stopDate := now.Day()
	stopMonth := now.Month()
	stopYear := now.Year()
	stopHour := now.Hour()
	stopMinute := now.Minute()
	loc := r.URL.Query().Get("loc")
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")
	date := r.URL.Query().Get("date")
	time := r.URL.Query().Get("time")
	cost := r.URL.Query().Get("cost")
	basefare := r.URL.Query().Get("basefare")
	from_address := r.URL.Query().Get("address")
	log.Println(loc)

	// data, _ := c.GetRawData()
	// jsonStream := string(data)
	// dec := json.NewDecoder(strings.NewReader(jsonStream))

	// for dec.More() {
	// 	var p Person
	// 	err := dec.Decode(&p)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("Hello %s\n", p.Loc)
	// }

	//res := loc
	data := make(map[string]interface{})
	data["lat"] = lat
	data["lon"] = lon
	data["date"] = date
	data["time"] = time
	data["stopDate"] = stopDate
	data["stopMonth"] = stopMonth
	data["stopYear"] = stopYear
	data["stopHour"] = stopHour
	data["stopMinute"] = stopMinute
	data["cost"] = cost
	data["basefare"] = basefare

	data["from_address"] = from_address

	type Ride struct {
		Distance string `json:"distance"`
		Wait     string `json:"wait"`
		Total    string `json:"total"`
	}
	if r.Method == "POST" {
		//decoding http request
		decoder := json.NewDecoder(r.Body)

		d := Ride{}

		// Decoder stores the parsed JSON into our user struct
		// fails on regular submit, pass on REST client submit.
		err := decoder.Decode(&d)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(d.Distance)
		fmt.Println(d.Wait)
		fmt.Println(d.Total)
	}
	if lat != "" {
		//

		insForm, err := db.Prepare("INSERT INTO ride_history(user_id, from_lat,from_lon,from_address,from_date,from_time,base_fare,cost) VALUES(?,?,?,?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm1, err := insForm.Exec(uid, lat, lon, from_address, date, time, basefare, cost)
		log.Println("INSERT: Date: " + date + " | lat: " + lat + " | lon: " + lon)

		lastinsertid, err := insForm1.LastInsertId()
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("lastinsertid: ", lastinsertid)

		data["lastinsertid"] = lastinsertid
	}
	//dialog.Alert("Custom Fare created!")
	//http.Redirect(w, r, "/customefare", 301)
	defer db.Close()

	ajax_post_data1 := r.FormValue("tolat")
	ajax_post_data2 := r.FormValue("tolat")
	fmt.Println(ajax_post_data1)
	fmt.Println(ajax_post_data2)

	tmpl.ExecuteTemplate(w, "Stopmap", data)

}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 301)
		return
	}
	//res := 0
	//var now time.Time
	currentTime := time.Now()
	startTime := time.Now().AddDate(0, 0, -1)
	fmt.Println("Current Time in String: ", currentTime.String())

	fmt.Println("MM/DD/YYYY : ", currentTime.Format("01/02/2006"))
	fmt.Println("MM/DD/YYYY : ", startTime.Format("01/02/2006"))

	user_id := session.Values["id"]
	db := dbConn()
	// selDB, err := db.Query("SELECT count(total) FROM  ride_history WHERE user_id=?", user_id)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// user := User{}
	// //emp := Employee{}
	// for selDB.Next() {
	// 	var id int
	// 	var name, email, mobile, password string
	// 	err = selDB.Scan(&id, &name, &email, &mobile, &password)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}

	// }

	var count, todayCount, yesterdayCount, totalTrip, duration int
	var distance float32
	data := make(map[string]interface{})

	err := db.QueryRow("SELECT COALESCE(SUM(total),0) as total FROM ride_history WHERE user_id=? AND to_date=?", user_id, currentTime.Format("01-02-2006")).Scan(&todayCount)
	switch {
	case err != nil:
		log.Fatal(err)
	default:

		data["todayCount"] = todayCount
		//res = count
		//fmt.Printf("Number of rows are %s\n", count)
	}

	err2 := db.QueryRow("SELECT COALESCE(SUM(total),0) as total ,COUNT(*), COALESCE(SUM(duration),0) as duration  ,COALESCE(SUM(distance),0) as distance  FROM ride_history WHERE user_id=? ", user_id).Scan(&count, &totalTrip, &duration, &distance)
	switch {
	case err2 != nil:
		log.Fatal(err2)
	default:

		data["count"] = count
		data["totalTrip"] = totalTrip
		data["duration"] = duration
		data["distance"] = distance
		//res = count
		//fmt.Printf("Number of rows are %s\n", count)
	}

	err3 := db.QueryRow("SELECT COALESCE(SUM(total),0) as total FROM ride_history WHERE user_id=? AND to_date=?", user_id, startTime.Format("01-02-2006")).Scan(&yesterdayCount)
	switch {
	case err3 != nil:
		log.Fatal(err3)
	default:

		data["yesterdayCount"] = yesterdayCount
		//res = count
		//fmt.Printf("Number of rows are %s\n", count)
	}

	tmpl.ExecuteTemplate(w, "Dashboard", data)

}
func Profile(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 301)
		return
	}
	username := session.Values["username"] // we get username from browser.
	id := session.Values["id"]             // we get id from browser.
	log.Println(username)
	if username == nil {
		http.Redirect(w, r, "/login", 301)
		//return c.Redirect(http.StatusSeeOther, "/login") // login firs.
	}
	//res := 0
	//myvar := map[string]interface{}{"MyVar": username, "ID": id}

	selDB, err := db.Query("SELECT id,name,email,mobile,password FROM users WHERE id=?", id)
	if err != nil {
		panic(err.Error())
	}
	user := User{}
	//emp := Employee{}
	for selDB.Next() {
		var id int
		var name, email, mobile, password string
		err = selDB.Scan(&id, &name, &email, &mobile, &password)
		if err != nil {
			panic(err.Error())
		}
		user.Id = id
		user.Name = name
		user.Mobile = mobile
		user.Email = email

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		user.Password = string(hashedPassword)
		fmt.Println(string([]byte(password)))
		fmt.Println(string(hashedPassword))
		//user.Password = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	}

	tmpl.ExecuteTemplate(w, "Profile", user)
	defer db.Close()

}
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 301)
		return
	}
	db := dbConn()
	if r.Method == "POST" {

		id := r.FormValue("uid")
		name := r.FormValue("name")
		email := r.FormValue("email")
		mobile := r.FormValue("mobile")
		password := r.FormValue("password")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error, unable to create your account.", 500)
			return
		}

		insForm, err := db.Prepare("UPDATE users SET name=?, email=?, mobile=? , password=?  WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, email, mobile, hashedPassword, id)
		log.Println("UPDATE: Name: " + name + " | Email: " + email + "| Mobile: " + mobile)
	}

	defer db.Close()
	http.Redirect(w, r, "/profile", 301)

}
func Ridehistory(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 301)
		return
	}
	user_id := session.Values["id"]
	db := dbConn()
	selDB, err := db.Query("SELECT id,COALESCE(from_date, '') as from_date,COALESCE(from_time,'') as from_time,COALESCE(from_address,'') as from_address ,COALESCE(to_address,'') as to_address, COALESCE(duration,'') as duration , COALESCE(distance,'') as distance ,COALESCE(total,'') as total  ,user_id FROM ride_history where user_id=? ORDER BY id DESC", user_id)
	if err != nil {
		panic(err.Error())
	}
	rideHistory := RideHistory{}
	res := []RideHistory{}
	for selDB.Next() {
		var id, user_id int
		var from_date, from_time, from_address, to_address, duration, distance, total string
		err = selDB.Scan(&id, &from_date, &from_time, &from_address, &to_address, &duration, &distance, &total, &user_id)
		if err != nil {
			panic(err.Error())
		}

		//}

		rideHistory.Id = id
		rideHistory.FromDate = from_date
		rideHistory.FromTime = from_time
		rideHistory.FromAddress = from_address
		rideHistory.ToAddress = to_address
		rideHistory.Duration = duration
		rideHistory.Distance = distance
		rideHistory.Total = total
		rideHistory.UserSession = user_id

		res = append(res, rideHistory)

	}

	tmpl.ExecuteTemplate(w, "Ridehistory", res)
	defer db.Close()

}
func Customefare(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 301)
		return
	}
	user_id := session.Values["id"]
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM custom_fare where user_id=? ORDER BY id DESC", user_id)
	if err != nil {
		panic(err.Error())
	}
	fare := Fare{}
	res := []Fare{}
	for selDB.Next() {
		var id, user_id int
		var fare_name, base_fare, min_fare, cost string
		err = selDB.Scan(&id, &fare_name, &base_fare, &min_fare, &cost, &user_id)
		if err != nil {
			panic(err.Error())
		}

		//}

		fare.Id = id
		fare.FareName = fare_name
		fare.BaseFare = base_fare
		fare.MinFare = min_fare
		fare.Cost = cost
		fare.UserId = user_id
		res = append(res, fare)

	}

	tmpl.ExecuteTemplate(w, "Customefare", res)
	defer db.Close()

}

func RideHistoryDdetail(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 301)
		return
	}
	db := dbConn()
	nId := r.URL.Query().Get("id")
	log.Println(nId)
	selDB, err := db.Query("SELECT id,COALESCE(from_time, '') as from_time, COALESCE(to_time, '') as to_time, COALESCE(from_address, '') as from_address, COALESCE(to_address, '') as to_address, COALESCE(duration, '') as duration, COALESCE(distance, '') as distance, COALESCE(total, '') as total, COALESCE(base_fare, '') as base_fare, COALESCE(cost, '') as cost, COALESCE(waiting, '') as waiting  FROM ride_history WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	// rideHistoryDetail := RideHistoryDetail{}
	// res1 := []RideHistoryDetail{}

	rideHistoryDetail := RideHistoryDetail{}
	for selDB.Next() {
		var id int
		var from_time, to_time, from_address, to_address, duration, distance, total, base_fare, cost, waiting string
		err = selDB.Scan(&id, &from_time, &to_time, &from_address, &to_address, &duration, &distance, &total, &base_fare, &cost, &waiting)
		if err != nil {
			panic(err.Error())
		}

		//}

		rideHistoryDetail.Id = id
		rideHistoryDetail.FromTime = from_time
		rideHistoryDetail.ToTime = to_time
		rideHistoryDetail.FromAddress = from_address
		rideHistoryDetail.ToAddress = to_address
		rideHistoryDetail.Duration = duration
		rideHistoryDetail.Distance = distance
		rideHistoryDetail.Total = total
		rideHistoryDetail.BaseFare = base_fare
		rideHistoryDetail.Cost = cost
		rideHistoryDetail.Waiting = waiting
		//rideHistory.UserSession = user_id

		//res1 = append(res1, rideHistoryDetail)

	}
	tmpl.ExecuteTemplate(w, "Ride-history-detail", rideHistoryDetail)
	defer db.Close()

}
func Home(w http.ResponseWriter, r *http.Request) {

	res := 0

	tmpl.ExecuteTemplate(w, "Home", res)

}
func FareSetting(w http.ResponseWriter, r *http.Request) {

	res := 0

	tmpl.ExecuteTemplate(w, "Fare-setting", res)

}
func FareSettingEdit(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 301)
		return
	}
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT id,fare_name,base_fare,min_fare,cost FROM custom_fare WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	fare := Fare{}
	for selDB.Next() {
		var id int
		var fare_name, base_fare, min_fare, cost string
		err = selDB.Scan(&id, &fare_name, &base_fare, &min_fare, &cost)
		if err != nil {
			panic(err.Error())
		}

		fare.Id = id
		fare.FareName = fare_name
		fare.BaseFare = base_fare
		fare.MinFare = min_fare
		fare.Cost = cost

	}

	tmpl.ExecuteTemplate(w, "Fare-setting-edit", fare)
	defer db.Close()

}
func FareInsert(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 301)
		return
	}
	db := dbConn()
	if r.Method == "POST" {

		user_id := session.Values["id"]
		fare_name := r.FormValue("fare_name")
		base_fare := r.FormValue("base_fare")
		min_fare := r.FormValue("min_fare")
		cost := r.FormValue("cost")

		insForm, err := db.Prepare("INSERT INTO custom_fare(fare_name, base_fare,min_fare,cost,user_id) VALUES(?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(fare_name, base_fare, min_fare, cost, user_id)
		log.Println("INSERT: Name: " + fare_name + " | Base_Fare: " + base_fare + " | Cost: " + cost)
	}
	//dialog.Alert("Custom Fare created!")
	message := "Custom Fare created Successfully"
	http.Redirect(w, r, "/success?success_message="+message+"&page=customefare", 301)
	// http.Redirect(w, r, "/customefare", 301)
	defer db.Close()

}
func FareUpdate(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 301)
		return
	}
	db := dbConn()
	if r.Method == "POST" {
		id := r.FormValue("id")
		fare_name := r.FormValue("fare_name")
		base_fare := r.FormValue("base_fare")
		min_fare := r.FormValue("min_fare")
		cost := r.FormValue("cost")

		insForm, err := db.Prepare("UPDATE custom_fare SET fare_name=?, base_fare=?, min_fare=?,cost=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(fare_name, base_fare, min_fare, cost, id)

		defer db.Close()
		// dialog.Alert("Custom Fare Updated!")
		// http.Redirect(w, r, "/customefare", 301)

		message := "Custom Fare Updated Successfully"
		http.Redirect(w, r, "/success?success_message="+message+"&page=customefare", 301)
	}
}
func FareDelete(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", 301)
		return
	}
	db := dbConn()
	id := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM custom_fare WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(id)
	// dialog.Alert("Custom Fare Deleted!")

	message := "Custom Fare Deleted Successfully"
	http.Redirect(w, r, "/success?success_message="+message+"&page=customefare", 301)

	defer db.Close()
	// http.Redirect(w, r, "/customefare", 301)
}

func Email(w http.ResponseWriter, r *http.Request) {
	// Sender data.
	from := "anand.kcet@gmail.com"
	password := "rkokaphbxflpjmkb"

	// Receiver email address.
	to := []string{
		"test@yopmail.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "25"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Puneet Singh",
		Message: "This is a test message in a HTML template",
	})

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Error(w, "Email Sent!", http.StatusAccepted)
	fmt.Println("Email Sent!")
}
func receiveAjax(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {

		tolat := r.FormValue("tolat")
		tolon := r.FormValue("tolon")
		stopdate := r.FormValue("stopdate")
		stoptime := r.FormValue("stoptime")
		distance := r.FormValue("distance")
		wait := r.FormValue("wait")
		duration := r.FormValue("duration")
		total := r.FormValue("total")
		to_address := r.FormValue("to_address")
		lastid := r.FormValue("lastid")

		insForm, err := db.Prepare("UPDATE ride_history SET to_lat=?, to_lon=?, to_address=?, to_date=?,to_time=?,distance=?,waiting=?,duration=?,total=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(tolat, tolon, to_address, stopdate, stoptime, distance, wait, duration, total, lastid)
		// fmt.Println(ajax_post_data1)
		// fmt.Println(ajax_post_data2)
		//ajax_post_data := r.FormValue("ajax_post_data")
		fmt.Println("Receive ajax post data string ", duration)
		w.Write([]byte("<h2>after<h2>"))
		defer db.Close()
	}
}
func Amount(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {

		amount_status := r.FormValue("amount_status")
		lastid := r.FormValue("lastid")

		insForm, err := db.Prepare("UPDATE ride_history SET amount_status=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(amount_status, lastid)
		// fmt.Println(ajax_post_data1)
		// fmt.Println(ajax_post_data2)
		//ajax_post_data := r.FormValue("ajax_post_data")
		//fmt.Println("Receive ajax post data string ", duration)
		w.Write([]byte("<h2>after<h2>"))
		defer db.Close()
	}
}
func MapTest(w http.ResponseWriter, r *http.Request) {

	res := 0
	tmpl.ExecuteTemplate(w, "MapTest", res)
}
func MapTest1(w http.ResponseWriter, r *http.Request) {

	res := 0
	tmpl.ExecuteTemplate(w, "MapTest1", res)
}

func main() {
	//fs := http.FileServer(http.Dir("/css/"))
	//http.Handle("/css/", fs)

	log.Println("Server started on: http://localhost:8381")
	//	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("http/css"))))

	http.Handle("/css/",
		http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	http.Handle("/Company/",
		http.StripPrefix("/Company/", http.FileServer(http.Dir("Company"))))

	http.Handle("/public/",
		http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	http.HandleFunc("/email", Email)
	http.HandleFunc("/startmap", Startmap)
	http.HandleFunc("/stopmap", Stopmap)

	http.HandleFunc("/dashboard", Dashboard)
	http.HandleFunc("/profile", Profile)
	http.HandleFunc("/ridehistory", Ridehistory)
	http.HandleFunc("/ride-history-detail", RideHistoryDdetail)

	http.HandleFunc("/customefare", Customefare)
	http.HandleFunc("/fare-setting", FareSetting)
	http.HandleFunc("/fare-insert", FareInsert)
	http.HandleFunc("/fare-setting-edit", FareSettingEdit)
	http.HandleFunc("/fare-update", FareUpdate)
	http.HandleFunc("/fare-delete", FareDelete)

	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/login-check", Logincheck)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/register-user", RegisterUser)

	http.HandleFunc("/success", Success)
	http.HandleFunc("/error", Error)

	http.HandleFunc("/user-update", UserUpdate)
	http.HandleFunc("/forgotpassword", Forgotpassword)
	http.HandleFunc("/", Home)

	http.HandleFunc("/receive", receiveAjax)
	http.HandleFunc("/amount", Amount)

	http.HandleFunc("/map-test", MapTest)
	http.HandleFunc("/map-test1", MapTest1)

	//http.ListenAndServe(":8380", nil)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8381" // Default port if not specified
	}
	http.ListenAndServe(":"+port, nil)
}
