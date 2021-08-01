## Run Locally

Clone the project in your go path

```bash
  cd ~/go/src/github.com/
  git clone https://github.com/sarismet/go-challenge
```

Go to the project directory

```bash
  cd go-challenge
```
Make sure 

```
  GO111MODULE="auto"
```
To make GO111MODULE="auto" run

```
  export GO111MODULE="auto"
```

Install dependencies

```bash
  go get github.com/go-redis/redis
  go get github.com/alicebob/miniredis
  go get github.com/stretchr/testify
  go get get go.mongodb.org/mongo-driver/mongo
  go get get go.mongodb.org/mongo-driver/bson
```

Set MongoDB url in docker-compose.yml file

```
environment:
      REDIS_URL: redis:6379
      MONGO_URL: [MongoDB_URL]
```

Please make sure that you have docker and docker compose in your machine.  
After that run these commands.
```
  docker-compose up -d
```

## Tech Stack

**Database:** [Redis](https://redis.io/)
```
I needed to choose a In-Memory database and I decided to use Redis since I beleive that it is efficient and easy to use.
```

## API Documentation

#### -Fetch Data from MongoDB

Request  

```http
  GET http://localhost:5000/fetchData or http://18.118.255.222:5000/fetchData
```

| Body | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `{"startDate": string,"endDate": string,"minCount": int,"maxCount": int}` | `json` | Fetch Data request model|

- Example
```
{
    "startDate": "2016-11-22",
    "endDate": "2016-12-05",
    "minCount": 100,
    "maxCount": 150
}
```
Response  

| Body | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `code` | `string` | response code|
| `msg` | `string` | response msg|
| `"records": [{"key": string,"createdAt": string,"totalCounts": int}]` | `array` | the records from mongodb|

- Example
```
{
    "code": 0,
    "msg": "Success",
    "records": [
        {
            "key": "ZrUxelLG",
            "createdAt": "2016-12-03T13:04:17.799Z",
            "totalCounts": 103
        },
        {
            "key": "JJlGewEB",
            "createdAt": "2016-12-03T06:01:50.089Z",
            "totalCounts": 115
        },
        {
            "key": "enaLTnTM",
            "createdAt": "2016-11-30T05:18:04.161Z",
            "totalCounts": 111
        }
    ]
}
```

#### In-MemoryDB Endpoints

##### POST

Request  

```http
  POST http://localhost:5000/in-memory or http://18.118.255.222:5000/in-memory
```

| Body | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `{"key": string,"value": string}` | `json` | json model which includes key and value|

- Example
```
{
    "key": "key1",
    "value": "value2"
}
```

Response  
| Body | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `"value": string` | `string` |post response value|

- Example
```
{
    "value": string
}
```

##### GET

Request  

```http
  GET http://localhost:5000/in-memory?key=key1 or http://18.118.255.222:5000/in-memory?key=key1
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `key` | `string` | key from which we get value via redis|

Response  

| Body | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `{"key": string,"value": string}` | `json` | json model which includes key and value|

- Example
```
{
    "key": "key1",
    "value": "value2"
}
```
## Note
- I added a postman collection json file. You can just basicly import to your postman so that you can see how I had tried endpoints.

## Running Tests

To run the endpoints one by one
```bash
  cd tests
  go test ./fetchData_test.go -v
  go test ./in_memory_test.go -v
```


