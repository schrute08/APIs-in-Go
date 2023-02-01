package StatusChecker

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmitHandler(t *testing.T) {
	// Test if the handler returns a "Method Not Allowed" error for a GET request
	req, err := http.NewRequest("GET", "/POST/websites", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SubmitHandler(NewFunc()))
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	// Test if the handler returns a "OK" status for a POST request
	var jsonStr []byte
	for website := range websitesMap {
		jsonStr = append(jsonStr, []byte(`"`+website+`",`)...)
	}
	jsonStr = append(append([]byte(`{"websites":[`), jsonStr...), []byte(`]}`)...)

	req, err = http.NewRequest("POST", "/POST/websites", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Test if the JSON payload from the POST request body is properly decoded and stored in the "websitesMap"
	for website := range websitesMap {
		if websitesMap[website].URL != website {
			t.Errorf("website not properly stored in websitesMap: got %v want %v", websitesMap[website].URL, website)
		}
	}
}

func TestStatusHandler(t *testing.T) {
	// Test if the handler returns a "Method Not Allowed" error for a POST request
	req, err := http.NewRequest("POST", "/GET/websites", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(StatusHandler(NewFunc()))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

	// Test if the handler returns a "OK" status for a GET request
	req, err = http.NewRequest("GET", "/GET/websites", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Test if the "id" parameter is present in the GET request query string and the status of the specified website is returned
	req, err = http.NewRequest("GET", "/GET/websites?id=www.medium.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected := `{"www.medium.com": "UP"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Test if the handler returns a "Not Found" error for an unknown website id in the GET request query string
	req, err = http.NewRequest("GET", "/status?id=www.medium.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}