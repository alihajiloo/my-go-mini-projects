package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type User struct {
	tableName struct{} `pg:"users_table"`
	ID        string   `json:"id" pg:"id,pk,default:gen_random_uuid()"`
	FirstName string   `json:"first_name,omitempty" pg:"first_name"`
	LastName  string   `json:"last_name,omitempty" pg:"last_name"`
	Status    string   `json:"status,omitempty" pg:"status"`
	Password  string   `json:"password" pg:"password"`
}

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
		Database: "plllatform_com_db",
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

	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	fmt.Println(newUser)
	InsertDB(NewDBConn(), newUser)

	// c.IndentedJSON(http.StatusCreated, newUser)
	c.JSON(http.StatusCreated, gin.H{
		"message": "OK",
	})
}

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
