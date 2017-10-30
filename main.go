package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	jobs = map[string]*Job{}
)

// Job struct
type Job struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"`
	Request  interface{} `json:"request,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/discover", withLog(handleDiscover))
	http.HandleFunc("/requests", withLog(handleRequests))
	http.HandleFunc("/requests/", withLog(handleRequest))
	http.HandleFunc("/jobs", withLog(handleJobs))
	http.HandleFunc("/jobs/", withLog(handleJob))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Listening on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func withLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := httptest.NewRecorder()
		fn(c, r)
		log.Printf("[%d] %-4s %s\n", c.Code, r.Method, r.URL.Path)

		for k, v := range c.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(c.Code)
		c.Body.WriteTo(w)
	}
}

func handleDiscover(w http.ResponseWriter, r *http.Request) {
	body, _ := json.Marshal(map[string]interface{}{"type": "discover"})
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	handleRequests(w, r)
}

func handleRequests(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Type    string
		Payload struct {
			ID string `json:"id"`
		}
	}
	json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()

	var jobID string
	for {
		hasher := md5.New()
		hasher.Write([]byte(strconv.FormatInt(time.Now().Unix()*rand.Int63(), 10)))
		jobID = hex.EncodeToString(hasher.Sum(nil))

		if _, ok := jobs[jobID]; !ok {
			break
		}
	}
	jobs[jobID] = &Job{
		ID:      jobID,
		Type:    request.Type,
		Request: request.Payload,
	}

	var response struct {
		ID string `json:"id"`
	}
	response.ID = jobID
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	id := regexp.MustCompile("/requests/([^/]+)").FindStringSubmatch(r.URL.Path)[1]

	job := jobs[id]

	var response struct {
		Type    string      `json:"type"`
		Payload interface{} `json:"payload"`
	}
	response.Type = job.Type
	response.Payload = job.Response
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func handleJobs(w http.ResponseWriter, r *http.Request) {
	var response struct {
		Jobs []*Job `json:"jobs"`
	}
	for _, job := range jobs {
		if job.Response != nil {
			continue
		}

		response.Jobs = append(response.Jobs, job)
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func handleJob(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Payload interface{}
	}
	json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()

	id := regexp.MustCompile("/jobs/([^/]+)").FindStringSubmatch(r.URL.Path)[1]

	job := jobs[id]
	job.Response = request.Payload
}
