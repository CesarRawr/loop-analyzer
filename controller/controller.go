package controller

import (
  "analyzer/models"
  "bytes"
  "encoding/json"
  "fmt"
  badger "github.com/dgraph-io/badger/v4"
  "net/http"
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

func GetAllLogs(db *badger.DB) ([]models.Log, error) {
  arr := []models.Log{}

  // Itera sobre todas las claves y valores
  err := db.View(func(txn *badger.Txn) error {
    opts := badger.DefaultIteratorOptions
    opts.PrefetchValues = true
    it := txn.NewIterator(opts)
    defer it.Close()

    for it.Rewind(); it.Valid(); it.Next() {
      item := it.Item()
      value, err := item.ValueCopy(nil)
      if err != nil {
        return err
      }

      var data models.Log
      err = json.Unmarshal(value, &data)
      if err != nil {
        return err
      }

      arr = append(arr, data)
    }
    return nil
  })

  if err != nil {
    return nil, err
  }

  return arr, nil
}

func DeleteExceptOne(db *badger.DB, protectedKey *string) error {
  err := db.Update(func(txn *badger.Txn) error {
    // Itera sobre las claves y elimina los datos excepto el que deseas conservar.
    opts := badger.DefaultIteratorOptions
    opts.PrefetchValues = false // No necesitamos los valores, solo las claves.

    it := txn.NewIterator(opts)
    defer it.Close()

    for it.Rewind(); it.Valid(); it.Next() {
      item := it.Item()
      key := item.KeyCopy(nil)

      // Verifica si es la clave que deseas conservar.
      if string(key) == *protectedKey {
        continue // Saltar la clave a conservar.
      }

      // Elimina la clave y su valor asociado.
      err := txn.Delete(key)
      if err != nil {
        return err
      }
    }

    return nil
  })

  if err != nil {
    return err
  }

  return nil
}

func SaveLogsOnDatabase(body *[]byte, url *string) error {
  fmt.Println(*url)
  resp, err := http.Post(*url, "application/json", bytes.NewBuffer(*body))
  if err != nil {
    return err
  }

  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("Bad Response")
  }

  return nil
}
