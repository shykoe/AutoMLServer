package main

import (
	"container/list"
	"database/sql"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
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
	S3HOST string
	S3AK string
	S3SK string
	TMPPATH string
	BOREURL string
	BUCKET string
	SKEY string
	BORELOGURL string
	BORESTATUSURL string
	DEV string
)
type AddExpJson struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	WorkDir string `form:"workDir" json:"workDir" xml:"workDir" binding:"required"`
	TunerType string `form:"tunerType" json:"tunerType" xml:"tunerType" binding:"required"`
	TunerArgs string `form:"tunerArgs" json:"tunerArgs" xml:"tunerArgs" binding:"required"`
	SearchSpace string `form:"SearchSpace" json:"SearchSpace" xml:"SearchSpace" binding:"required"`
	Parallel int `form:"Parallel" json:"Parallel" xml:"Parallel" `
	MaxTrialNum int `form:"MaxTrialNum" json:"MaxTrialNum" xml:"MaxTrialNum" `
	BoreFile string `form:"BoreFile" json:"BoreFile" xml:"BoreFile" `
	OptimizeParam string `form:"OptimizeParam" json:"OptimizeParam" xml:"OptimizeParam"`
}
func (a *AddExpJson) toString() string{
	str := fmt.Sprintf("User: %s\nWorkDir: %s\nTunerType: %s\nTunerArgs: %s\nSearchSpace: %s\nParallel: %d", a.User, a.WorkDir, a.TunerType ,a.TunerArgs, a.SearchSpace, a.Parallel)
	return str
}

func main() {
	port := flag.Int("port",8989,"Server port")
	configFile := flag.String("config", "./config.yml", "config file")
	log.SetReportCaller(true)
	err := initConfig(*configFile)
	if err!=nil{
		log.Fatal(err)
		return
	}
	err = initDB()
	if err!=nil{
		log.Fatal(err)
		return
	}
	var exp = make(map[string] *experiment)
	r := gin.Default()
	r.POST("/addExp", func(context *gin.Context) {
		var json AddExpJson
		if err := context.ShouldBindJSON(&json); err !=nil{
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		rand.Seed(time.Now().UnixNano())
		ids := fmt.Sprintf("%s_%d",json.User ,rand.Int31())
		log.Info("get json ", json.toString())
		var parallelNum int
		if json.Parallel!=0{
			parallelNum = json.Parallel
		}else{
			parallelNum = 1
		}
		ok := IsValidOptimizeParam(json.OptimizeParam)
		if !ok{
			context.JSON(http.StatusOK, gin.H{"status": "error", "msg":"Error in OptimizeParam!"})
			return
		}
		newExp :=  &experiment{
			tunerType:json.TunerType,
			tunerArgs:json.TunerArgs,
			ll : list.New(),
			trials : list.New(),
			output: make(chan struct{}),
			trialChan: make(chan string, 100),
			searchSpace:json.SearchSpace,
			parallel:parallelNum,
			expName:ids,
			workDir:json.WorkDir,
			runner:json.User,
			maxTrialNum:json.MaxTrialNum,
			boreFile: json.BoreFile,
			optimizeParam: json.OptimizeParam,
		}
		go newExp.run()
		exp[ids] = newExp
		context.JSON(http.StatusOK, gin.H{"status": "success", "id":ids})

	})
	r.Run(fmt.Sprintf(":%d", *port))

}