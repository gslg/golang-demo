package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


type middleWareHandler struct {
	r *mux.Router
	l *ConnLimiter
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.l.GetConn() {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many request"))
		return
	}
	m.r.ServeHTTP(w,r)
	defer m.l.ReleaseConn()
}

func NewMiddleWareHandler(r *mux.Router,limit int) http.Handler{
	return middleWareHandler{r,NewConnLimiter(limit)}
}

func RegisterHandlers() *mux.Router{
	router := mux.NewRouter()

	router.HandleFunc("/videos/{videoId}",streamHandler).Methods("GET")
	router.HandleFunc("/videos",uploadHandler).Methods("POST")


	return router
}

func main(){
	router := RegisterHandlers()
	m := NewMiddleWareHandler(router,2)
	log.Fatal(http.ListenAndServe(":8080", m))
}


