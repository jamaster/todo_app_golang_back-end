# todo_app_golang_back-end

## Description
This is a simple todo app back-end written in Golang. It uses a bolt database to store the data.


# Requirements
* go version 1.20 or hire
  * [GoLang](https://go.dev/)


# Run tests
* run tests `go test ./...`

# Local development
* clone repository `git clone https://github.com/jamaster/todo_app_golang_back-end.git`
* run application `go run main.go`
    * visit http://localhost:8080 in your browser

# Run in production
* build application `go build -o todoapp main.go`
* run application `./todoapp`

# Run in docker
* docker build -t my-todo-app .
* docker run -p 8080:8080 my-todo-app
* visit http://localhost:8080 in your browser

# Run with docker-compose
* docker-compose up
* visit http://localhost:8080 in your browser
