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

func getConfiguration(file string) (configuration, error) {
  var config configuration

  //Open and read configuration from file
  configFile := file
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
