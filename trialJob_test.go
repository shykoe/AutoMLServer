package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httptest"
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
func TestBoreKill(t *testing.T){
	l, _ := net.Listen("tcp", "127.0.0.1:12123")
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		// Send response to be tested
		result := make(map[string]interface{})
		result["ret_code"]  = 0
		result["err_msg"] = ""
		b, _ := json.Marshal(result)
		rw.Write(b)
	}))
	server.Listener.Close()
	server.Listener = l
	server.Start()
	defer server.Close()
	initConfig("./config.yml")
	newExp := &experiment{
		runner:        "testA",
	}
	job := &trial{
		jobId:     "kwinsheng_1511407985_00000000",
		startTime: time.Now(),
		exp:newExp,
	}
	job.killBore()
	server.Close()
}
