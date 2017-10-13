package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
	"regexp"
)

type job struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type jenkinsJobResponse struct {
	Jobs []job `json:"jobs"`
}

type jenkinsBuildResponse struct {
	Result string `json:"result"`
}

func main() {
	var user, token, server, prefix, filter, graphiteHost string
	var graphitePort int
	flag.StringVar(&user, "u", "", "one of the configured jenkins user")
	flag.StringVar(&token, "t", "", "one of the configured jenkins users token")
	flag.StringVar(&server, "s", "", "the actual jenkins server")
	flag.StringVar(&prefix, "p", "", "prefix to use for graphite. e.g., backend.marketing.jenkins")
	flag.StringVar(&filter, "f", ".*", "regular expression to filter by job names")
	flag.IntVar(&graphitePort, "gp", 3002, "the port to use to talk to graphite")
	flag.StringVar(&graphiteHost, "gh", "127.0.0.1", "the server address to use to talk to graphite")

	flag.Parse()

	filterRE, err := regexp.Compile(filter)
	if err != nil {
		fmt.Printf("Invalid filter expression: %s\n", err)
		os.Exit(1)
	}

	if user == "" || token == "" || prefix == "" || server == "" {
		fmt.Println("Please specify jenkins user (-u), token (-t), server (-s) and prefix (-p)")
		os.Exit(1)
	}

	connectionAddress := fmt.Sprintf("%s:%d", graphiteHost, graphitePort)
	con, err := net.Dial("tcp", connectionAddress)
	if err != nil {
		fmt.Println("Can not connect to", connectionAddress)
		os.Exit(1)
	}
	defer con.Close()

	response, err := http.Get(fmt.Sprintf("https://%s:%s@%s/api/json", user, token, server))
	if err != nil {
		fmt.Printf("Failed to download json from https://%s/api/json\n", server)
		os.Exit(1)
	}
	defer response.Body.Close()

	jobs := jenkinsJobResponse{}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jobs)
	if err != nil {
		fmt.Println("Failed to decode json")
		os.Exit(1)
	}

	wg := sync.WaitGroup{}

	for _, currentJob := range jobs.Jobs {
		if !filterRE.MatchString(currentJob.Name) {
			continue
		}

		wg.Add(1)
		go func(name string) {
			defer wg.Done()

			buildResponse, err := http.Get(fmt.Sprintf("https://%s:%s@%s/job/%s/lastBuild/api/json", user, token, server, name))
			if err != nil {
				fmt.Println("Failed to get job data for", name)
				return
			}
			defer buildResponse.Body.Close()

			buildDecoder := json.NewDecoder(buildResponse.Body)
			status := jenkinsBuildResponse{}
			err = buildDecoder.Decode(&status)
			if err != nil {
				fmt.Println("Failed to decode json for", name)
				return
			}

			value := 1
			if status.Result == "SUCCESS" || status.Result == "building" {
				value = 0
			}

			con.Write([]byte(fmt.Sprintf("%s.%s %d %d\n", prefix, name, value, time.Now().Unix())))
		}(currentJob.Name)
	}

	wg.Wait()
}
