package main

import (
  "analyzer/controller"
  "analyzer/models"
  "analyzer/routines"
  "analyzer/utils"
  "fmt"
  badger "github.com/dgraph-io/badger/v4"
  "log"
  "os"
  "os/signal"
  "syscall"
  "time"
)

func main() {
  var isClockStopped bool
  done := make(chan bool)
  logID := utils.GenerateID(25)
  fmt.Println("pseudoID")
  fmt.Println(logID)

  // Se abre la base de datos
  db, err := badger.Open(badger.DefaultOptions(".data"))
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  // Se obtienen los datos para enviar
  logsToSend, err := controller.GetAllLogs(db)
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(logsToSend)

  // Se crea y configura el canal para señales de interrupción
  // para cuando sea detenido el proceso.
  signalChannel := make(chan os.Signal, 2)
  signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

  logFile, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
  defer logFile.Close()
  log.SetOutput(logFile)

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

  // Nuevo log creado
  initialLog := models.Log{
    ID:        logID,
    PCNAME:    *hostname,
    MAC:       mac,
    StartDate: time.Now().Format(time.RFC3339),
  }

  go routines.SignalListener(&signalChannel, db, &isClockStopped, &done)
  go routines.Clock(db, &initialLog, &isClockStopped)
  go routines.ScheduledRequest(db, &isClockStopped, &logsToSend, &logID)

  <-done
}
