package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

const (
	READY = "READY"
	RUNNING = "RUNNING"
	DONE = "DONE"
	ERROR = "ERROR"
	USERCANCEL = "USERCANCEL"
)
var(
	//boreUrl = "http://schedulerproxy.wsd.com/api/v1/app/task/submit"
	boreUrl = "http://localhost:12123"

	jobUrl = "http://localhost:12123/job"
)
type trial struct {
	jobId        string
	startTime    time.Time
	endTime      time.Time
	parameters   map[string]interface {}
	jobFile      string
	status       string
	s3           string
	experimentId string
}
func (t *trial) callBore(fileRoot string) error{

	filePtr, err := os.Open(path.Join(fileRoot, "bore.json"))
	defer func() {
		if filePtr == nil{
			return
		}else {
			filePtr.Close()
		}
	}()
	if err != nil {
		log.Error("Open file failed Err:", err.Error())
		return err
	}
	var boreMap map[string] interface{}
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&boreMap)
	if err != nil {
		filePtr.Close()
		return err
	}
	boreMap["appinstance_name"] = fmt.Sprintf("t_nni_%s", t.jobId)

	b, err := json.Marshal(boreMap)
	if err != nil{
		log.Error("BoreJson error: ",err)
		return err
	}
	resp, err := http.Post(boreUrl,"application/json;charset=UTF-8", bytes.NewBuffer(b) )
	if err!=nil{
		log.Error("Http error ", err)
		return err
	}
	response, _ := ioutil.ReadAll(resp.Body)
	log.Info("http Body:",  string(response))
	return nil
}
func (t *trial) getStatus(){
	//for ; ;  {
	//	req, err := http.NewRequest("GET",jobUrl,nil)
	//}
}
func (t *trial)run()  {
	paraStr, err := json.Marshal(t.parameters)
	if err!=nil{
		log.Error(err)
		t.status = ERROR
	}
	_,err = DB.Exec("INSERT INTO `automl`.`t_trials_info`(`trial_name`," +
		"`s3_path`,`parameter`,`start_time`,`end_time`,`status`,`experiment_id`) VALUES ( ?, ?, ?, ?, ?, ?, ? ) ",
		t.jobId, t.s3, paraStr, time.Now(), nil, t.status,  t.experimentId )
	if err!=nil{
		log.Error(err)
		t.close()
		return
	}
	log.Info("jobId: ",t.jobId, "\nstartTime", t.startTime.String(), "\nparameters:", t.parameters, "\nS3:", t.jobFile)
	err = t.callBore("/Users/shykoe/go/src/auto/mock/20191223/")
	if err!= nil{
		log.Error("call bore error: ",err)
		t.status = ERROR
	}
}
func (t *trial) close(){
	log.Info(" Trial close ", t.jobId)
}