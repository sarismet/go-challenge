package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-challenge/db"
	"github.com/go-challenge/models"
)

func Hello() {
	fmt.Println("Hello from endpoints")
}

type App struct {
	RedisDB *db.RedisDatabase
	MongoDB *db.MongoDB
}

func Init() {
	app := App{}
	RedisDB, err := db.NewRedisDatabase()
	if err != nil || RedisDB == nil {
		log.Fatal("Error as conencting to Redis")
		return
	}

	MongoDB, err := db.ConnectMongoDB()
	if err != nil || MongoDB == nil {
		log.Fatal("Error as conencting to Sql")
		return
	}
	app.RedisDB = RedisDB
	app.MongoDB = MongoDB

	http.HandleFunc("/fetchData", app.fetchData)
	http.HandleFunc("/postKey", app.postKey)
	http.HandleFunc("/getKey", app.getKey)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func (app *App) fetchData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the fetchData!")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", body)
	fetchRequestModel := models.FetchRequestModel{}
	json.Unmarshal([]byte(body), &fetchRequestModel)
	fmt.Printf("StartDate : %s", fetchRequestModel.StartDate)

	app.MongoDB.FetchDataFromMongoDB()
}

func (app *App) postKey(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the postKey!")
}

func (app *App) getKey(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the getKey!")
}
