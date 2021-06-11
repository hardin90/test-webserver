package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"log"
	"net/http"
	"regexp"
	"encoding/json"
)

func SetMyCookie(response http.ResponseWriter){
	cookie := http.Cookie{Name: "cookiename", Value:"cookievalue"}
	http.SetCookie(response, &cookie)
}


func GenericHandler(response http.ResponseWriter, request *http.Request){

	SetMyCookie(response)
	response.Header().Set("Content-type", "text/plain")

	err := request.ParseForm()
	if err != nil {
		http.Error(response, fmt.Sprintf("error parsing url %v", err), 500)
	}

	fmt.Fprint(response,  "FooWebHandler says ... \n")
	fmt.Fprintf(response, " request.Method     '%v'\n", request.Method)
	fmt.Fprintf(response, " request.RequestURI '%v'\n", request.RequestURI)
	fmt.Fprintf(response, " request.URL.Path   '%v'\n", request.URL.Path)
	fmt.Fprintf(response, " request.Form       '%v'\n", request.Form)
	fmt.Fprintf(response, " request.Cookies()  '%v'\n", request.Cookies())
}

func HomeHandler(response http.ResponseWriter, request *http.Request){
	response.Header().Set("Content-type", "text/html")
	webpage, err := ioutil.ReadFile("home.html")
	if err != nil { 
		http.Error(response, fmt.Sprintf("home.html file error %v", err), 500)
	}
	fmt.Fprint(response, string(webpage));
}

func ItemHandler(response http.ResponseWriter, request *http.Request){

	SetMyCookie(response)
	response.Header().Set("Content-type", "application/json")

	data := map[string]string { "what" : "item", "name" : "" }

	var itemURL = regexp.MustCompile(`^/item/(\w+)$`)
	var itemMatches = itemURL.FindStringSubmatch(request.URL.Path)
	if len(itemMatches) > 0 {
		data["name"] = itemMatches[1] 
		json_bytes, _ := json.Marshal(data)
		fmt.Fprintf(response, "%s\n", json_bytes)

	} else {
		http.Error(response, "404 page not found", 404)
	}
}

func main(){
	port := 8098
	portstring := strconv.Itoa(port)

	mux := http.NewServeMux()
	mux.Handle("/home", http.HandlerFunc( HomeHandler ))
	mux.Handle("/item/", http.HandlerFunc( ItemHandler ))
	mux.Handle("/generic/", http.HandlerFunc( GenericHandler ))

	log.Print("Listening on port " + portstring + " ... ")
	err := http.ListenAndServe(":" + portstring, mux)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

