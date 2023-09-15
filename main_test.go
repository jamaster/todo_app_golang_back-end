package main

// unit test for addTask method using httptest package and mock database
import (
	"bytes"
	"github.com/boltdb/bolt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock bolt database
type MockBoltDB struct {
	*bolt.DB
}

// implement method Update
func (db *MockBoltDB) Update(fn func(*bolt.Tx) error) error {
	return nil
}

// implement method Close
func (db *MockBoltDB) Close() error {
	return nil
}

// implement view method
func (db *MockBoltDB) View(fn func(*bolt.Tx) error) error {
	return nil
}

// test addTask method
func TestAddTask(t *testing.T) {
	// create mock database
	db := &MockBoltDB{}
	// create request to add task
	req, err := http.NewRequest("POST", "/add", bytes.NewBufferString("task=task1"))
	if err != nil {
		t.Fatal(err)
	}
	// create response recorder to record the response
	rr := httptest.NewRecorder()
	// create handler
	handler := http.HandlerFunc(addTask(db))
	// call handler
	handler.ServeHTTP(rr, req)
	// check for status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// test listTask method
func TestListTask(t *testing.T) {
	// create mock database
	db := &MockBoltDB{}
	// create request to list task
	req, err := http.NewRequest("GET", "/list", nil)
	if err != nil {
		t.Fatal(err)
	}
	// create response recorder to record the response
	rr := httptest.NewRecorder()
	// create handler
	handler := http.HandlerFunc(listTask(db))
	// call handler
	handler.ServeHTTP(rr, req)
	// check for status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
