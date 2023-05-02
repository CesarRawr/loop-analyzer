package utils

import (
  "os"
  "net"
  "fmt"
  "io/ioutil"
  "encoding/json"
  "analyzer/models"
)

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
  content, err := ioutil.ReadFile("./config.json")
  if err != nil {
    return nil, err
  }

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
  err := ioutil.WriteFile("test.json", file, 0644)
  if err != nil {
    return err
  }

  return nil
}
