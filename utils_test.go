package main

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestCreateTar(t *testing.T) {
	createTar("/Users/shykoe/go/src/auto/mock/tmp/kwinsheng_1427131847_00000000", "/Users/shykoe/go/src/auto/mock/tmp/test.tar")

}
func TestUploadS3(t *testing.T) {
}
func TestGetBoreStatus(t *testing.T) {
	initConfig("./config.yml")
	data, _ := getBoreStatus("t_nni_kwinsheng_413687274_00000000")
	log.Info(data)
}
func TestParseMetric(t *testing.T) {
	data := `INFO[0013] DEBUG|12-26 13:10:05.649 |CommonUtils$1.run:196|2019-12-26 13:10:05: metric_count_ = 56 |ver:'2019_12_26_131005' data_info:'2019_12_19_08' total_num = 2179072 |loglosss=0.478716 |auc=0.7555 |predict_avg=0.252039 |real_avg=0.250997 |copc=0.995866 `
	metric := parseMetric(data)
	log.Info(metric)
}
func TestGetBoreLog(t *testing.T) {
	initConfig("./config.yml")
	var offset int = 0
	var everyLen = 1000
	for {
		Pos, result, str, _ := getBoreLog("t_nni_kwinsheng_413687274_00000000", "driver", "E_STDOUT", offset, everyLen)
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
			log.Info(str)
		}
		offset = Pos
	}
}
func TestExtractTar(t *testing.T) {
	initConfig("./config.yml")
	extractTar("/Users/shykoe/go/src/auto/mock/20191223/numerous_v2_beta.20191223.814959378.tar", "./mock/tmp/tmp1")
}
func TestLoginit(t *testing.T){
	initLog("testlog.log")
	log.Info("testttt")
}