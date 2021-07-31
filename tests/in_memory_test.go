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
	bodyJson := `{"key": "active-tabs","value":"getir"}`
	requestPost, _ := http.NewRequest(http.MethodPost, "http://localhost:5000/in-memory", bytes.NewBuffer([]byte(bodyJson)))
	responsePost := httptest.NewRecorder()
	app.In_memory(responsePost, requestPost)
	postResponseModel := models.PostResponseModel{}
	json.Unmarshal(responsePost.Body.Bytes(), &postResponseModel)
	assert.Equal(t, responsePost.Code, 200)
	assert.Equal(t, postResponseModel.Value, "OK")
	requestGet, _ := http.NewRequest(http.MethodGet, "http://localhost:5000/in-memory?key=active-tabs", nil)
	responseGet := httptest.NewRecorder()
	app.In_memory(responseGet, requestGet)
	getResponseModel := models.GetResponseModel{}
	json.Unmarshal(responseGet.Body.Bytes(), &getResponseModel)
	assert.Equal(t, responseGet.Code, 200)
	assert.Equal(t, getResponseModel.Key, "active-tabs")
	assert.Equal(t, getResponseModel.Value, "getir")
}
