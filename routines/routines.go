package routines

import (
  "analyzer/controller"
  "analyzer/models"
  "analyzer/utils"
  "encoding/json"
  "fmt"
  badger "github.com/dgraph-io/badger/v4"
  "os"
  "syscall"
  "time"
)

func Clock(db *badger.DB, initialLog *models.Log, stop *bool) {
  tdr := time.Tick(5 * time.Second)

  for {
    select {
    case actualDate := <-tdr:
      if !*stop {
        fmt.Println(actualDate)
        initialLog.EndDate = actualDate.Format(time.RFC3339)
        controller.Update(db, initialLog)
      }
    }
  }
}

func SignalListener(signalChannel *chan os.Signal, db *badger.DB, stop *bool, done *chan bool) {
  sig := <-*signalChannel
  switch sig {
  case os.Interrupt:
    utils.CloseAll(db, stop, done)
  case syscall.SIGTERM:
    utils.CloseAll(db, stop, done)
  }
}

func ScheduledRequest(db *badger.DB, stop *bool, logs *[]models.Log, actualKey *string, url *string) {
  jsonData, err := json.Marshal(logs)
  if err != nil {
    fmt.Println("Error al convertir a JSON:", err)
    return
  }

  isDataSavedOnServer := false
  tdr := time.Tick(60 * time.Second)
  for !isDataSavedOnServer {
    select {
    case <-tdr:
      *stop = true

      // Petición para guardar los datos
      err = controller.SaveLogsOnDatabase(&jsonData, url)
      if err != nil {
        fmt.Println("No se pudo enviar los logs a la base de datos")
      }

      if err == nil {
        isDataSavedOnServer = true
        err = controller.DeleteExceptOne(db, actualKey)
        if err != nil {
          fmt.Println("Algo salió mal al eliminar los datos de badger")
        }
      }

      *stop = false
    }
  }
}
