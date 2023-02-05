package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	var wg sync.WaitGroup
	var linkIds = []int{1, 2, 3, 4, 5}
	var todos = map[int]Todo{}
	todoChan := make(chan Todo, len(linkIds))

	for _, i := range linkIds {
		wg.Add(1)
		go fetchTodos(&wg, todoChan, i)
	}

	go func() {
		wg.Wait()
		close(todoChan)
	}()

	for todo := range todoChan {
		todos[todo.ID] = todo
	}

	// Pretty printing
	b, err := json.Marshal(todos)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func fetchTodos(wg *sync.WaitGroup, todoChan chan Todo, id int) {
	defer wg.Done()
	response := Todo{}
	link := fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", id)
	fmt.Println("Fetching todo with id ", id)

	req, err := http.NewRequest(http.MethodGet, link, &bytes.Buffer{})
	if err != nil {
		fmt.Printf("Error in making the http request %s", err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("error in doing the http request %s", err.Error())
		return
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error in reading the response body %s ", err.Error())
	}
	fmt.Println("Received response for ", id)

	switch res.StatusCode {
	case http.StatusOK:
		err := json.Unmarshal(responseBody, &response)
		if err != nil {
			fmt.Printf("error in reading the unmarshalling body %s ", err.Error())
		}
		todoChan <- response
	default:
		fmt.Printf("Unexpected code ,  Status code %d| Response %s", res.StatusCode, string(responseBody))
	}

}
