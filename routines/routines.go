package routines

import (
  "os"
  "fmt"
  "time"
  "syscall"
  "analyzer/models"
  "analyzer/controller"
  "analyzer/utils"
  badger "github.com/dgraph-io/badger/v4"
)

func Clock(db *badger.DB, initialLog *models.Log, stop *chan bool) {
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

func SignalListener(signalChannel *chan os.Signal, db *badger.DB, stop *chan bool) {
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