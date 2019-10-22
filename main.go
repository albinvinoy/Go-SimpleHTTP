package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Person object
type Person struct {
	Name       string
	Age        int
	Profession string
	HairColor  string
}

var peopleMap = make(map[string]Person)

//Get all people info
func PostJson(w http.ResponseWriter, req *http.Request) {
	//fmt.Println("This is a post json data")
	switch req.Method {
	case "GET":
		var names []Person
		for k := range peopleMap {
			names = append(names, (peopleMap[k]))
		}
		val, err := json.Marshal(names)
		if err != nil {
			fmt.Fprintf(w, "Marshal error : %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "%s", val)

	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Println(w, "ParseForm() error : %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		decoder := json.NewDecoder(req.Body)
		var p Person
		err := decoder.Decode(&p)
		if err != nil {
			fmt.Fprintf(w, "Decoder error : %s", err)
			return
		}
		peopleMap[p.Name] = p
		file, err := json.MarshalIndent(peopleMap, "", "\t")
		if err != nil {
			fmt.Fprintf(w, "Json Marshalling Error : %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		err = ioutil.WriteFile("test.json", file, 0644)
		if err != nil {
			fmt.Fprintf(w, "Failed to write file : %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

//return by Name
func getByName(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Path[8:]
	val, err := json.Marshal(peopleMap[name])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", val)
}

func main() {
	http.HandleFunc("/people", PostJson)
	http.HandleFunc("/people/", getByName)
	http.ListenAndServe(":8080", nil)
}
