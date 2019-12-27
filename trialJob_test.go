package main

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestTrialJob(t *testing.T) {
	job := &trial{
		jobId:     "test",
		startTime: time.Now(),
	}
	job.run()
}
func TestGetMetric(t *testing.T) {
	initConfig("./config.yml")
	job := &trial{
		jobId:     "kwinsheng_1511407985_00000000",
		startTime: time.Now(),
	}
	job.getMetric()
}
func TestRuneSearch(t *testing.T) {
	teststr := `asdasdasdasdasd
asdasdasd
asdasd
asda
sda
sdas
das
das
升
挖到
ww
阿萨德
das
das
d:
asd:"
asd
:asd
"`
	runes := []rune(teststr)
	ind := runeSearch(runes, "\n")
	data := string(runes[:ind])
	log.Info(data)
}
