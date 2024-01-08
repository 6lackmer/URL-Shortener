package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRedirectHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Id", "Original_url", "Short_url"}).
		AddRow(1, "http://example.com", "abc123")

	mock.ExpectQuery("^SELECT \\* FROM url_shortener WHERE Short_url = \\?$").
		WithArgs("abc123").
		WillReturnRows(rows)

	app := &MyApp{db: &MySQLDatabase{DB: db}}

	req := httptest.NewRequest("GET", "/abc123", nil)
	rr := httptest.NewRecorder()

	app.redirectHandler(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	expectedLocation := "http://example.com"
	location := rr.Header().Get("Location")
	if location != expectedLocation {
		t.Errorf("handler returned unexpected location: got %v want %v", location, expectedLocation)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestFormHandler_NonPostRequest(t *testing.T) {
	db, _, _ := sqlmock.New()
	app := &MyApp{db: &MySQLDatabase{DB: db}}

	req := httptest.NewRequest("GET", "/submit", nil)
	rr := httptest.NewRecorder()

	app.formHandler(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}
}

func TestFormHandler_EmptyUserInput(t *testing.T) {
	db, _, _ := sqlmock.New()
	app := &MyApp{db: &MySQLDatabase{DB: db}}

	req := httptest.NewRequest("POST", "/submit", nil)
	rr := httptest.NewRecorder()

	app.formHandler(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}
}

func TestFormHandler_InvalidURL(t *testing.T) {
	db, _, _ := sqlmock.New()
	app := &MyApp{db: &MySQLDatabase{DB: db}}

	form := strings.NewReader("textInput=invalidurl")
	req := httptest.NewRequest("POST", "/submit", form)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	app.formHandler(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}
}

func TestFormHandler_ExistingURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Id", "Original_url", "Short_url"}).
		AddRow(1, "http://example.com", "abc123")

	mock.ExpectQuery("^SELECT \\* FROM url_shortener WHERE Original_url = \\?$").
		WithArgs("http://example.com").
		WillReturnRows(rows)

	app := &MyApp{db: &MySQLDatabase{DB: db}}

	form := strings.NewReader("textInput=http://example.com")
	req := httptest.NewRequest("POST", "/submit", form)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	app.formHandler(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFormHandler_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("^SELECT \\* FROM url_shortener WHERE Original_url = \\?$").
		WithArgs("https://example.com").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery("^SELECT \\* FROM url_shortener$").
		WillReturnRows(sqlmock.NewRows([]string{"Short_url"}))

	mock.ExpectExec("^INSERT INTO url_shortener").
		WithArgs(sqlmock.AnyArg(), "https://example.com", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	app := &MyApp{db: &MySQLDatabase{DB: db}}

	form := strings.NewReader("textInput=https://example.com")
	req := httptest.NewRequest("POST", "/submit", form)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	app.formHandler(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	location := rr.Header().Get("Location")
	if !strings.Contains(location, "/?success=shortened") {
		t.Errorf("handler returned unexpected location header: got %v want '/?success=shortened'", location)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestViewUrlsHandler_TemplateExecutionError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Id", "Original_url", "Short_url"}).
		AddRow(1, "http://example.com", "xyz123").
		AddRow(2, "http://example.org", "abc123")
	mock.ExpectQuery("^SELECT \\* FROM url_shortener$").WillReturnRows(rows)

	tmpl, err := template.New("viewurls.html").Parse("{{range .}}{{.NonExistentField}}{{end}}")
	if err != nil {
		t.Fatalf("Failed to parse mock template: %v", err)
	}

	app := &MyApp{db: &MySQLDatabase{DB: db}, tmpl: tmpl}

	req := httptest.NewRequest("GET", "/viewurls", nil)
	rr := httptest.NewRecorder()

	app.viewUrlsHandler(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestViewUrlsHandler_DatabaseRetrievalError(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a mock database connection", err)
    }
    defer db.Close()

    mock.ExpectQuery("^SELECT \\* FROM url_shortener$").WillReturnError(sql.ErrNoRows)

    tmpl, err := template.New("viewurls.html").Parse("{{range .}}{{.}}{{end}}")
    if err != nil {
        t.Fatalf("Failed to create mock template: %v", err)
    }

    app := &MyApp{db: &MySQLDatabase{DB: db}, tmpl: tmpl}

    req := httptest.NewRequest("GET", "/viewurls", nil)
    rr := httptest.NewRecorder()

    app.viewUrlsHandler(rr, req)

    // Check if the status code is 500 Internal Server Error
    if status := rr.Code; status != http.StatusInternalServerError {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
    }

    // Verify that all expectations on the mock were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}

func TestIndexHandler_Redirect(t *testing.T) {
    db, _, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a mock database connection", err)
    }
    defer db.Close()

    app := &MyApp{
        db: &MySQLDatabase{DB: db},
    }

    req, err := http.NewRequest("GET", "/somepath", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(app.indexHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusNotFound)
    }
}

func TestNewMyApp(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()

	tmpl, err := template.New("test").Parse("{{.}}")
	if err != nil {
		t.Fatalf("Failed to create mock template: %v", err)
	}

	myApp := NewMyApp(&MySQLDatabase{DB: db}, tmpl)

	if myApp.db == nil {
		t.Errorf("NewMyApp did not correctly initialize the db field")
	}

	if myApp.tmpl == nil {
		t.Errorf("NewMyApp did not correctly initialize the tmpl field")
	}
}