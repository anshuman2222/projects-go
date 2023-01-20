package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

type User struct {
	Name string
	UID  string
}

type UserJson struct {
	Users []struct {
		Name string `json:"name"`
		UID  string `json:"uid"`
	}
}

type Stack struct {
	Cap int
	Arr []interface{}
}

var stack Stack

var users map[string]*User

var push_ch = make(chan interface{})
var pop_ch = make(chan interface{})
var is_pop = make(chan bool)

func handle_user_create_oper(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering user create handle function ")

	payload, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal("Failed in getting paylaod, ", err)
		return
	}
	var user_json UserJson

	err = json.Unmarshal(payload, &user_json)

	if err != nil {
		log.Fatal("Failed in unmarshling data, err: ", err)
		return
	}
	for i := 0; i < len(user_json.Users); i++ {
		user := new(User)
		user.Name = user_json.Users[i].Name
		user.UID = user_json.Users[i].UID

		fmt.Printf("User is getting add: %v", user)
		push_ch <- user
	}
}

func handle_user_get_oper(w http.ResponseWriter, r *http.Request) {

	user_id := strings.Split(r.RequestURI, "/")[4]
	fmt.Println("User id: ", user_id)

	//user_data := make(map[string]*User)

	is_pop <- true
	user_stack_data := <-pop_ch

	//user_data[user_id] = users[user_id]

	byte_data, err := json.Marshal(user_stack_data)

	w.Header().Set("content-type", "application/json")

	if err != nil {
		log.Fatal("Failed in marshling data. err: ", err)
		return
	}
	w.Write(byte_data)
}

func (s *Stack) push(ch chan interface{}, m *sync.Mutex) {
	for {
		select {
		case data := <-ch:
			if s.Cap == len(s.Arr) {
				fmt.Println("Stack is overflow")
				return
			} else {
				m.Lock()
				s.Arr = append(s.Arr, data)
				m.Unlock()
			}
		}
	}
}

func (s *Stack) pop(is_pop chan bool, pop_data chan interface{}, m *sync.Mutex) {
	for {
		select {
		case <-is_pop:
			if len(s.Arr) == 0 {
				fmt.Println("Stack is underflow")
				return
			} else {
				m.Lock()
				pop_data <- s.Arr[len(s.Arr)-1]
				s.Arr = s.Arr[0 : len(s.Arr)-1]
				m.Unlock()
			}
		}
	}
}

func init() {
	users = make(map[string]*User)
	stack.Cap = 10
	stack.Arr = make([]interface{}, 0)
}

func main() {
	var wg sync.WaitGroup
	var m sync.Mutex

	wg.Add(2)

	go stack.push(push_ch, &m)
	go stack.pop(is_pop, pop_ch, &m)

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/users", handle_user_create_oper).Methods("POST")
	r.HandleFunc("/api/v1/users/{id}", handle_user_get_oper).Methods("GET")

	fmt.Println("listening on : 127.0.0.1:8000")
	http.ListenAndServe("127.0.0.1:8000", r)
}

