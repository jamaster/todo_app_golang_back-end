package main

import (
	"github.com/boltdb/bolt"
	"log"
	"net/http"
	"time"
)

func main() {
	// create bolt database connection
	db := NewBoltDB("my.db")
	defer db.Close()

	// create bucket in bolt database	to store data of task with name of the task and time of creataion of task
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Tasks"))
		if err != nil {
			return err
		}
		return nil
	})

	// create http server
	httpServer(db)

}

// implement method NewBoltDB
func NewBoltDB(dbPath string) *bolt.DB {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// implement method httpServer
func httpServer(db *bolt.DB) {
	http.HandleFunc("/", handler)
	http.HandleFunc("/add", addTask(db))
	http.HandleFunc("/delete", deleteTask(db))
	//http.HandleFunc("/update", updateTask(db))
	//http.HandleFunc("/list", listTask(db))
	http.ListenAndServe(":8080", nil)
}

// implement method handler to show the home page of the application with imput box to add task. Submit button will call addTask method
func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<html>
			<head>
				<title>Task Manager</title>
			</head>
			<body>
				<form action="/add" method="post">
					<input type="text" name="task">
					<input type="submit" value="Add">
				</form>
			</body>
		</html>
	`))
}

// implement method addTask to add task (with the created time ) in the database and redirect to home page
func addTask(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		task := r.PostForm.Get("task")
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Tasks"))
			err := b.Put([]byte(task), []byte(time.Now().Format(time.RFC850)))
			return err
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// implement method deleteTask to delete task from the database and redirect to home page
func deleteTask(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		task := r.PostForm.Get("task")
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Tasks"))
			err := b.Delete([]byte(task))
			return err
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}
}