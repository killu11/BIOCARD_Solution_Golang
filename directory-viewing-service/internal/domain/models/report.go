package models

import "time"

type Report struct {
	Filename  string
	Timestamp time.Time
	Msg       string
}

func NewReport(filename, msg string) *Report {
	return &Report{Filename: filename, Msg: msg}
}
