package models

type Log struct {
  ID        string `json:"id"`
  MAC       string `json:"mac"`
  Name      string `json:"name"`
  StartDate string `json:"startDate"`
  EndDate   string `json:"endDate"`
}
