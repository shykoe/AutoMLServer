package main

import (
	"testing"
	"time"
)

func TestTrialJob(t *testing.T){
	job := &trial{
		jobId:"test",
		startTime:time.Now(),
		jobFile: "121",
	}
	job.run()
}