package main

import (
	"fmt"
	"net/http"
	"flag"
	"strconv"
	"io/ioutil"
	"os"
	"encoding/json"
	"time"
	"github.com/davecgh/go-spew/spew"
)

func main(){
	fmt.Println("web fish")

	var port int
	var file string
	flag.IntVar(&port, "port", 80, "--port=80")
	flag.StringVar(&file, "file", "", "--file=<path/to/file>")
	flag.Parse()

	if file == ""{
		fmt.Println("missing arg: file")
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){

		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil{
			fmt.Println(err)
			return
		}

		payload := &GitHubPayload{}
		err = json.Unmarshal(body, payload)
		if err != nil{
			fmt.Println(err)
			return
		}

		container := ""

		spew.Dump(payload)

		switch payload.Repository.FullName {
			case "teo-mateo/flbrowser":
				fmt.Println("Adding timestamp for FLBROWSER repo.")
				container = "rtorrent"
			case "teo-mateo/webfish":
				fmt.Println("Adding timestamp for WEBFISH repo.")
				container = "webfish"
			default:
				fmt.Println("Unknown repo: " + payload.Repository.FullName)
		}

		if container != ""{
			err = appendContainerToFile(file, container)
			if err != nil{
				fmt.Println(err)
			}
		}
	})

	var host = "0.0.0.0:"+strconv.Itoa(port)
	fmt.Println("Listening at: " + host)
	fmt.Println("Shared (command) file: " + file)
	http.ListenAndServe(host, nil)
}

func appendContainerToFile(file string, repo string) error {
	now := time.Now()
	s := fmt.Sprintf("%s|%s", now.Format("20060102030405"), repo)

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil{
		return err
	}

	defer f.Close()

	f.WriteString(s)
	return nil
}

type GitHubPayload struct{
	Ref string `json:"ref"`
	Repository struct {
		Name string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
}
