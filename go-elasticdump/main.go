package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	initurl = "http://localhost:9200/fazhi_company/_search?scroll=1m"
	url = "http://localhost:9200/_search/scroll"
	query = `{"size":1000,"query": {"match_all" : {}}}`
)


type scroll struct {
	scrollId string
	datas []interface{}
}

func getScroll(url string,query []byte,ch chan interface{}) scroll {

	req,err := http.NewRequest("POST",url,bytes.NewBuffer(query))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp,err :=client.Do(req)

	if err != nil{
		log.Fatalf("Error geeting response: %s",err)
	}

	defer resp.Body.Close()

	body,_:= ioutil.ReadAll(resp.Body)

	log.Printf("请求返回%s",body)
	var result map[string]interface{}
	json.Unmarshal([]byte(body),&result)

	scrollId := result["_scroll_id"].(string)

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})

	datas := make([]interface{},0)

	if len(hits)>0 {
		for _,value := range hits {
			v := value.(map[string]interface{})
			source := v["_source"]
			ch <- source
			datas = append(datas,source)
		}
	}

    return scroll{
    	scrollId:scrollId,
    	datas:datas,
	}
}

type Query struct {
	Scroll string `json:"scroll"`
	ScrollId string `json:"scroll_id"`
}

func SaveJsonFile(datas []interface{},sequence int){
	log.Println("开始保存文件........")
	filename := "./data2/data"+strconv.Itoa(sequence)+".json"

	f,err := os.Create(filename)

	defer f.Close()

	if err != nil {
		fmt.Println("os Create error: ", err)
		return
	}

	bw := bufio.NewWriter(f)

	for _,v := range datas {
		r, _ := json.Marshal(v)
		bw.WriteString(string(r) + "\n")
	}
	bw.Flush()
}

func AppendToFile(filename string,ch chan interface{}){
	f,err := os.OpenFile(filename,os.O_RDWR|os.O_CREATE|os.O_APPEND,0777)
	if err != nil {
		fmt.Println("os OpenFile error: ", err)
		return
	}
	defer f.Close()

	for data := range ch {
		log.Printf("收到数据:%s",data)
		r,err := json.Marshal(data)
		if err!=nil {
			log.Fatalf("转json字符串出错:%s",err)
		}

		f.WriteString(string(r)+"\n")
	}
}

func getData(ch chan interface{}){

	var jsonStr = []byte(query)
	reslut := getScroll(initurl,jsonStr,ch)
	var sequence = 0
	for ;len(reslut.datas) > 0 ;{

		scrollQuery := Query{
			Scroll:"1m",
			ScrollId:reslut.scrollId,
		}
		log.Println(reslut.scrollId)
		sequence += 1

		/*for index,value := range reslut.datas {

			log.Printf("index=%d,value=%s",index,value)

		}*/
		q,err := json.Marshal(scrollQuery)
		if err != nil{
			log.Fatalf("Error %s",err)
		}
		log.Println("请求参数:"+ string(q))
		reslut = getScroll(url,q,ch)
	}
	close(ch)
}

func main() {
	ch := make(chan interface{},1000)

	go getData(ch)

	AppendToFile("./datas/data2.json",ch)
	//var sequence = 0
	/*for data := range ch {
		sequence += 1
		log.Printf("收到数据:%s",data)
		SaveJsonFile(data,sequence)
	}*/



}