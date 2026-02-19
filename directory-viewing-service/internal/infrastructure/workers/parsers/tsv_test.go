package parsers

import (
	"os"
	"testing"
)

var p *TSVParser

func TestMain(m *testing.M) {
	testDir := "/Users/egorsuslov/GolandProjects/tz/test-files/insh-data.tsv"
	p = NewTSVParser(testDir)
	m.Run()
}

func TestTSVParser_Parse(t *testing.T) {
	f, err := os.Open(p.outPath)
	if err != nil {
		t.Error(err)
		return
	}
	fds, err := p.Parse(f, f.Name())
	if err != nil {
		t.Error(err)
		return
	}
	for _, r := range fds {
		t.Logf(
			"file: %s\n"+
				"unit_guid: %s\n"+
				"msg_id: %s\n"+
				"text: %s\n"+
				"context: %s\n"+
				"class: %s\n"+
				"level: %s\n"+
				"area: %s\n"+
				"addr: %s\n"+
				"block: %s\n"+
				"type: %s\n"+
				"bit: %s\n"+
				"invert_bit: %s\n\n",
			r.Filename,
			r.UnitGuid,
			r.MsgID,
			r.MsgText,
			r.Context,
			r.Class,
			r.Level,
			r.Area,
			r.Addr,
			r.Block,
			r.Type,
			r.Bit,
			r.InvertBit,
		)
	}
}
