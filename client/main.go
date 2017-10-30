package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ramin0/guicy/client/job"
)

var (
	pollURL      = flag.String("poll", "http://localhost:3000", "The base URL to poll")
	pollInterval = flag.Duration("every", 5*time.Second, "The time interval between consequent polls")

	fns = []jobFn{
		job.StudentData,
		job.SendNotification,
	}
)

type jobFn interface {
	Name() string
	Description() string
	Inputs() []map[string]string
	Outputs() []map[string]string
	Exec(interface{}) (interface{}, error)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	flag.Parse()

	for {
		log.Printf("Polling %s", *pollURL)
		poll()
		time.Sleep(*pollInterval)
	}
}

func poll() error {
	res, err := http.Get(*pollURL + "/jobs")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var response struct {
		Jobs []struct {
			ID      string
			Type    string
			Request interface{}
		}
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return err
	}

	if l := len(response.Jobs); l == 0 {
		log.Print("- No jobs found")
	} else {
		log.Printf("- Received %d job(s):", l)
	}

jobsLoop:
	for _, j := range response.Jobs {
		for _, fn := range fns {
			if j.Type == "discover" || j.Type == fmt.Sprint(fn) {
				log.Printf("  * Executing %q as %q", j.ID, j.Type)

				var (
					payload interface{}
					err     error
				)
				if j.Type == "discover" {
					payload = discoverPayload()
				} else {
					payload, err = fn.Exec(j.Request)
				}
				if err != nil {
					log.Printf("    '--> %v", err)
					continue jobsLoop
				}

				request := map[string]interface{}{"payload": payload}
				body, err := json.Marshal(request)
				if err != nil {
					log.Printf("    '--> %v", err)
					continue jobsLoop
				}

				_, err = http.Post(
					*pollURL+"/jobs/"+j.ID,
					"application/json",
					bytes.NewReader(body),
				)
				if err != nil {
					log.Printf("    '--> %v", err)
					continue jobsLoop
				}

				log.Print("    '--> Done")
				continue jobsLoop
			}
		}

		log.Printf("  * Skipping %q as %q", j.ID, j.Type)
	}

	return nil
}

func discoverPayload() interface{} {
	jobs := make([]map[string]interface{}, 0, len(fns))
	for _, fn := range fns {
		jobs = append(jobs, map[string]interface{}{
			"id":          fmt.Sprint(fn),
			"name":        fn.Name(),
			"description": fn.Description(),
			"inputs":      fn.Inputs(),
			"outputs":     fn.Outputs(),
		})
	}
	return jobs
}
