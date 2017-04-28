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

	user := flag.String("u", "", "one of the configured jenkins user")
	token := flag.String("t", "", "one of the configured jenkins users token")
	server := flag.String("s", "", "the actual jenkins server")
	prefix := flag.String("p", "backend.marketing.jenkins", "prefix to use for graphite")
	graphitePort := flag.Int("gp", 3002, "the port to use to talk to graphite")
	graphiteHost := flag.String("gh", "127.0.0.1", "the server address to use to talk to graphite")

	flag.Parse()

	if *user == "" || *token == "" {
		os.Exit(1)
	}

	con, conErr := net.Dial("tcp", fmt.Sprintf("%s:%d", *graphiteHost, *graphitePort))

	if conErr != nil {
		os.Exit(1)
	}

	response, _ := http.Get(fmt.Sprintf("https://%s:%s@%s/api/json", *user, *token, *server))
	jobs := jenkinsJobResponse{}
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&jobs)
	if err != nil {
		os.Exit(1)
	}

	wg := sync.WaitGroup{}

	for _, currentJob := range jobs.Jobs {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()

			buildResponse, _ := http.Get(fmt.Sprintf("https://%s:%s@%s/job/%s/lastBuild/api/json", *user, *token, *server, name))
			defer buildResponse.Body.Close()

			buildDecoder := json.NewDecoder(buildResponse.Body)
			status := jenkinsBuildResponse{}
			err = buildDecoder.Decode(&status)

			if err != nil {
				return
			}

			value := 1

			if status.Result == "SUCCESS" || status.Result == "building" {
				value = 0
			}

			line := fmt.Sprintf("%s.%s %d %d\n", *prefix, name, value, time.Now().Unix())
			con.Write([]byte(line))

		}(currentJob.Name)
	}

	wg.Wait()

	response.Body.Close()
	con.Close()
}
