package main

import (
  "os"
  "log"
  "fmt"
  "time"
  "syscall"
  "os/signal"
  "analyzer/models"
  "analyzer/utils"
  "analyzer/routines"
  badger "github.com/dgraph-io/badger/v4"
)

func main() {
  stop := make(chan bool)
  signalChannel := make(chan os.Signal, 2)
  signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

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

  // Actualizar ID del run actual
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
  // Nuevo log creado
  initialLog := models.Log{
    ID: logID,
    MAC: mac,
    StartDate: time.Now().Format(time.RFC3339),
  }

  go routines.SignalListener(&signalChannel, db, &stop)
  go routines.Clock(db, &initialLog, &stop)
  go routines.Test(&initialLog)
  
  <-stop
}
