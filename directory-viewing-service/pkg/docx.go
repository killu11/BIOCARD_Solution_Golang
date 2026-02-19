package pkg

import (
	"fmt"
	"path/filepath"

	docx "github.com/mmonterroca/docxgo/v2"
)

func AddRow(path, row string) error {
	name := filepath.Base(path)
	doc, err := docx.OpenDocument(path)
	if err != nil {
		doc = docx.NewDocument()
	}

	ph, err := doc.AddParagraph()
	if err != nil {
		return fmt.Errorf("add paragraph in %s", name)
	}

	r, err := ph.AddRun()
	if err != nil {
		return fmt.Errorf("run parahraph in %s", name)
	}

	if err = r.AddText(row); err != nil {
		return fmt.Errorf("add text in %s", name)
	}

	if err = doc.SaveAs(path); err != nil {
		return fmt.Errorf("save file to %s", path)
	}
	return nil
}
