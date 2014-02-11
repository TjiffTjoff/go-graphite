package main

import (
  "net"
  "log"
  "os"
  "strings"
  "os/exec"
  "time"
  "fmt"
)

func main() {

  //Set up log file
  logFile, err := os.OpenFile("graphite.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  if err != nil {
    log.Println(err)
    os.Exit(1)
  }
  defer logFile.Close()

  log.SetOutput(logFile)

  //Read configuration from json
  config, err := getConfiguration()
  if err != nil {
    log.Printf("Configuration error: %s\n", err)
    os.Exit(1)
  }

  hostname, _ := os.Hostname()

  //Initiate loop for collecting and sending counters
  for {

    //Execute counters script that returns a csv of values
    output, err := exec.Command(config.Script).Output()
    if err != nil {
      log.Printf("Script: %s - %s\n", string(output), err)
      os.Exit(1)
    }

    //Trim spaces and split csv string into array
    valuesRaw := strings.TrimSpace(string(output))
    values := strings.Split(valuesRaw, ",")

    //Get name of the counters from configuration
    //The order of the names must correspond to the order in the csv
    names := config.Counters

    //Iterate names and counters and send to Graphite
    for i := 0; i < len(names); i++ {
      timestamp := time.Now().Unix()
      msg := fmt.Sprintf("%s.%s %s %d\n", hostname, names[i], strings.Trim(values[i],"\""), timestamp)
      if err := sendToUdp(config.Server, config.Port, msg); err != nil {
        log.Printf("UDP Error: %s\n", err)
      }
    }

  //Sleep for specified time in configuration before iteration
  time.Sleep( time.Duration(config.Interval) * time.Second)
  }
}


//Open connection to Graphite server and send message over UDP
func sendToUdp(host string, port string, msg string) error {
  addr, err := net.ResolveUDPAddr("udp", host + ":" + port)
  if err != nil {
    return err
  }
  conn, err := net.DialUDP("udp", nil, addr)
  if err != nil {
    return err
  }
  if _, err := conn.Write([]byte(msg)); err != nil {
    return err
  }

  conn.Close()

  return nil
}

