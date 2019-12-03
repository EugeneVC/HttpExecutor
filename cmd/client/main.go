package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"models"
	"net/http"
	"os"
	"sync"
)

func createTask(serverAddress string,taskURL []string){

	var wg sync.WaitGroup

	for _,method := range [...]string{"GET","POST"} {
		for _, val := range taskURL {

			wg.Add(1)

			go func(method,url string) {

				task := models.Task{Method: method, URL: url}

				bodyRequest, err := json.Marshal(task)
				if err != nil {
					log.Println(task,err)
					return
				}
				//fmt.Println(string(body))

				resp, err := http.Post(fmt.Sprintf("http://%s/task", serverAddress), "application/json",bytes.NewBuffer([]byte(bodyRequest)))
				if err != nil {
					log.Println(task,err)
					return
				}

				defer resp.Body.Close()

				if resp.StatusCode!=http.StatusOK{
					log.Println(task,errors.New(fmt.Sprintf("Status code: %d",resp.StatusCode)))
					return
				}

				err = json.NewDecoder(resp.Body).Decode(&task)
				if resp.StatusCode!=http.StatusOK{
					log.Println(task,err)
					return
				}

				log.Printf("SUCCESS: %#v",task)

				wg.Done()
			}(method,val)
		}
	}

	wg.Wait()
}

func main(){

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("ADDR")

	log.Println("Task creation")
	taskURL := []string{"https://ru.wikipedia.org/wiki/%D0%92%D0%B8%D0%BA%D0%B8","https://yandex.ru/",
		"https://gobyexample.com","https://stackoverflow.com", "https://habr.com","https://www.google.com/",
		"https://golang.org"}

	createTask(addr,taskURL)
}