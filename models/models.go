package models

type Log struct {
  ID        string `json:"id"`
  MAC       string `json:"mac"`
  StartDate string `json:"startDate"`
  EndDate   string `json:"endDate"`
}

type Config struct {
  ActualID int `json:"actualID"`
}
