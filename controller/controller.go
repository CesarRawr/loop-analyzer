package controller

import (
  "log"
  "fmt"
  "encoding/json"
  "analyzer/models"
  badger "github.com/dgraph-io/badger/v4"
)

func Create(db *badger.DB, p *models.Log) error {
  // Codificar la estructura como JSON
  jsonData, err := json.Marshal(p)
  if err != nil {
    return err
  }

  // Guardar la estructura en BadgerDB
  return db.Update(func(txn *badger.Txn) error {
    return txn.Set([]byte(p.ID), jsonData)
  })
}

func Read(db *badger.DB, id string) (*models.Log, error) {
  p := models.Log{}
  err := db.View(func(txn *badger.Txn) error {
    item, err := txn.Get([]byte(id))
    if err != nil {
      return err
    }

    return item.Value(func(val []byte) error {
      return json.Unmarshal(val, &p)
    })
  })

  if err != nil {
    return nil, err
  }
  
  return &p, nil
}

func Update(db *badger.DB, data *models.Log) error {
  tx := db.NewTransaction(true)
  defer tx.Discard()

  // Serializar la estructura log en un slice de bytes
  bytes, err := json.Marshal(data)
  if err != nil {
    return err
  }

  // Guardar el slice de bytes como valor de la clave correspondiente
  // en la base de datos Badger
  key := []byte(fmt.Sprintf(data.ID))
  err = tx.Set(key, bytes)
  if err != nil {
    return err
  }

  // Commit la transacción
  err = tx.Commit()
  if err != nil {
    return err
  }

  return nil
}

func PrintAllLogs(db *badger.DB) {
  // Itera sobre todas las claves y valores
  err := db.View(func(txn *badger.Txn) error {
    opts := badger.DefaultIteratorOptions
    opts.PrefetchValues = true
    it := txn.NewIterator(opts)
    defer it.Close()

    for it.Rewind(); it.Valid(); it.Next() {
      item := it.Item()
      key := item.Key()
      value, err := item.ValueCopy(nil)
      if err != nil {
        return err
      }
      fmt.Printf("key=%s, value=%s\n", key, value)
    }
    return nil
  })

  if err != nil {
    log.Fatal(err)
  }
}
