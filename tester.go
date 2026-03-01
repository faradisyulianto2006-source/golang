package main

// This file contains a tester function that sends a GET request to the /projects endpoint of the server and prints the response. It is used for testing the functionality of the server and can be called from the main function or any other part of the code to verify that the server is working as expected.
import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func tester () {
	// define the URL for the GET request to the /projects endpoint
	url := "http://localhost:8080/projects"

	// create a new HTTP request with the specified URL and method
	req, _ := http.NewRequest("GET", url, nil)

	// add necessary headers to the request, such as authorization and cache control
	req.Header.Add("Authorization", "Basic YWRtaW46YWRtaW4=")
	req.Header.Add("Cache-control", "no-cache")

	// send the HTTP request and receive the response
	res, _ := http.DefaultClient.Do(req)

	// read the response body and print the status code and body content
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res)
	fmt.Println(string(body))
}