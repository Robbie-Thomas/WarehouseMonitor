package api

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/robbie-thomas/fullstack/api/controllers"
	"github.com/robbie-thomas/fullstack/api/seed"
	"log"
	"os"
)

var server = controllers.Server{}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialiser(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	//server.Initialiser("mysql","root","password","3306","127.0.0.1","fullstack_api")
	seed.Load(server.DB)

	server.Run(":8080")

}
