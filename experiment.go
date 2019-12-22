package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"sync"
	"time"
)
type experment struct {
	write4out *os.File
	read4nni *os.File
	write4nni *os.File
	read4out *os.File
	tunerType string
	tunerArgs string
	ll *list.List
	trials *list.List
	output chan struct{}
	mu    sync.RWMutex
	searchSpace string
	parallel int
	expId string
	S3 string
}
func (e *experment) listen(){
	for ; ;  {
		out := make([]byte, 10000)
		_, err := e.read4out.Read(out)

		if err != nil{
			log.Warn("read error", err)
			break
		}
		err, cmdData := convert(out)
		log.Info("listent: ",cmdData.cmdContent)
		e.ll.PushBack(cmdData)
		if err != nil{
			log.Warn("convert error", err)
			break
		}
		e.output <- struct{}{}
	}
}
func (e *experment) work()  {
	for ; ;  {
		select {
		case <- e.wait():
			for{
				data := e.get()
				if data!= nil{
					log.Info("input data !! ",data.cedType, data.cmdContent)
					switch data.cedType {
					case NewTrialJob:
						jsonData := make(map[string] interface{})
						json.Unmarshal([]byte(data.cmdContent),&jsonData)
						params,ok := jsonData["parameters"].(map[string]interface {})
						if !ok{
							log.Error("params Decode Error")
							return
						}
						ind,ok:= jsonData["parameter_id"].(int)
						if !ok{
							log.Error("Parameter_id Decode Error")
							return
						}
						job := &trial{
							jobId:fmt.Sprintf("%s_%d", e.expId, ind),
							startTime: time.Now(),
							parameters:params,
							S3:e.S3,
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
func(e *experment) wait() <- chan struct{}{
	return e.output
}
func(e *experment) get() *IpcData {
	e.mu.Lock()
	defer e.mu.Unlock()
	if elem := e.ll.Front(); elem!=nil{
		data := elem.Value.(*IpcData)
		e.ll.Remove(elem)
		return data
	}
	return nil
}
func(e *experment) keepAlive()  {
	for ; ;  {
		alive := IpcData{
			cedType:Ping,
			cmdContent:"",
		}
		e.send(alive)
		time.Sleep(1 * time.Second)
	}
}
func (e *experment) run(){
	cmd := exec.Command("python", "-m", "nni", "--tuner_class_name", e.tunerType, "--tuner_args", e.tunerArgs)
	cmd.Env = os.Environ()
	r1, w1, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	r2, w2, err := os.Pipe()
	if err != nil{
		panic(err)
	}
	pip :=make([]*os.File,2)
	pip[0] = r1
	pip[1] = w2
	cmd.ExtraFiles = pip
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	e.read4nni = r1
	e.write4out = w1
	e.write4nni = w2
	e.read4out = r2
	go cmd.Run()
	go e.listen()
	go e.keepAlive()
	go e.work()
	initExp := IpcData{
		cedType:Initialize,
		cmdContent:e.searchSpace,
	}
	e.send(initExp)
	for ; ;  {
		if e.trials.Len() < e.parallel{
			//should add jobs
			need := e.parallel - e.trials.Len()
			for i:=0; i<need ;i++  {
				initExp := IpcData{
					cedType:RequestTrialJobs,
					cmdContent:"1",
				}
				e.send(initExp)
			}

		}
		time.Sleep(2*time.Second)
	}
}
func (e *experment) send(data IpcData){
	log.Info("send", data.cedType, data.cmdContent)
	byteContent := data.decode()
	log.Info("send ",byteContent)
	_, err := e.write4out.Write(byteContent)
	if err != nil {
		panic(err)
	}
}