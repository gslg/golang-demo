package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type MyUrl struct {
	ID       string `json:"id,omitempty"`
	LongUrl  string `json:"longUrl,omitempty"`
	ShortUrl string `json:"shortUrl,omitempty"`
}

type Config struct {
	Database struct{
		Address string `json:"address"`
		password string `json:"password,omitempty"`
		DB string `json:"db"`
	} `json:"database"`
	Host string `json:"host"`
	Port string `json:"port"`
}

var appConfig *Config

func LoadConfiguration(filename string) (*Config,error){

	var config Config
	configFile,err := os.Open(filename)

	if err != nil{
		return &config,err
	}

	json.NewDecoder(configFile).Decode(&config)

	return &config,nil
}

var client *redis.Client
var db string

func ExpandEndpoint(w http.ResponseWriter, r *http.Request) {
   params := r.URL.Query()
   shortUrl := params.Get("shortUrl")

	re := client.HGetAll(db).Val()

	for _, v := range re {
		var item MyUrl
		json.Unmarshal([]byte(v), &item)
		if item.ShortUrl == shortUrl {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(401)
	w.Write([]byte("不存在"))
}
func CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	var url MyUrl
	json.NewDecoder(r.Body).Decode(&url)
	re := client.HGetAll(db).Val()

	for _, v := range re {
		var item MyUrl
		json.Unmarshal([]byte(v), &item)
		if item.LongUrl == url.LongUrl {
			w.Write([]byte("已经存在"))
			return
		}
	}

	hd := hashids.NewData()
	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	url.ID, _ = h.Encode([]int{int(now.Unix())})
	url.ShortUrl = strings.Join([]string{"http://",appConfig.Host,":",appConfig.Port,"/",url.ID},"")

	bytes, _ := json.Marshal(url)
	client.HSet(db, url.ID, bytes)

	json.NewEncoder(w).Encode(url)

}
func RootEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result, _ := client.HGet(db, params["id"]).Result()
	var url MyUrl
	json.Unmarshal([]byte(result),&url)

	http.Redirect(w,r,url.LongUrl,301)
}

func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     appConfig.Database.Address,
		Password: appConfig.Database.password, // no password set
		DB:       0,  // use default DB
	})

	return client
}

func main() {
	appConfig,_ = LoadConfiguration("config.json")

	client = GetRedisClient()
	db = appConfig.Database.DB

	router := mux.NewRouter()
	router.HandleFunc("/create",CreateEndpoint).Methods("PUT")
	router.HandleFunc("/expand",ExpandEndpoint).Methods("GET")
	router.HandleFunc("/{id}",RootEndpoint).Methods("GET")

	log.Fatal(http.ListenAndServe(appConfig.Host+":"+appConfig.Port,router))
}
