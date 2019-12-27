package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type metric struct {
	count      int
	ver        string
	dataInfo   string
	loss       float32
	auc        float32
	predictAvg float32
	realAvg    float32
	copc       float32
	metricType string
}

func (m *metric) store(trialDBId int64) error {

	_, err := DB.Exec("insert INTO t_metric_info(`loss`, `auc`, `predictAvg`, `realAvg`,"+
		" `copc`, `metricType`, `update_time`, `trial_id`, `sequence_id`) values(?,?,?,?,?,?,?,?,?) ",
		m.loss, m.auc, m.predictAvg, m.realAvg, m.copc, m.metricType, time.Now(), trialDBId, m.count)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
