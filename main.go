package main

import (
  "os"
  "log"
  "fmt"
  "net"
  "time"
  "analyzer/models"
  "analyzer/controller"
  badger "github.com/dgraph-io/badger/v4"
)

func clock() {
  tdr := time.Tick(1 * time.Second)
  for actualHour := range tdr {
    fmt.Println(actualHour)
  }
}

func show(db *badger.DB, id string) {
  oldLog, err := controller.Read(db, id)
  if err != nil {
    log.Fatal(err)
  }
  
  fmt.Println(*oldLog)
}

func getHostname() (*string, error) {
  hostname, err := os.Hostname()
  if err != nil {
    return nil, err
  }

  return &hostname, nil
}

func getMac() (string, error) {
  ifas, err := net.Interfaces()
  if err != nil {
    return "", err
  }

  for _, ifa := range ifas {
    a := ifa.HardwareAddr.String()
    if a != "" {
      return a, nil
    }
  }

  return "", fmt.Errorf("could not find a valid network interface with a MAC address")
}

func main() {
  hostname, err := getHostname();
  if err != nil {
    log.Fatal(err)
  }

  mac, err := getMac()
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(*hostname);
  fmt.Println(mac);

  initialLog := models.Log{}
  fmt.Println(initialLog)

  db, err := badger.Open(badger.DefaultOptions("./.data"))
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  controller.PrintAllLogs(db)
}
