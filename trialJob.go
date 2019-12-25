package main

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	READY      = "READY"
	RUNNING    = "RUNNING"
	DONE       = "DONE"
	ERROR      = "ERROR"
	USERCANCEL = "USERCANCEL"
)

type trial struct {
	jobId      string
	startTime  time.Time
	endTime    time.Time
	parameters map[string]interface{}
	status     string
	s3         string
	expId      int64
	boreFile   string
	endDir     string
}

func (t *trial) callBore(boreFile string) error {

	filePtr, err := os.Open(boreFile)
	defer func() {
		if filePtr == nil {
			return
		} else {
			filePtr.Close()
		}
	}()
	if err != nil {
		log.Error("Open file failed Err:", err.Error())
		return err
	}
	var boreMap map[string]interface{}
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&boreMap)
	if err != nil {
		filePtr.Close()
		return err
	}
	boreMap["appinstance_name"] = t.jobId
	boreMap["app_name"] = t.jobId
	b, err := json.Marshal(boreMap)
	if err != nil {
		log.Error("BoreJson error: ", err)
		return err
	}
	resp, err := http.Post(BOREURL, "application/json;charset=UTF-8", bytes.NewBuffer(b))
	if err != nil {
		log.Error("Http error ", err)
		return err
	}
	response, _ := ioutil.ReadAll(resp.Body)
	log.Info("http Body:", string(response))
	return nil
}
func (t *trial) getMetric() error{
	const url = "http://localhost:8111/getlog"
	var current int
	const everyLen = 100
	var Err int = 0
	var maxErr int = 10
	for ; ;  {
		req, err := http.NewRequest("GET",url,nil)
		if err!=nil{
			log.Error(err)
			return err
		}
		q := req.URL.Query()
		q.Add("start", strconv.Itoa(current) )
		q.Add("length", strconv.Itoa(everyLen))
		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err!=nil{
			if Err < maxErr{
				log.Error(err)
				Err += 1
				continue
			}
			return err
		}
		var data = make(map[string] string)
		response, _ := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(response,&data); err!=nil{
			log.Error(err)
			return err
		}
		logStr := []rune(data["data"])
		for ; ;  {
			aLen := runeSearch(logStr,"\n")
			if aLen == -1{
				break
			}
			current += aLen
			log.Info(string(logStr[:aLen]))
			logStr = logStr[aLen+1:]
		}



	}
}
func (t *trial) getStatus() {
	//for ; ;  {
	//	req, err := http.NewRequest("GET",jobUrl,nil)
	//}
}
func (t *trial) run() {
	paraStr, err := json.Marshal(t.parameters)
	if err != nil {
		log.Error(err)
		t.status = ERROR
	}
	_, err = DB.Exec("INSERT INTO `t_trials_info`(`trial_name`,"+
		"`s3_path`,`parameter`,`start_time`,`end_time`,`status`,`experiment_id`) VALUES ( ?, ?, ?, ?, ?, ?, ? )  ",
		t.jobId, t.s3, paraStr, time.Now(), nil, t.status, t.expId)
	if err != nil {
		log.Error(err)
		t.close()
		return
	}
	log.Info("jobId: ", t.jobId, "\nstartTime", t.startTime.String(), "\nparameters:", t.parameters, "\ns3:", t.s3)
	err = t.callBore(t.boreFile)
	if err != nil {
		log.Error("call bore error: ", err)
		t.status = ERROR
	}
}
func (t *trial) close() {
	log.Info(" Trial close ", t.jobId)
}
