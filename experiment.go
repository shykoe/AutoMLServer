package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type experiment struct {
	write4out   *os.File
	read4nni    *os.File
	write4nni   *os.File
	read4out    *os.File
	tunerType   string
	tunerArgs   string
	ll          *list.List
	trials      *list.List
	output      chan struct{}
	mu          sync.RWMutex
	searchSpace string
	parallel    int
	maxTrialNum int
	expName     string
	workDir     string
	tuner       *exec.Cmd
	jobFile     string
	runner      string
	expId       int64
	currentNum  int
	boreFile    string
}

func (e *experiment) listen() {
	for {
		out := make([]byte, 10000)
		_, err := e.read4out.Read(out)

		if err != nil {
			log.Warn("read error", err)
			break
		}
		err, cmdData := convert(out)
		log.Info("listent: ", cmdData.cmdContent)
		e.ll.PushBack(cmdData)
		if err != nil {
			log.Warn("convert error", err)
			break
		}
		e.output <- struct{}{}
	}
}
func (e *experiment) prepareTrial(trialId string, workSpace string, params *map[string]interface{}) (*trial, error) {
	trialPath := path.Join(TMPPATH, trialId)
	var trialTar string = fmt.Sprintf("%s/%s.tar", TMPPATH, trialId)
	err := copy.Copy(e.workDir, trialPath)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("write params file")
	f, err := os.Create(path.Join(trialPath, "automl.py"))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer f.Close()
	//miniBatch param is in bore.json, save it
	var miniBatch int = 0
	for k, v := range *params {
		var innerData string
		if k == "minibatch" {
			miniBatch = v.(int)
		}
		if reflect.TypeOf(v).String() == "string" {
			innerData = fmt.Sprintf("%s='%s'", k, v)
		} else {
			innerData = fmt.Sprintf("%s=%g", k, v)
		}
		f.WriteString(innerData + "\n")
	}
	f.Close()

	//modify bore.json
	numerousJson := make(map[string]string)

	jsonFile, err := os.Open(path.Join(trialPath, "numerous.json"))
	if err != nil {
		return nil, err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal([]byte(byteValue), &numerousJson)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	var endDir string
	for k, _ := range numerousJson {
		if k == "minibatch" && miniBatch != 0 {
			numerousJson[k] = strconv.Itoa(miniBatch)
		}
		if k == "model_name" {
			numerousJson[k] = trialId
		}
		if k == "fs.train_end_dir" {
			endDir = numerousJson[k]
		}
	}

	jsonStr, err := json.Marshal(numerousJson)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = ioutil.WriteFile(path.Join(trialPath, "numerous.json"), jsonStr, os.ModePerm)
	if err != nil {
		return nil, err
	}
	err = createTar(trialPath, trialTar)
	if err != nil {
		return nil, err
	}
	s3Path := fmt.Sprintf("%s/%s", "testAutoml", trialTar)
	err = uploadS3(trialTar, BUCKET, s3Path, "public-read")
	if err!=nil{
		return nil, err
	}
	job := &trial{
		jobId:      trialId,
		startTime:  time.Now(),
		parameters: *params,
		status:     READY,
		expId:      e.expId,
		boreFile:   e.boreFile,
		endDir:     endDir,
		s3: s3Path,
	}
	err = os.RemoveAll(trialPath)
	if err !=nil{
		return nil, err
	}
	err = os.RemoveAll(trialTar)
	if err !=nil{
		return nil, err
	}
	return job, nil
}
func (e *experiment) work() {
	for {
		select {
		case <-e.wait():
			for {
				data := e.get()
				if data != nil {
					log.Info("input data !! ", data.cedType, data.cmdContent)
					switch data.cedType {
					case NewTrialJob:
						jsonData := make(map[string]interface{})
						json.Unmarshal([]byte(data.cmdContent), &jsonData)
						params, ok := jsonData["parameters"].(map[string]interface{})
						if !ok {
							log.Error("params Decode Error")
							e.close()
							return
						}
						ind, ok := jsonData["parameter_id"].(float64)
						if !ok {
							log.Error("Parameter_id Decode Error")
							log.Info()
							e.close()
							return
						}
						jobId := fmt.Sprintf("%s_%08.0f", e.expName, ind)
						job, err := e.prepareTrial(jobId, e.workDir, &params)
						if err != nil {
							log.Error(err)
							e.close()
							return
						}
						go job.run()
						e.trials.PushBack(job)
					case Initialized:
						break
					default:
						break

					}
				}
				break
			}

		}
	}
}
func (e *experiment) wait() <-chan struct{} {
	return e.output
}
func (e *experiment) get() *IpcData {
	e.mu.Lock()
	defer e.mu.Unlock()
	if elem := e.ll.Front(); elem != nil {
		data := elem.Value.(*IpcData)
		e.ll.Remove(elem)
		return data
	}
	return nil
}
func (e *experiment) keepAlive() {
	for {
		alive := IpcData{
			cedType:    Ping,
			cmdContent: "",
		}
		if err := e.send(alive); err != nil {
			log.Warn("keepAlive send error! ")
			break
		}
		time.Sleep(1 * time.Second)
	}
}
func (e *experiment) getS3File() string {

	return "/Users/shykoe/go/src/auto/mock/20191223/s3"
}
func (e *experiment) run() {
	cmd := exec.Command("python", "-m", "nni", "--tuner_class_name", e.tunerType, "--tuner_args", e.tunerArgs)
	cmd.Env = os.Environ()
	r1, w1, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	r2, w2, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	pip := make([]*os.File, 2)
	pip[0] = r1
	pip[1] = w2
	cmd.ExtraFiles = pip
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	e.read4nni = r1
	e.write4out = w1
	e.write4nni = w2
	e.read4out = r2
	e.tuner = cmd
	go cmd.Run()
	go e.listen()
	e.jobFile = e.getS3File()
	//go e.keepAlive()
	go e.work()
	initExp := IpcData{
		cedType:    Initialize,
		cmdContent: e.searchSpace,
	}
	e.send(initExp)
	result, err := DB.Exec("insert INTO t_experiment_info(`experiment_name`, `runner`, `search_space`, `start_time`,"+
		" `trial_concurrency`, `max_trial_num`, `algorithm_type`, `algorithm_content`, `status`) values(?,?,?,?,?,?,?,?,?) ",
		e.expName, e.runner, e.searchSpace, time.Now(), e.parallel, e.maxTrialNum, e.tunerType, e.tunerArgs, READY)
	if err != nil {
		log.Error(err)
		e.close()
		return
	}
	e.expId, err = result.LastInsertId()
	if err != nil {
		log.Error(err)
		e.close()
		return
	}
	e.currentNum = 0
	var total int
	for {
		if e.currentNum < e.parallel && total < e.maxTrialNum {
			//should add jobs
			need := e.parallel - e.currentNum
			for i := 0; i < need; i++ {
				initExp := IpcData{
					cedType:    RequestTrialJobs,
					cmdContent: "1",
				}
				if err := e.send(initExp); err != nil {
					log.Warn("send error! Need to close!")
					e.close()
					return
				}
				e.currentNum += 1
				total += 1
			}

		}
		time.Sleep(2 * time.Second)
	}
}
func (e *experiment) send(data IpcData) error {
	log.Info("send", data.cedType, data.cmdContent)
	byteContent := data.decode()
	_, err := e.write4out.Write(byteContent)
	return err
}
func (e *experiment) close() {
	if e.tuner != nil {
		e.tuner.Process.Kill()
	}
	if e.read4out != nil {
		e.read4out.Close()
	}
	if e.read4nni != nil {
		e.read4nni.Close()
	}
	if e.write4out != nil {
		e.write4out.Close()
	}
	if e.write4nni != nil {
		e.write4nni.Close()
	}
}
