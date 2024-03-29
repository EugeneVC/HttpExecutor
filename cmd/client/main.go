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

func createTask(serverAddress string, taskURL []string) {

	var wg sync.WaitGroup

	for _, method := range [...]string{"GET", "POST"} {
		for _, val := range taskURL {

			wg.Add(1)

			go func(method, url string) {

				defer wg.Done()

				task := models.Task{Method: method, URL: url}

				bodyRequest, err := json.Marshal(task)
				if err != nil {
					log.Println(task, err)
					return
				}
				//fmt.Println(string(body))

				resp, err := http.Post(fmt.Sprintf("http://%s/task", serverAddress),
					"application/json", bytes.NewBuffer([]byte(bodyRequest)))
				if err != nil {
					log.Println(task, err)
					return
				}

				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					log.Println(task, errors.New(fmt.Sprintf("Status code: %d", resp.StatusCode)))
					return
				}

				err = json.NewDecoder(resp.Body).Decode(&task)

				log.Printf("SUCCESS: %#v", task)
			}(method, val)
		}
	}

	wg.Wait()
}

func getTasks(serverAddress string) {
	resp, err := http.Get(fmt.Sprintf("http://%s/task", serverAddress))
	if err != nil {
		log.Fatal(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Status code: ", resp.StatusCode)
		return
	}

	var tasks []*models.Task
	err = json.NewDecoder(resp.Body).Decode(&tasks)

	log.Printf("Task count: %d", len(tasks))

	for _, task := range tasks {
		log.Printf("%#v", task)
	}
}

func getTasksByPage(serverAddress string, pageNumber, pageSize int) {
	resp, err := http.Get(fmt.Sprintf("http://%s/task?pagenumber=%d&pagesize=%d", serverAddress,
		pageNumber, pageSize))
	if err != nil {
		log.Fatal(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Status code: ", resp.StatusCode)
		return
	}

	var tasks []*models.Task
	err = json.NewDecoder(resp.Body).Decode(&tasks)

	log.Printf("Task count: %d", len(tasks))

	for _, task := range tasks {
		log.Printf("%#v", task)
	}
}

func deleteTask(serverAddress string, key int64) {

	client := &http.Client{}

	request, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/task/%d", serverAddress, key),
		nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("ERROR: Status code: ", resp.StatusCode)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("ADDR")

	taskURL := []string{"https://ru.wikipedia.org/wiki/%D0%92%D0%B8%D0%BA%D0%B8", "https://yandex.ru/",
		"https://gobyexample.com", "https://stackoverflow.com", "https://habr.com", "https://www.google.com/",
		"https://golang.org"}

	log.Println("Tasks CREATE ########################")
	createTask(addr, taskURL)

	log.Println("Tasks LIST ALL  ########################")
	getTasks(addr)

	key := int64(5)
	log.Println(fmt.Sprintf("Tasks DELETE %d ########################", key))
	deleteTask(addr, key)

	log.Println("Tasks LIST ALL ########################")
	getTasks(addr)

	pageNumber, pageSize := 1, 7
	log.Println("Tasks PAGE  ########################")
	getTasksByPage(addr, pageNumber, pageSize)
}
