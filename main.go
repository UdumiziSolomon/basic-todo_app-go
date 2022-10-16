package main

import (
	"net/http"
	"errors"
	"os"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/UdumiziSolomon/todo-app/util"
)

type todo struct {
	ID          string   `json: "id"`
	Item        string   `json: "item"`
	Completed   bool     `json: "completed"`
}

var todos = []todo {
	{ ID: "1", Item: "Fix Bug", Completed: true },
	{ ID: "2", Item: "Add Feature", Completed: false },
	{ ID: "3", Item: "Open PR", Completed: true },
}

//  GET ALL TODOS
func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)   
}

// ADD A TODO TO THE ARRAY
func addTodos(context *gin.Context) {
	var newTodo todo     // new todo instance from client

	if err := context.BindJSON(&newTodo); err != nil {
		return
	}

	// append the new todo to the todo array
	todos = append(todos, newTodo)
	// return the todo
	context.IndentedJSON(http.StatusCreated, todos)
}

// GE A TODO BY ID    
func getTodoById(id string) (*todo, error) {    //returns either the todo or error
	for i, td := range todos {
		if td.ID == id {
			return &todos[i], nil
		}
	}

	return nil, errors.New("Todo was not found")
}

func getSingleTodo(context *gin.Context){
	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "TODO not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, todo) 
}

func toggleTodoStatus(context *gin.Context){
	var newTodo todo     // new todo instance from client

	if err := context.BindJSON(&newTodo); err != nil {
		return
	}

	id := context.Param("id")
	todo, err := getTodoById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "TODO not found"})   // custom Error instead of conventiona err.Error()
	}

	todo.Completed = newTodo.Completed  // patch the data with new ones
	todo.ID = newTodo.ID
	todo.Item = newTodo.Item

	context.IndentedJSON(http.StatusOK, todo)
}

// RETRIEVE LOGS FOR REQUESTS
func retrieveRequestLogs() {
	file, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
}

func main() {

	// LOAD ENV VARIANLE FROM VIPER
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config variables", err)
	}

	router := gin.New();
	
	// gin mthd to recover from panic
	router.Use(gin.Recovery(), gin.Logger())       // ==> router := gin.Default()
	
	retrieveRequestLogs()   // func to retrieve every request ==> gin.log file

	router.GET("/todos", getTodos)
	router.POST("/todo", addTodos)
	router.GET("/todos/:id", getSingleTodo)
	router.PATCH("/todos/:id", toggleTodoStatus)


	// initializing the server
	router.Run(config.ServerPort)
}