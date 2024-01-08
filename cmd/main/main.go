package main

import (
	"cmd/main/internal"
	"cmd/main/pkg"
	"cmd/main/pkg/Storage/MySql"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type UrlShortener struct {
	Id           int
	Original_url string
	Short_url    string
}

type Database interface {
	GetByWhere(table string, whereClause string, args []interface{}, dest interface{}) error
	GetAll(table string, dest interface{}) error
	Save(table string, data interface{}) error
}

type MySQLDatabase struct {
	DB *sql.DB
}

func (m *MySQLDatabase) GetByWhere(table string, whereClause string, args []interface{}, dest interface{}) error {
	return MySql.GetByWhere(m.DB, table, whereClause, args, dest)
}

func (m *MySQLDatabase) GetAll(table string, dest interface{}) error {
	return MySql.GetAll(m.DB, table, dest)
}

func (m *MySQLDatabase) Save(table string, data interface{}) error {
	return MySql.Save(m.DB, table, data)
}

type MyApp struct {
	db   *MySQLDatabase
	tmpl *template.Template
}

func NewMyApp(db *MySQLDatabase, tmpl *template.Template) *MyApp {
	return &MyApp{
		db:   db,
		tmpl: tmpl,
	}
}

// Handles the form submission and validation of user input 
func (app *MyApp) formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	userInput := r.FormValue("textInput") 
	
	if userInput == "" {
		http.Redirect(w, r, "/?error=no_input", http.StatusSeeOther)
		return
	} else if !pkg.IsValidURL(userInput) {
		http.Redirect(w, r, "/?error=invalid_url", http.StatusSeeOther)
		return
	}

	var existingUrlShortener UrlShortener
	err := app.db.GetByWhere("url_shortener", "Original_url = ?", []interface{}{userInput}, &existingUrlShortener)
	if err == nil {
		log.Println("URL already exists in database: " + userInput)
		http.Redirect(w, r, "/?error=url_exists", http.StatusSeeOther)
		return
	}

	var allShortUrls []string
	var urlShortenerData []UrlShortener
	err = app.db.GetAll("url_shortener", &urlShortenerData)
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range urlShortenerData {
		allShortUrls = append(allShortUrls, result.Short_url)
	}

	shortUrl := pkg.GetUniqueShortUrl(allShortUrls, 5)
	newUrlShortener := UrlShortener{
		Id:           len(urlShortenerData) + 1,
		Original_url: userInput,
		Short_url:    shortUrl,
	}

	err = app.db.Save("url_shortener", &newUrlShortener)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/?success=shortened", http.StatusSeeOther)
}

// Handles the redirecting of the user to the original url
func (app *MyApp) redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path[1:]
	var urlShortener UrlShortener
	err := app.db.GetByWhere("url_shortener", "Short_url = ?", []interface{}{shortUrl}, &urlShortener)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, urlShortener.Original_url, http.StatusFound)
}

// handles the viewurls route. Allowing the user to view all the urls and their shortened versions
func (app *MyApp) viewUrlsHandler(w http.ResponseWriter, r *http.Request) {
	var urlShortenerData []UrlShortener
	err := app.db.GetAll("url_shortener", &urlShortenerData)
	if err != nil {
		log.Printf("Error retrieving data: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = app.tmpl.ExecuteTemplate(w, "viewurls.html", urlShortenerData)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// setupRoutes sets up the routes for the application
func (app *MyApp) setupRoutes() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/submit", app.formHandler)
	http.HandleFunc("/", app.indexHandler)
	http.HandleFunc("/viewurls", app.viewUrlsHandler)
}

// indexHandler handles the root route
func (app *MyApp) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.redirectHandler(w, r)
		return
	}
	http.ServeFile(w, r, "static/templates/index.html")
}

func main() {
	db, err := internal.ConnectToMySqlDB() // connect to database
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	internal.InitMySqlDB(db) // Make sure database is set up 

	tmpl := template.Must(template.ParseGlob("static/templates/*.html"))	// parse the templates
	myApp := NewMyApp(&MySQLDatabase{DB: db}, tmpl) 

	myApp.setupRoutes() // set up routes

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

