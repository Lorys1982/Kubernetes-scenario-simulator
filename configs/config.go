package configs

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
)

var conf *Config
var ni *nodeInfo

type nodeInfo struct {
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
}

type Node struct {
	ConfigName string `yaml:"name"`
	Replicas   int    `yaml:"replicas"`
}

func (node Node) GetConfName() string {
	return node.ConfigName
}

func (node Node) GetName() string {
	yamlfile, err := os.ReadFile(node.GetConfName())

	if err != nil {
		log.Fatal("ERROR: readfile")
	}
	err = yaml.Unmarshal(yamlfile, &ni)
	if err != nil {
		log.Fatal("ERROR: conversion")
	}
	return ni.Metadata.Name
}

func (node Node) GetReplicas() int {
	return node.Replicas
}

type Config struct {
	ClusterName string `yaml:"clusterName"`
	Scheduler   string `yaml:"schedulerConfig"`
	Nodes       []Node `yaml:"nodes"`
	Audit       string `yaml:"auditLoggingConfig"`
}

func GetClusterName() string {
	return conf.ClusterName
}

func GetSchedulerConf() string {
	return conf.Scheduler
}

func GetNodesConf() []Node {
	return conf.Nodes
}

func GetAuditConf() string {
	return conf.Audit
}

func fixFilePath() {
	if conf.Scheduler != "" {
		conf.Scheduler = path.Join("configs", "topology", conf.Scheduler)
	}
	if conf.Audit != "" {
		conf.Audit = path.Join("configs", "topology", conf.Audit)
	}
	for i, node := range conf.Nodes {
		conf.Nodes[i].ConfigName = path.Join("configs", "topology", node.ConfigName)
	}
}

func NewConfig() {
	yamlfile, err := os.ReadFile("configs/config.yaml")

	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlfile, &conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	fixFilePath()
}
