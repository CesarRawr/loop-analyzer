package models

type Log struct {
  ID        string `json:"id"`
  MAC       string `json:"mac"`
  PCNAME    string `json:"pcname"`
  StartDate string `json:"startDate"`
  EndDate   string `json:"endDate"`
}

type Config struct {
  ActualID int `json:"actualID"`
}
