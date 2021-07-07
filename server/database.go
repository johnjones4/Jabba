package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

func connect() {
	_pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
	}
	pool = _pool
}

func timeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(time.Second*10))
}

func insertJobRun(jobRun *JobRun) error {
	ctx, cancel := timeoutContext()
	defer cancel()
	err := pool.QueryRow(ctx, "INSERT INTO jobruns (created, job) VALUES ($1,$2) RETURNING \"id\"",
		jobRun.Created,
		jobRun.Job,
	).Scan(&jobRun.Id)
	if err != nil {
		return err
	}
	for i, alert := range jobRun.Alerts {
		var id int
		err := pool.QueryRow(ctx, "INSERT INTO alerts (jobrun_id, order, line, rule, description) VALUES ($1,$2,$3,$4,$5) RETURNING \"id\"",
			jobRun.Id,
			alert.Order,
			alert.Line,
			alert.Rule,
			alert.Description,
		).Scan(&id)
		if err != nil {
			return err
		}
		jobRun.Alerts[i].Id = id
	}
	return nil
}

func queryAlerts(jobRun JobRun) ([]Alert, error) {
	ctx, cancel := timeoutContext()
	defer cancel()
	rows, err := pool.Query(ctx, "SELECT id, order, line, rule, description FROM alerts ORDER BY order")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	alerts := make([]Alert, 0)
	for rows.Next() {
		alert := Alert{}
		rows.Scan(
			&alert.Id,
			&alert.Order,
			&alert.Line,
			&alert.Rule,
			&alert.Description,
		)
		alerts = append(alerts, alert)
	}
	return alerts, nil
}

func queryJobRuns(limit int, offset int) ([]JobRun, error) {
	ctx, cancel := timeoutContext()
	defer cancel()
	rows, err := pool.Query(ctx, "SELECT id, created, job FROM jobruns ORDER BY created DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	jobruns := make([]JobRun, 0)
	for rows.Next() {
		jr := JobRun{}
		rows.Scan(
			&jr.Id,
			&jr.Created,
			&jr.Job,
		)
		alerts, err := queryAlerts(jr)
		if err != nil {
			return nil, err
		}
		jr.Alerts = alerts
		jobruns = append(jobruns, jr)
	}
	return jobruns, nil
}

func queryJobRun(id int) (JobRun, error) {
	ctx, cancel := timeoutContext()
	defer cancel()
	var jobRun JobRun
	err := pool.QueryRow(ctx, "SELECT id, created, job FROM jobruns WHERE id = $1",
		jobRun.Id,
	).Scan(&jobRun.Id, &jobRun.Created, jobRun.Job)
	if err != nil {
		return jobRun, err
	}
	alerts, err := queryAlerts(jobRun)
	if err != nil {
		return jobRun, err
	}
	jobRun.Alerts = alerts
	return jobRun, nil
}
