package utils

import (
  "analyzer/models"
  "encoding/json"
  "fmt"
  badger "github.com/dgraph-io/badger/v4"
  "io/ioutil"
  "math/rand"
  "net"
  "os"
)

func CloseAll(db *badger.DB, stop *bool, done *chan bool) {
  *stop = true
  db.Close()
  *done <- true
}

// Obtener el hostname
func GetHostname() (*string, error) {
  hostname, err := os.Hostname()
  if err != nil {
    return nil, err
  }

  return &hostname, nil
}

// Obtener mac
func GetMac() (string, error) {
  ifas, err := net.Interfaces()
  if err != nil {
    return "", err
  }

  // Se busca la interfaz con la mac
  for _, ifa := range ifas {
    a := ifa.HardwareAddr.String()
    if a != "" {
      return a, nil
    }
  }

  return "", fmt.Errorf("could not find a valid network interface with a MAC address")
}

func GenerateID(n int) string {
  var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
  b := make([]rune, n)
  for i := range b {
    b[i] = letter[rand.Intn(len(letter))]
  }

  return string(b)
}

// Obtener la configuración del archivo de configuración
func GetConfig() (*models.Config, error) {
  if _, err := os.Stat("config.json"); os.IsNotExist(err) {
    // Si el archivo no existe, crearlo con una configuración por defecto
    defaultConfig := models.Config{
      URL: "https://spcc-ccfe1.vercel.app/api/v1/logs",
    }

    configBytes, err := json.MarshalIndent(defaultConfig, "", "  ")
    if err != nil {
      return nil, err
    }

    err = ioutil.WriteFile("config.json", configBytes, 0644)
    if err != nil {
      return nil, err
    }
  }

  // leer el json
  content, err := ioutil.ReadFile("config.json")
  if err != nil {
    return nil, err
  }

  // Obtener objeto del json
  payload := models.Config{}
  err = json.Unmarshal(content, &payload)
  if err != nil {
    return nil, err
  }
  return &payload, nil
}
