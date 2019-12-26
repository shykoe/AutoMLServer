use automl;
DROP TABLE IF EXISTS `t_experiment_info`;
CREATE TABLE `t_experiment_info` (

    `experiment_id` int(11) NOT NULL AUTO_INCREMENT,
    `experiment_name` varchar(100) NOT NULL ,
    `runner` varchar(100) NOT NULL,
    `search_space` varchar(100) DEFAULT NULL ,
    `start_time` datetime DEFAULT NULL,
    `end_time` datetime DEFAULT NULL,
    `trial_concurrency` int(11) DEFAULT 1,
    `max_trial_num` int(11) DEFAULT 10,
    `algorithm_type` varchar(100) NOT NULL ,
    `algorithm_content` varchar(1000) NOT NULL ,
    `status` varchar(100) NOT NULL,
    `optimize_param` varchar(100) DEFAULT NULL,
    KEY `experiment_name` (`experiment_name`),
  PRIMARY KEY (`experiment_id`)

) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `t_trials_info`;
CREATE TABLE `t_trials_info` (

    `trial_id` int(11) NOT NULL AUTO_INCREMENT,
    `trial_name` varchar(100) NOT NULL ,
    `s3_path` varchar(100) NOT NULL,
    `parameter` varchar(1000) DEFAULT NULL ,
    `start_time` datetime DEFAULT NULL,
    `end_time` datetime DEFAULT NULL,
    `status` varchar(100) NOT NULL,
    `experiment_id` int(11)  DEFAULT NULL,
    KEY `trial_name` (`trial_name`),
    KEY `experiment_id` (`experiment_id`),
    CONSTRAINT `t_experiment_id` FOREIGN KEY (`experiment_id`) REFERENCES `t_experiment_info` (`experiment_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  PRIMARY KEY (`trial_id`)

) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `t_metric_info`;
CREATE TABLE `t_metric_info` (

    `metric_id` int(11) NOT NULL AUTO_INCREMENT,
	`dataInfo` varchar(100)  NULL,
	`loss` double DEFAULT NULL,
	`auc` double DEFAULT NULL,
	`predictAvg` double DEFAULT NULL,
	`realAvg` double DEFAULT NULL,
	`copc` double DEFAULT NULL,
	`metricType` varchar(100) NOT NULL ,
    `update_time` datetime DEFAULT NULL,
    `trial_id` int(11) NOT NULL,
    `sequence_id` int(11)  NULL,
    KEY `trial_id` (`trial_id`),
    KEY `sequence_id` (`sequence_id`),
    CONSTRAINT `t_trial_id` FOREIGN KEY (`trial_id`) REFERENCES `t_trials_info` (`trial_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  PRIMARY KEY (`metric_id`)

) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8;
