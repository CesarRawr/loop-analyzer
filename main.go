package main

import (
  "log"
  "fmt"
  "time"
  "analyzer/models"
  "analyzer/controller"
  "analyzer/utils"
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

func main() {
  // Se abre la base de datos
  db, err := badger.Open(badger.DefaultOptions("./.data"))
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  // Se obtiene la configuración
  config, err := utils.GetConfig()
  if err != nil {
    log.Fatal(err)
  }

  config.ActualID++
  err = utils.SaveConfig(config);
  if err != nil {
    log.Fatal(err)
  }

  // Obtener hostname
  hostname, err := utils.GetHostname()
  if err != nil {
    log.Fatal(err)
  }

  // Obtener mac
  mac, err := utils.GetMac()
  if err != nil {
    log.Fatal(err)
  }

  // log unique id
  logID := fmt.Sprintf("%d_%s", config.ActualID, *hostname)

  fmt.Println(logID)
  fmt.Println(*hostname)
  fmt.Println(mac)

  initialLog := models.Log{
    ID: logID,
    MAC: mac,
  }

  fmt.Println(initialLog)
  controller.PrintAllLogs(db)
}
