package main

import (
	"fmt"
	"net/http"
	"flag"
	"strconv"
	"io/ioutil"
	"github.com/davecgh/go-spew/spew"
)

func main(){
	fmt.Println("web fish")

	var port int
	flag.IntVar(&port, "port", 80, "--port=80")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){

		defer r.Body.Close()
		spew.Dump(r.Header)

		body, err := ioutil.ReadAll(r.Body)
		if err == nil{
			fmt.Println("BODY:")
			fmt.Println(string(body))
		}

		w.Write([]byte("Yes webfish"))
	})
	var host = "localhost:"+strconv.Itoa(port)
	fmt.Println("Listening at " + host)

	http.ListenAndServe(host, nil)

}