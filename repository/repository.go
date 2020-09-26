package repository

import (
	"container/list"
	"database/sql"
	"encoding/json"
	"fmt"
	"golang-api/cache"
	"golang-api/db"

	"github.com/go-redis/redis/v7"
	"github.com/rs/zerolog/log"
)

//Repo - abstract repository interface
type Repo interface {
	Select(doc interface{}) (*list.List, error)
	Insert(doc interface{}) (int64, error)
	Update(doc interface{}) (int64, error)
	Remove(doc interface{}) (int64, error)
}

//Data struct for  db
type Data struct {
	Id   int
	Name string
}

//DbConnection - Database Connectin Pool
var DbConnection *sql.DB

//RedisClient for redis client
var RedisClient *redis.Client

//SetupRepo - setup database connections
func SetupRepo() (err error) {
	DbConnection, err = db.GetDatabase()
	_, err = DbConnection.Exec("CREATE DATABASE testDB")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully created database..")
	}
	_, err = DbConnection.Exec("CREATE TABLE IF NOT EXISTS mytable (id INT NOT NULL PRIMARY KEY, name TEXT NOT NULL)")
	if err != nil {
		panic(err)
	}
	RedisClient = cache.InitRedisConfig("localhost:6379")
	return
}

//CloseRepo - close database connections
func CloseRepo() {
	if DbConnection != nil {
		defer DbConnection.Close()
	}
}

// ReadData from cache/db
func ReadData(value []string, page int) []string {
	var record []string
	record = getFromCache(value)
	if len(record) < 1 {
		rows, err := DbConnection.Query("select * from mytable limit ?", page)
		if err != nil {
			panic(err)
		}
		data := Data{}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&data.Id, &data.Name)
			if err != nil {
				panic(err)
			}
			fmt.Println(data)
		}
		cacheAddData(data)
		record = getFromCache(value)
	}
	return record
}

// InsertData into db
func InsertData(data []byte) bool {
	datalist := &[]Data{}
	if err := json.Unmarshal([]byte(data), datalist); err != nil {
		log.Error().Err(err).Msg("Marshaling Error prevTagsFromCache")
	}
	for _, d := range *datalist {
		_, err := DbConnection.Exec("INSERT INTO mytable name VALUES (?)", d.Id, d.Name)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}
	log.Info().Msg("Data added to db")
	return true
}

func cacheAddData(data Data) {
	if err := RedisClient.HMSet("data", data).Err(); err != nil {
		log.Error().Err(err).Msg("Unable to store tags into redis")
	}
}

func getFromCache(value []string) []string {
	var record []string
	data, err := RedisClient.HMGet("data", value...).Result()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching data from Redis")
	}
	if len(data) > 0 {
		for _, t := range data {
			record = append(record, t.(string))
		}
	}
	return record
}

//UpdateCache data when we get notification from kafka
func UpdateCache() {
	rows, err := DbConnection.Query("select * from mytable")
	if err != nil {
		panic(err)
	}
	data := Data{}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&data.Id, &data.Name)
		if err != nil {
			panic(err)
		}
	}
	cacheAddData(data)
	log.Info().Msg("Kafka notification:: Cache updated")
}
