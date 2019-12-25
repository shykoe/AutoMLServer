package main

import (
	"container/list"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

var DB *sql.DB
var (
	USERNAME string
	PASSWORD string
	NETWORK  = "tcp"
	SERVER   string
	PORT     string
	DATABASE string
	Count int64
)
type AddExpJson struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	WorkDir string `form:"workDir" json:"workDir" xml:"workDir" binding:"required"`
	TunerType string `form:"tunerType" json:"tunerType" xml:"tunerType" binding:"required"`
	TunerArgs string `form:"tunerArgs" json:"tunerArgs" xml:"tunerArgs" binding:"required"`
	SearchSpace string `form:"SearchSpace" json:"SearchSpace" xml:"SearchSpace" binding:"required"`
	Parallel int `form:"Parallel" json:"Parallel" xml:"Parallel" `
	MaxTrialNum int `form:"MaxTrialNum" json:"MaxTrialNum" xml:"MaxTrialNum" `
}
func (a *AddExpJson) toString() string{
	str := fmt.Sprintf("User: %s\nWorkDir: %s\nTunerType: %s\nTunerArgs: %s\nSearchSpace: %s\nParallel: %d", a.User, a.WorkDir, a.TunerType ,a.TunerArgs, a.SearchSpace, a.Parallel)
	return str
}
func initDB() error{
	var err error
	config := make(map[string] string)
	data, err := ioutil.ReadFile("./config.yml")
	err = yaml.Unmarshal(data,&config)
	if err != nil{
		log.Fatal(err)
		return err
	}
	USERNAME = config["username"]
	PASSWORD = config["password"]
	SERVER = config["server"]
	DATABASE = config["database"]
	PORT = config["port"]
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s",USERNAME,PASSWORD,NETWORK,SERVER,PORT,DATABASE)
	DB, err = sql.Open("mysql",dsn)
	if err!=nil{
		log.Fatal(err)
		return err
	}
	DB.SetConnMaxLifetime(100*time.Second)  //最大连接周期，超过时间的连接就close
	DB.SetMaxOpenConns(1000)//设置最大连接数
	DB.SetMaxIdleConns(0) //设置闲置连接数
	return nil
}
func main() {
	log.SetReportCaller(true)
	err := initDB()
	if err!=nil{
		log.Fatal(err)
		return
	}
	var exp = make(map[string] *experment )
	r := gin.Default()
	r.POST("/addExp", func(context *gin.Context) {
		var json AddExpJson
		//data,err :=context.GetRawData(); if err== nil{
		//	panic(err)
		//}
		//log.Info(string(data))
		if err := context.ShouldBindJSON(&json); err !=nil{
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ids := fmt.Sprintf("%s_%d",json.User ,rand.Int31())
		log.Info("get json ", json.toString())
		var parallelNum int
		if json.Parallel!=0{
			parallelNum = json.Parallel
		}else{
			parallelNum = 1
		}
		newExp :=  &experment{
			tunerType:json.TunerType,
			tunerArgs:json.TunerArgs,
			ll : list.New(),
			trials : list.New(),
			output: make(chan struct{}),
			searchSpace:json.SearchSpace,
			parallel:parallelNum,
			expName:ids,
			workDir:json.WorkDir,
			runner:json.User,
			maxTrialNum:json.MaxTrialNum,
		}
		go newExp.run()
		exp[ids] = newExp
		context.JSON(http.StatusOK, gin.H{"status": "success", "id":ids})

	})
	r.Run(":8989")

}