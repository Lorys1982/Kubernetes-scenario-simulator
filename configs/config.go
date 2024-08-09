package configs

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
)

var conf *Config

type Config struct {
	ClusterName string `yaml:"Cluster-name"`
	Scheduler   string `yaml:"Scheduler-conf"`
}

func GetClusterName() string {
	return conf.ClusterName
}

func GetScheduler() string {
	if conf.Scheduler == "" {
		return ""
	} else {
		return path.Join("configs", "topology", conf.Scheduler)
	}
}

func NewConfig() {
	yamlfile, err := os.ReadFile("configs/config.yaml")

	if err != nil {
		log.Fatal("ERROR: readfile")
	}
	//fmt.Println(string(yamlfile))
	err = yaml.Unmarshal(yamlfile, &conf)
	if err != nil {
		log.Fatal("ERROR: conversion")
	}
}
