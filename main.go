package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"tawesoft.co.uk/go/dialog"
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

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "solotaxi"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func LoginPage(w http.ResponseWriter, r *http.Request) {

	res := 0

	tmpl.ExecuteTemplate(w, "Login", res)

}
func Logincheck(res http.ResponseWriter, req *http.Request) {

	db := dbConn()
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

		dialog.Alert("Username Incorrect")
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		dialog.Alert("Password incorrect")
		http.Redirect(res, req, "/login", 301)
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
	//res.Write([]byte("Hello" + databaseUsername))
	http.Redirect(res, req, "/startmap", 301)

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

	var user string

	err := db.QueryRow("SELECT email FROM users WHERE email=?", email).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO users(name,mobile,email, password) VALUES(?, ?, ?, ?)", name, mobile, email, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}
		defer db.Close()
		//res.Write([]byte("User created!"))
		dialog.Alert("User created!")
		http.Redirect(res, req, "/register", 301)
		return

	case err != nil:
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(res, req, "/register", 301)

	}

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

	tmpl.ExecuteTemplate(w, "Stopmap", data)

}

func Dashboard(w http.ResponseWriter, r *http.Request) {

	res := 0

	tmpl.ExecuteTemplate(w, "Dashboard", res)

}
func Profile(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	session, _ := store.Get(r, "cookie-name")
	username := session.Values["username"] // we get email from browser.
	id := session.Values["id"]             // we get email from browser.
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

	res := 0

	tmpl.ExecuteTemplate(w, "Ridehistory", res)

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

	res := 0

	tmpl.ExecuteTemplate(w, "Ride-history-detail", res)

}
func Home(w http.ResponseWriter, r *http.Request) {

	res := 0

	tmpl.ExecuteTemplate(w, "Home", res)

}
func FareSetting(w http.ResponseWriter, r *http.Request) {

	res := 0

	tmpl.ExecuteTemplate(w, "Fare-setting", res)

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
	dialog.Alert("Custom Fare created!")
	http.Redirect(w, r, "/customefare", 301)
	defer db.Close()

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

func main() {
	//fs := http.FileServer(http.Dir("/css/"))
	//http.Handle("/css/", fs)

	log.Println("Server started on: http://localhost:8380")
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

	http.HandleFunc("/customefare", Customefare)
	http.HandleFunc("/fare-setting", FareSetting)
	http.HandleFunc("/fare-insert", FareInsert)

	http.HandleFunc("/ride-history-detail", RideHistoryDdetail)

	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/login-check", Logincheck)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/register-user", RegisterUser)
	http.HandleFunc("/user-update", UserUpdate)
	http.HandleFunc("/forgotpassword", Forgotpassword)
	http.HandleFunc("/", Home)

	//http.ListenAndServe(":8380", nil)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8380" // Default port if not specified
	}
	http.ListenAndServe(":"+port, nil)
}
