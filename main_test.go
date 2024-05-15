package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/get", nil) //(method, url , body reader)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}


func TestCreateRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	loginJson := []byte(`{
		"id":220,
		"name":"test-name-110"
	}`)
	bodyReader := bytes.NewReader(loginJson)

	//NewRequest returns a new incoming server Request, suitable for passing to an http.Handler for testing.
	req, err := http.NewRequest(http.MethodPost, "/create", bodyReader) //(method, url , body reader)

	if err != nil {
		fmt.Printf("Could request, Error : %s \n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}


//delete/:id
func TestDeleteRoute(t *testing.T) {
	router := setupRouter()

	req, err := http.NewRequest(http.MethodDelete, "/delete/220", nil) //(method, url , body reader)

	if err != nil {
		fmt.Printf("Could request, Error : %s \n", err)
		os.Exit(1)
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}


//update/:id
func TestUpdateRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	updateJson := []byte(`{
		"name":"updated-220",
		"age":20,
		"address":"updated-Mumbai220"
	}`)
	bodyReader := bytes.NewReader(updateJson)

	req, err := http.NewRequest(http.MethodPut, "/update/220", bodyReader)

	if err != nil {
		fmt.Printf("Could request, Error : %s \n", err)
		os.Exit(1)
	}

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}



//auth
func TestGetAuthRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/auth", nil)

	req.Header.Add("Role","admin")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetSessionTestRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/session-test", nil)
	req.Header.Add("Role","admin")

	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestLogoutRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/logout", nil)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}


