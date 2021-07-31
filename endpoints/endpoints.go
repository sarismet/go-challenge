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
		log.Fatal("Error as conencting to MongoDB")
		return
	}
	app.RedisDB = RedisDB
	app.MongoDB = MongoDB
	http.HandleFunc("/fetchData", app.FetchData)
	http.HandleFunc("/in-memory", app.In_memory)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func (app *App) FetchData(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fetchRequestModel := models.FetchRequestModel{}
	json.Unmarshal([]byte(body), &fetchRequestModel)
	records, msg, resulstCode := app.MongoDB.FetchDataFromMongoDB(fetchRequestModel.StartDate, fetchRequestModel.EndDate, fetchRequestModel.MinCount, fetchRequestModel.MaxCount)
	json, err := json.Marshal(models.FetchResponseModel{Code: resulstCode, Msg: msg, Records: records})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (app *App) In_memory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		key := r.URL.Query().Get("key")
		res, errCode := app.RedisDB.GetKeyFromRedis(key)
		if errCode == 404 {
			json, err := json.Marshal(models.ErrorModel{"Not Found", errCode})
			if err != nil {
				http.Error(w, string(json), http.StatusInternalServerError)
				return
			}
			http.Error(w, string(json), errCode)
			return
		} else if errCode == 500 {
			json, err := json.Marshal(models.ErrorModel{"Internal Server Error", errCode})
			if err != nil {
				http.Error(w, string(json), http.StatusInternalServerError)
				return
			}
			http.Error(w, string(json), errCode)
			return
		}
		getResponseModel := models.GetResponseModel{}
		getResponseModel.Key = key
		getResponseModel.Value = string(res)
		json, err := json.Marshal(getResponseModel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(json)
	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		postRequestModel := models.PostRequestModel{}
		json.Unmarshal([]byte(body), &postRequestModel)
		response := app.RedisDB.InsertKeyToRedis(postRequestModel.Key, postRequestModel.Value)
		postResponseModel := models.PostResponseModel{}
		postResponseModel.Value = response
		json, err := json.Marshal(postResponseModel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(json)
	} else {
		fmt.Fprintf(w, "Method is not allowed!")
	}
}
