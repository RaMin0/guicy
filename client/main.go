package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/ramin0/guicy/client/job"
)

var (
	pollURL      = flag.String("poll", "http://localhost:3000", "The base URL to poll")
	pollInterval = flag.Duration("every", 5*time.Second, "The time interval between consequent polls")

	fns = map[string]jobFn{
		"student-data": job.StudentData,
	}
)

type jobFn func(interface{}) (interface{}, error)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	flag.Parse()

	pollTicker := time.NewTicker(*pollInterval)
	for range pollTicker.C {
		log.Printf("Polling %s", *pollURL)
		if err := poll(); err != nil {
			log.Printf("- %v", err)
		}
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

	for _, job := range response.Jobs {
		if fn, ok := fns[job.Type]; ok {
			log.Printf("  * Executing %q as %q", job.ID, job.Type)

			payload, err := fn(job.Request)
			if err != nil {
				log.Printf("    '--> %v", err)
				continue
			}

			request := map[string]interface{}{"payload": payload}
			body, err := json.Marshal(request)
			if err != nil {
				log.Printf("    '--> %v", err)
				continue
			}

			_, err = http.Post(
				*pollURL+"/jobs/"+job.ID,
				"application/json",
				bytes.NewReader(body),
			)
			if err != nil {
				log.Printf("    '--> %v", err)
				continue
			}

			log.Print("    '--> Done")
			continue
		}

		log.Printf("  * Skipping %q as %q", job.ID, job.Type)
	}

	return nil
}
