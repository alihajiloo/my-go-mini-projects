package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// user represents data about a record user.
// type user struct {
// 	ID     string  `json:"id"`
// 	Title  string  `json:"title"`
// 	Artist string  `json:"artist"`
// 	Price  float64 `json:"price"`
// }

// User data type
type User struct {
	tableName struct{} `pg:"users_table"`
	ID        string   `json:"id" pg:"id,pk"`
	FirstName string   `json:"first_name,omitempty" pg:"first_name"`
	LastName  string   `json:"last_name,omitempty" pg:"last_name"`
	Status    string   `json:"status,omitempty" pg:"status"`
}

// type User struct {
// 	Id     int64
// 	Name   string
// 	Emails []string
// }

// users slice to seed record user data.
// var users = []User{
// 	{ID: "1", FirstName: "Blue Train", LastName: "John Coltrane", Status: "56.99"},
// 	{ID: "2", FirstName: "Jeru", LastName: "Gerry Mulligan", Status: "17.99"},
// 	{ID: "3", FirstName: "Sarah Vaughan and Clifford Brown", LastName: "Sarah Vaughan", Status: "39.99"},
// }

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)
	// router.GET("/users/:id", getUserByID)
	router.POST("/users", postUsers)

	db := NewDBConn()
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		fmt.Println(err)
	}

	router.Run("localhost:8080")
}

func NewDBConn() (con *pg.DB) {
	address := fmt.Sprintf("%s:%s", "localhost", "5432")
	options := &pg.Options{
		User: "alihajiloo",
		// Password: "12345678",
		Addr:     address,
		Database: "termestyle_db",
		PoolSize: 50,
	}
	con = pg.Connect(options)
	if con == nil {
		fmt.Println("cannot connect to postgres")
	}

	return
}

func InsertDB(pg *pg.DB, user User) error {
	_, err := pg.Model(&user).Insert()
	return err
}

func SelectDBUser(pg *pg.DB) ([]User, error) {
	var users []User
	err := pg.Model(&users).Select()
	fmt.Println(users)
	return users, err
}

// getUsers responds with the list of all users as JSON.
func getUsers(c *gin.Context) {
	// var users []User

	users, err := SelectDBUser(NewDBConn())
	if err != nil {
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

// postUsers adds an user from JSON received in the request body.
func postUsers(c *gin.Context) {
	var newUser User

	// user := c.BindJSON(&newUser)

	// Call BindJSON to bind the received JSON to
	// newUser.

	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	fmt.Println(newUser)
	InsertDB(NewDBConn(), newUser)

	// Add the new user to the slice.
	// users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

// getUserByID locates the user whose ID value matches the id
// parameter sent by the client, then returns that user as a response.
// func getUserByID(c *gin.Context) {
// 	id := c.Param("id")

// 	// Loop through the list of users, looking for
// 	// an user whose ID value matches the parameter.
// 	for _, a := range users {
// 		if a.ID == id {
// 			c.IndentedJSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
// }

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
