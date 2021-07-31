package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-challenge/db"
	"github.com/go-challenge/endpoints"
	"github.com/go-challenge/models"
	"github.com/stretchr/testify/assert"
)

type App struct {
	RedisDB *db.RedisDatabase
	MongoDB *db.MongoDB
}

func TestFetchData(t *testing.T) {
	app := endpoints.App{}
	MongoDB, err := db.ConnectMongoDB()
	if err != nil || MongoDB == nil {
		log.Fatal("Error as conencting to MongoDB")
		return
	}
	app.MongoDB = MongoDB
	bodyJson := `{"startDate":"2016-11-26","endDate":"2016-12-05","minCount":100,"maxCount":150}`
	request, _ := http.NewRequest(http.MethodGet, "/fetchData", bytes.NewBuffer([]byte(bodyJson)))
	response := httptest.NewRecorder()
	app.FetchData(response, request)
	fetchResponseModel := models.FetchResponseModel{}
	json.Unmarshal(response.Body.Bytes(), &fetchResponseModel)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, fetchResponseModel.Msg, "Success")
	assert.Equal(t, len(fetchResponseModel.Records), 5)
	assert.Equal(t, Contains(fetchResponseModel.Records, "ZrUxelLG"), true)
}

func Contains(records []models.FetchRecordsArrayModel, key string) bool {
	result := false
	for _, item := range records {
		if item.Key == key {
			result = true
			break
		}
	}
	return result
}
