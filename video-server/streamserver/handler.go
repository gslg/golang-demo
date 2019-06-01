package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func streamHandler(w http.ResponseWriter, r *http.Request) {
	vid := mux.Vars(r)
	videoLocation := VIDEO_DIR + vid["videoId"]

	video, err := os.Open(videoLocation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("播放视频出错" + err.Error()))
		log.Printf("读取文件报错:%v",err)
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "视频", time.Now(), video)

	defer video.Close()
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte("文件太大了........"))
		return
	}

	file, _, err := r.FormFile("file")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("读取文件报错!"))
		log.Printf("读取文件报错:%v",err)
		return
	}

	data,err := ioutil.ReadAll(file)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("读取文件报错!"))
		log.Printf("读取文件报错:%v",err)
		return
	}

	name := r.FormValue("name")

	err = ioutil.WriteFile(VIDEO_DIR+name,data,0666)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("上传文件报错!"))
		log.Printf("上传文件报错:%v",err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("上传成功!"))
}
