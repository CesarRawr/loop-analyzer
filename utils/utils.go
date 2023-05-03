package utils

import (
  "os"
  "net"
  "fmt"
  "io/ioutil"
  "encoding/json"
  "analyzer/models"
  badger "github.com/dgraph-io/badger/v4"
)

func CloseAll(db *badger.DB, stop *chan bool) {
  *stop <- true
  db.Close()
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

// Obtener la configuración del archivo de configuración
func GetConfig() (*models.Config, error) {
  if _, err := os.Stat("config.json"); os.IsNotExist(err) {
    // Si el archivo no existe, crearlo con una configuración por defecto
    defaultConfig := models.Config{
      ActualID: 0,
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

// Guardar configuración
func SaveConfig(data *models.Config) error {
  file, _ := json.MarshalIndent(data, "", " ")
  err := ioutil.WriteFile("config.json", file, 0644)
  if err != nil {
    return err
  }

  return nil
}
