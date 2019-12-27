package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	READY      = "READY"
	RUNNING    = "RUNNING"
	SUCCESS    = "success"
	ERROR      = "ERROR"
	USERCANCEL = "USERCANCEL"
)

type trial struct {
	ind        int
	jobId      string
	startTime  time.Time
	endTime    time.Time
	parameters map[string]interface{}
	status     string
	s3         string
	expId      int64
	boreFile   string
	endDir     string
	metricll   *list.List
	exp        *experiment
	dbId       int64
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
	s3Url := fmt.Sprintf("%s:%s", "http://s3sz.sumeru.mig/algbaseserviceapi", t.s3)
	boreMap["program_urls"] = []string{s3Url}

	log.Info(boreMap)
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
func (t *trial) getMetric() error {
	var offset int = 0
	var everyLen = 1000
	for {
		if t.status == SUCCESS || t.status == ERROR || t.status == USERCANCEL {
			return nil
		}
		Pos, result, str, _ := getBoreLog(t.jobId, "driver", "E_STDOUT", offset, everyLen)
		log.Info("offset: ", offset)
		if str == "" {
			time.Sleep(time.Second * 5)
			continue
		}
		if offset == Pos {
			everyLen *= 2
			continue
		}
		for _, str := range result {
			currentMetric := parseMetric(str)
			if currentMetric != nil {
				if currentMetric.dataInfo == t.endDir {
					currentMetric.metricType = "FINAL"
				} else {
					currentMetric.metricType = "PERIODICAL"
				}
				err := currentMetric.store(t.dbId)
				if err != nil {
					panic(err)
				}
				t.metricll.PushBack(currentMetric)
				time.Sleep(2 * time.Second)
				err = t.exp.updateMetric(currentMetric, t.ind)
				if err != nil {
					return err
				}
			}
		}
		offset = Pos
		time.Sleep(1 * time.Second)
	}
	return nil
}
func (t *trial) updateStatus(status string) error {
	_, err := DB.Exec("UPDATE `t_trials_info` SET `status` = ? WHERE  `trial_id` = ?", status, t.dbId)
	if err != nil {
		return err
	}
	return nil
}
func (t *trial) getStatus() {
	for {
		statusData, err := getBoreStatus(t.jobId)
		status := parseStatus(statusData["appinstance_status"])
		err = t.updateStatus(status)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		t.status = status
		if status == ERROR || status == USERCANCEL {
			t.exp.trialChan <- t.jobId
			return
		} else if status == SUCCESS {
			//wait for update metric
			time.Sleep(10 * time.Second)
			t.exp.trialChan <- t.jobId
			return
		}
		time.Sleep(2 * time.Second)
	}
}
func (t *trial) run() {
	paraStr, err := json.Marshal(t.parameters)
	if err != nil {
		log.Error(err)
		t.status = ERROR
	}
	result, err := DB.Exec("INSERT INTO `t_trials_info`(`trial_name`,"+
		"`s3_path`,`parameter`,`start_time`,`end_time`,`status`,`experiment_id`) VALUES ( ?, ?, ?, ?, ?, ?, ? )  ",
		t.jobId, t.s3, paraStr, time.Now(), nil, t.status, t.expId)
	if err != nil {
		log.Error(err)
		t.close()
		return
	}
	t.dbId, err = result.LastInsertId()
	if err != nil {
		log.Error(err)
		t.close()
		return
	}
	log.Info("jobId: ", t.jobId, "\nstartTime", t.startTime.String(), "\nparameters:", t.parameters, "\ns3:", t.s3)
	err = t.callBore(t.boreFile)
	go t.getStatus()
	go t.getMetric()
	if err != nil {
		log.Error("call bore error: ", err)
		t.status = ERROR
	}
}
func (t *trial) close() {
	log.Info(" Trial close ", t.jobId)
}
