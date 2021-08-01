package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-challenge/db"
	"github.com/go-challenge/models"
)

type App struct {
	RedisDB *db.RedisDatabase
	MongoDB *db.MongoDB
}

//Init method is designed to create new Redis and MongoDB clients.
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

// FetchData is designed to fetch data from MongoDB considering FetchRequestModel fields.
func (app *App) FetchData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fetchRequestModel := models.FetchRequestModel{}
	json.Unmarshal([]byte(body), &fetchRequestModel)
	//we set a date format. The date does not matter. It is set for time.Parse method to understand what date format we want to have
	const (
		layoutISO = "2006-01-02"
	)
	startTime, err := time.Parse(layoutISO, fetchRequestModel.StartDate)
	//If we have an issue parsing the StartDate field then we return an error message to client here.
	if err != nil {
		json, err := json.Marshal(models.FetchResponseModel{Code: 404, Msg: "StartDate is not a valid date"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(404)
		w.Write(json)
		return
	}
	//If we have an issue parsing the EndDate field then we return an error message to client here.
	endTime, err := time.Parse(layoutISO, fetchRequestModel.EndDate)
	if err != nil {
		json, err := json.Marshal(models.FetchResponseModel{Code: 404, Msg: "EndDate is not a valid date"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(404)
		w.Write(json)
		return
	}
	//NOTE: Here wo could have implemented a checking mechanishm for minCount and maxCOunt fields but I though that there is
	// no issue for them to be negative or zero. Even if the client does not set the body properly they would be set zero and it is okey.
	records, msg, resulstCode := app.MongoDB.FetchDataFromMongoDB(startTime, endTime, fetchRequestModel.MinCount, fetchRequestModel.MaxCount)
	json, err := json.Marshal(models.FetchResponseModel{Code: resulstCode, Msg: msg, Records: records})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func (app *App) In_memory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		key := r.URL.Query().Get("key")
		//if key parameter is not found then we return an error message to client
		if key == "" {
			json, err := json.Marshal(models.ErrorModel{"parameter named key is not set properly", 404})
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(404)
			w.Write(json)
			return
		}
		res, errCode := app.RedisDB.GetKeyFromRedis(key)
		//if we encounter an issue as getting from redis then we response not found
		if errCode == 404 {
			json, err := json.Marshal(models.ErrorModel{"Not Found", errCode})
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(errCode)
			w.Write(json)
			return
		} else if errCode == 500 {
			//if we encounter an issue as getting from redis then we response internal server error
			//Here we do not care about error which comes from Marshall method since the error code would be the same.
			json, _ := json.Marshal(models.ErrorModel{"Internal Server Error", errCode})
			w.WriteHeader(errCode)
			w.Write(json)
			return
		} else {
			getResponseModel := models.GetResponseModel{}
			getResponseModel.Key = key
			getResponseModel.Value = string(res)
			json, err := json.Marshal(getResponseModel)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(json)
		}
	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		postRequestModel := models.PostRequestModel{}
		json.Unmarshal([]byte(body), &postRequestModel)
		//here we check if the key field is set properly
		if postRequestModel.Key == "" {
			json, err := json.Marshal(models.ErrorModel{"Key is not set properly", 404})
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(404)
			w.Write(json)
			return
		} else if postRequestModel.Value == "" {
			//here we check if the value field is set properly
			json, err := json.Marshal(models.ErrorModel{"Value is not set properly", 404})
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(404)
			w.Write(json)
			return
		}
		// we insert the key value pair in InsertKeyToRedis and get the reponse from it
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
		//if there is another method type other than GET and POST then we just simply return Method is not allowed!
		json, err := json.Marshal(models.ErrorModel{"Method is not allowed!", 404})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(404)
		w.Write(json)
	}
}
