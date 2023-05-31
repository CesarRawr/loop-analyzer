package utils

import (
  "fmt"
  badger "github.com/dgraph-io/badger/v4"
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
