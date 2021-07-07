package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func getStatus(w http.ResponseWriter, req *http.Request) {
	jsonResponse(w, Message{false, "ok"})
}

func postJobRun(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}

	jobRun := JobRun{
		Job:     vars["job"],
		Created: time.Now().UTC(),
		Log:     string(body),
	}

	if getJobDefinition(jobRun) == nil {
		errorResponse(w, fmt.Errorf("invalid job %s", jobRun.Job))
		return
	}
	err = generateAlerts(&jobRun)
	if err != nil {
		errorResponse(w, err)
		return
	}
	err = insertJobRun(&jobRun)
	if err != nil {
		errorResponse(w, err)
		return
	}
	err = saveJobRunLog(jobRun)
	if err != nil {
		errorResponse(w, err)
		return
	}
	err = transmitStatus(jobRun)
	if err != nil {
		errorResponse(w, err)
		return
	}
	jsonResponse(w, jobRun)
}

func getJobRuns(w http.ResponseWriter, req *http.Request) {
	limit, offset, err := getLimitOffset(req)
	if err != nil {
		errorResponse(w, err)
		return
	}
	jobRuns, err := queryJobRuns(limit, offset)
	if err != nil {
		errorResponse(w, err)
		return
	}
	jsonResponse(w, jobRuns)
}

func getJobRun(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		errorResponse(w, err)
		return
	}
	jobRun, err := queryJobRun(id)
	if err != nil {
		errorResponse(w, err)
		return
	}
	err = loadJobRunLog(&jobRun)
	if err != nil {
		errorResponse(w, err)
		return
	}
	jsonResponse(w, jobRun)
}

func displayJobRuns(w http.ResponseWriter, req *http.Request) {
	page := getPage(req)
	perPage := 10
	jobRuns, err := queryJobRuns(perPage, perPage*page)
	if err != nil {
		errorPage(w, err)
		return
	}
	tmpl, err := template.ParseFiles("templates/list.html")
	if err != nil {
		errorPage(w, err)
		return
	}
	tmpl.Execute(w, ListPageData{Page: page, NextPage: page + 1, PreviousPage: page - 1, JobRuns: jobRuns})
}

func displayJobRun(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		errorPage(w, err)
		return
	}
	jobRun, err := queryJobRun(id)
	if err != nil {
		errorPage(w, err)
		return
	}
	err = loadJobRunLog(&jobRun)
	if err != nil {
		errorPage(w, err)
		return
	}
	tmpl, err := template.ParseFiles("templates/detail.html")
	if err != nil {
		errorPage(w, err)
		return
	}
	tmpl.Execute(w, DetailPageData{JobRun: jobRun})
}

func runHttpServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", displayJobRuns)
	r.HandleFunc("/jobrun/{id}", displayJobRun)
	r.HandleFunc("/api", getStatus)
	r.HandleFunc("/api/jobrun/{job}", postJobRun).Methods("POST")
	r.HandleFunc("/api/jobrun", getJobRuns).Methods("GET")
	r.HandleFunc("/api/jobrun/{id}", getJobRun).Methods("GET")

	var handler http.Handler = r
	handler = logRequestHandler(handler)

	srv := &http.Server{
		Addr:    os.Getenv("HTTP_SERVER"),
		Handler: handler,
	}
	srv.ListenAndServe()
}
