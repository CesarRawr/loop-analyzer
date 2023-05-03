package main

import (
  "os"
  "log"
  "fmt"
  "time"
  "syscall"
  "os/signal"
  "analyzer/models"
  "analyzer/controller"
  "analyzer/utils"
  badger "github.com/dgraph-io/badger/v4"
)

func clock(db *badger.DB, initialLog *models.Log, stop *chan bool) {
  tdr := time.Tick(5 * time.Second)
  for {
    select {
    case <-*stop:
      return
    case actualDate := <-tdr:
      initialLog.EndDate = actualDate.Format(time.RFC3339)
      controller.Update(db, initialLog)
    }
  }
}

func signalListener(signalChannel *chan os.Signal, db *badger.DB, stop *chan bool) {
  sig := <-*signalChannel
  switch sig {
  case os.Interrupt:
    fmt.Println("SIGINT Signal")
    utils.CloseAll(db, stop)
  case syscall.SIGTERM:
    fmt.Println("SIGTERM Signal")
    utils.CloseAll(db, stop)
  }
}

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

  go signalListener(&signalChannel, db, &stop)
  go clock(db, &initialLog, &stop)
  <-stop
}
