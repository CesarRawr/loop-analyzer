package main

import (
  "log"
  "fmt"
  "time"
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

func main() {
  db, err := badger.Open(badger.DefaultOptions("./.data"))
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  controller.PrintAllLogs(db)
}
