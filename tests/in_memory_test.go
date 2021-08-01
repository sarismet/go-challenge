package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-challenge/db"
	"github.com/go-challenge/endpoints"
	"github.com/go-challenge/models"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestIn_memory(t *testing.T) {
	app := endpoints.App{}
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	newRedisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	if newRedisClient == nil {
		log.Fatalf("RedisDB is nil")
	}
	RedisDB := &db.RedisDatabase{
		Client: newRedisClient,
	}
	app.RedisDB = RedisDB
	//Since it is not set yet it response "not found" with 404 error code
	requestGet, _ := http.NewRequest(http.MethodGet, "http://localhost:5000/in-memory?key=key1", nil)
	responseGet := httptest.NewRecorder()
	app.In_memory(responseGet, requestGet)
	errorModel := models.ErrorModel{}
	json.Unmarshal(responseGet.Body.Bytes(), &errorModel)
	assert.Equal(t, responseGet.Code, 404)
	assert.Equal(t, errorModel.Msg, "Not Found")
	assert.Equal(t, errorModel.Code, 404)
	//Since it is not set yet it response "not found" with 404 error code
	requestGet, _ = http.NewRequest(http.MethodGet, "http://localhost:5000/in-memory?keyNotValid=key1", nil)
	responseGet = httptest.NewRecorder()
	app.In_memory(responseGet, requestGet)
	errorModel = models.ErrorModel{}
	json.Unmarshal(responseGet.Body.Bytes(), &errorModel)
	assert.Equal(t, responseGet.Code, 404)
	assert.Equal(t, errorModel.Msg, "parameter named key is not set properly")
	assert.Equal(t, errorModel.Code, 404)
	//we do not set value field and expect that it response an error
	bodyJson := `{"key": "key1"}`
	requestPost, _ := http.NewRequest(http.MethodPost, "http://localhost:5000/in-memory", bytes.NewBuffer([]byte(bodyJson)))
	responsePost := httptest.NewRecorder()
	app.In_memory(responsePost, requestPost)
	errorModel = models.ErrorModel{}
	json.Unmarshal(responsePost.Body.Bytes(), &errorModel)
	assert.Equal(t, responsePost.Code, 404)
	assert.Equal(t, errorModel.Msg, "Value is not set properly")
	assert.Equal(t, errorModel.Code, 404)
	//we do not set key field and expect that it response an error
	bodyJson = `{"value": "value1"}`
	requestPost, _ = http.NewRequest(http.MethodPost, "http://localhost:5000/in-memory", bytes.NewBuffer([]byte(bodyJson)))
	responsePost = httptest.NewRecorder()
	app.In_memory(responsePost, requestPost)
	errorModel = models.ErrorModel{}
	json.Unmarshal(responsePost.Body.Bytes(), &errorModel)
	assert.Equal(t, responsePost.Code, 404)
	assert.Equal(t, errorModel.Msg, "Key is not set properly")
	assert.Equal(t, errorModel.Code, 404)
	//we construct a valid post request valid body
	bodyJson = `{"key": "key1","value":"value1"}`
	requestPost, _ = http.NewRequest(http.MethodPost, "http://localhost:5000/in-memory", bytes.NewBuffer([]byte(bodyJson)))
	responsePost = httptest.NewRecorder()
	app.In_memory(responsePost, requestPost)
	postResponseModel := models.PostResponseModel{}
	json.Unmarshal(responsePost.Body.Bytes(), &postResponseModel)
	assert.Equal(t, responsePost.Code, 200)
	assert.Equal(t, postResponseModel.Value, "OK")
	//we construct a valid get request valid parameter
	requestGet, _ = http.NewRequest(http.MethodGet, "http://localhost:5000/in-memory?key=key1", nil)
	responseGet = httptest.NewRecorder()
	app.In_memory(responseGet, requestGet)
	getResponseModel := models.GetResponseModel{}
	json.Unmarshal(responseGet.Body.Bytes(), &getResponseModel)
	assert.Equal(t, responseGet.Code, 200)
	assert.Equal(t, getResponseModel.Key, "key1")
	assert.Equal(t, getResponseModel.Value, "value1")
}
