package models

import (
	"fmt"
	"time"
)

type FileTask struct {
	ID       int
	Filename string
	Status   string
	CreateAt time.Time
}
type ParseData struct {
	Number    string `json:"number,omitempty"`
	Mqtt      string `json:"mqtt,omitempty"`
	InvID     string `json:"inv_id,omitempty"`
	UnitGuid  string `json:"unit_guid,omitempty"`
	MsgID     string `json:"msg_id,omitempty"`
	MsgText   string `json:"msg_text,omitempty"`
	Context   string `json:"context,omitempty"`
	Class     string `json:"class,omitempty"`
	Level     string `json:"level,omitempty"`
	Area      string `json:"area,omitempty"`
	Addr      string `json:"addr,omitempty"`
	Block     string `json:"block,omitempty"`
	Type      string `json:"type,omitempty"`
	Bit       string `json:"bit,omitempty"`
	InvertBit string `json:"invert_bit,omitempty"`
}
type FileData struct {
	Filename string `json:"-"`
	*ParseData
}

func (f *FileData) ToRow() string {
	return fmt.Sprintf(
		"%s %s %s %s %s %s %s %s %s %s %s %s %s %s\n",
		f.Number,
		f.Mqtt,
		f.InvID,
		f.MsgID,
		f.MsgText,
		f.Context,
		f.Class,
		f.Level,
		f.Area,
		f.Addr,
		f.Block,
		f.Type,
		f.Bit,
		f.InvertBit,
	)
}
