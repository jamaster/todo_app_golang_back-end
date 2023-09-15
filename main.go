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
	http.HandleFunc("/list", listTask(db))
	http.ListenAndServe(":8080", nil)
}

// implement method handler to show the home page of the application with imput box to add task. Submit button will call addTask method
func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<html>
			<head>
				<title>Task Manager</title>
			</head>
            <script>
			  function fetchTaskList() {
				var xhr = new XMLHttpRequest();
				xhr.onreadystatechange = function() {
				  if (xhr.readyState === 4 && xhr.status === 200) {
					var taskList = document.getElementById("taskList");
					taskList.innerHTML = ""; // Очистити список перед оновленням
			
					// Розділити текст на рядки (передполагається, що кожен рядок - це одне завдання)
					var taskLines = xhr.responseText.split('\n');
			
					// Додати кожне завдання до списку
					for (var i = 0; i < taskLines.length; i++) {
					  var task = taskLines[i];
					  if (task.trim() !== "") {
						var listItem = document.createElement("li");
						listItem.textContent = task;
                        //add delete button to delete task
						var deleteButton = document.createElement("button");
						deleteButton.textContent = "Delete";
						deleteButton.onclick = function() {
						<!-- how to get task name here -->
						  var taskText = this.parentNode.textContent;
                          // task is substring between Task: and Created:
						  var task = taskText.substring(6, taskText.indexOf("Created:") - 1);
                 		  var xhr = new XMLHttpRequest();
						  xhr.open("POST", "/delete", true);
						  xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
						  xhr.send("task=" + encodeURIComponent(task));
                          // redirect to home page
						  window.location.href = "/";

						};
						listItem.appendChild(deleteButton);
						taskList.appendChild(listItem);
					  }
					}
				  }
				};
				xhr.open("GET", "/list", true);
				xhr.send();
			  }
			</script>

			<body onload="fetchTaskList()">
				<h1>Task Manager</h1>
				<h2>Add Task</h2>
       
				<form action="/add" method="post">
					<input type="text" name="task">
					<input type="submit" value="Add">
				</form>


				<h2>Tasks</h2>               
				<ul id="taskList"></ul>
			</body>
		</html>
	`))
}

// BoldDB interface
type BoldDB interface {
	Update(fn func(*bolt.Tx) error) error
	Close() error
	View(fn func(*bolt.Tx) error) error
}

// implement method addTask to add task (with the created time ) in the database and redirect to home page
func addTask(db BoldDB) http.HandlerFunc {
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

// implement method listTask to list all the tasks from the database
func listTask(db BoldDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Tasks"))
			b.ForEach(func(k, v []byte) error {
				w.Write([]byte("Task: " + string(k) + " Created: " + string(v) + "\n"))
				return nil
			})
			return nil

		})
	}
}
