package main

import (
  "io/ioutil"
  "encoding/json"
)

type configuration struct {
  Server string
  Port string
  Script string
  Interval int
  Counters []string
}

func getConfiguration() (configuration, error) {
  var config configuration

  //Open and read configuration from file
  configFile := "graphite.json"
  configJson, err := ioutil.ReadFile(configFile)
  if err != nil {
    return config, err
  }

  //Parse configuration from json to configuration struct
  if err := json.Unmarshal(configJson, &config); err != nil {
    return config, err
  }

  return config, nil

}
