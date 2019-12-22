package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type trial struct {
	jobId  string
	startTime time.Time
	endTime time.Time
	parameters map[string]interface {}
	S3 string
}

func (t *trial)run()  {
	log.Info("jobId: ",t.jobId, "\nstartTime", t.startTime.String(), "\nparameters:", t.parameters, "\nS3:", t.S3)
}