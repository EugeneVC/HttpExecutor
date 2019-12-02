package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"models"
	"net/http"
	"os"
)

func SimpleRequest(serverAddress string) error{
	task := models.Task{Method:"GET",URL:"https://ru.wikipedia.org/wiki/%D0%92%D0%B8%D0%BA%D0%B8"}

	body, err := json.Marshal(task)
	if err!=nil{
		return err
	}
	fmt.Println(string(body))

	resp,err := http.Post(fmt.Sprintf("http://%s/tasks",serverAddress),"application/json",bytes.NewBuffer(body))
	if err!=nil{
		return err
	}

	defer resp.Body.Close()

	fmt.Println("Status",resp.StatusCode)

	body,err = ioutil.ReadAll(resp.Body)
	if err!=nil{
		return err
	}

	fmt.Println(string(body))

	return nil
}

func main(){

	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":8887"
	}

	err := SimpleRequest(addr)
	fmt.Println(err)
}