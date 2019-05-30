package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"strings"
)

type Person struct {
	ID string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname string `json:"lastname,omitempty"`
	Social []SocialMedia `json:"socialmedia,omitempty"`
}

type SocialMedia struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

const (
	db = "people"
)

var client *redis.Client

func GetPeopleEndpoint(w http.ResponseWriter,req *http.Request){
	all,_ := client.HGetAll(db).Result()
	var people []string
	for _,v := range all {
		people = append(people,v)
	}
	var p []Person

	json.Unmarshal([]byte("["+strings.Join(people,",")+"]"),&p)
	json.NewEncoder(w).Encode(p)
}

func GetPersonEndpoint(w http.ResponseWriter,req *http.Request){
    params := mux.Vars(req)
	p,_ := client.HGet(db, params["id"]).Bytes()
	var person Person
	json.Unmarshal(p,&person)
	json.NewEncoder(w).Encode(person)
}

func CreatePersonEndpoint(w http.ResponseWriter,req *http.Request){
	var person Person
	json.NewDecoder(req.Body).Decode(&person)
	if person.ID == "" {
		person.ID = uuid.Must(uuid.NewV4()).String()
	}

	bytes, _ := json.Marshal(person)
	err := client.HSet(db, person.ID, bytes).Err()

	if err != nil{
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(person.ID)
}

func UpdatePersonEndpoint(w http.ResponseWriter,req *http.Request){
	params := mux.Vars(req)

	data,err := client.HGet(db,params["id"]).Bytes()
	if err != nil{
		w.WriteHeader(404)
		if len(data) == 0 {
			w.Write([]byte("用户不存在"))
			return
		}
		w.Write([]byte(err.Error()))
		return
	}

	var person,old Person
	json.NewDecoder(req.Body).Decode(&person)

	json.Unmarshal(data,&old)
	person.ID = old.ID

	bytes, _ := json.Marshal(person)
	client.HSet(db,person.ID,bytes)

	json.NewEncoder(w).Encode(&person)
}
func DeletePersonEndpoint(w http.ResponseWriter,req *http.Request){
	params := mux.Vars(req)
	_, e := client.HDel(db, params["id"]).Result()

	if e != nil{
		w.WriteHeader(401)
		w.Write([]byte(e.Error()))
		return
	}

	json.NewEncoder(w).Encode("删除成功!")
}

func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}


func ExampleNewClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	person := Person{
		ID:        "1",
		Firstname: "刘国",
		Lastname:  "李菲蓉",
		Social: []SocialMedia{
			{Title: "Github", Link: "http://gslg.github.com"},
			{Title: "Twitter", Link: "http://www.twitter.com/liuguo"},
		},
	}
	data, _ := json.Marshal(person)
	err := client.HSet("people", "1", data).Err()

	if err!=nil{
		log.Fatalf("Put Error:%s",err)
	}

	res,err := client.HGet("people","1").Result()


	if err!=nil{
		log.Fatalf("Get Error:%s",err)
	}

	var p Person
	json.Unmarshal([]byte(res), &p)

	log.Println(p)



}


func main() {

	client=GetRedisClient()

	router := mux.NewRouter()
	router.HandleFunc("/people",GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}",GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/person",CreatePersonEndpoint).Methods("PUT")
	router.HandleFunc("/person/{id}",UpdatePersonEndpoint).Methods("POST")
	router.HandleFunc("/person/{id}",DeletePersonEndpoint).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":12345",router))
	//ExampleNewClient()
}
