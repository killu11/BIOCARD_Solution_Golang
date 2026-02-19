package parsers

import (
	"bufio"
	"directory-viewing-service/internal/domain/models"
	"directory-viewing-service/pkg"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
)

const packageName = "infrastructure/workers/parsers"

var (
	ErrFileIsEmpty   = errors.New("file is empty")
	ErrStructureFile = errors.New("unexcepted file's structure")
)

type TSVParser struct {
	in  string
	out string
}

func (p *TSVParser) parseHeader(headerLine string) (map[string]int, error) {
	tabParts := strings.Split(headerLine, "\t")

	columnMap := make(map[string]int)
	currentIndex := 0

	for _, part := range tabParts {
		part = strings.TrimSpace(part)

		if part == "" {
			continue
		}

		if strings.Contains(part, " ") {
			subColumns := strings.Fields(part)
			for _, subCol := range subColumns {
				cleanCol := strings.TrimSpace(subCol)
				if cleanCol != "" {
					columnMap[cleanCol] = currentIndex
					currentIndex++
				}
			}
		} else {
			columnMap[part] = currentIndex
			currentIndex++
		}
	}
	if len(columnMap) != 15 {
		return nil, ErrStructureFile
	}
	return columnMap, nil
}

// Parse
// Парсить табличные данные документов в формате .tsv
// Дефолтные значения можно изменить при необходимости в конструкторе `models.FileData`

func (p *TSVParser) Parse(rd io.Reader, filename string) ([]*models.FileData, error) {
	brd := bufio.NewReader(rd)
	_, err := brd.ReadString('\n')
	if err == io.EOF {
		return nil, ErrFileIsEmpty
	}

	if err != nil {
		return nil, pkg.PackageError(packageName, "skip Cyrillic line", err)
	}

	cr := csv.NewReader(brd)
	records, err := cr.ReadAll()
	if err != nil {
		return nil, pkg.PackageError(packageName, "read records from file", err)
	}
	if len(records) == 1 {
		return nil, ErrFileIsEmpty
	}

	header, err := p.parseHeader(records[0][0])
	if err != nil {
		return nil, err
	}

	fds := make([]*models.FileData, 0, len(records))
	get := func(record []string, name, def string) string {
		colIdx, ok := header[name]
		if !ok || len(record) <= colIdx {
			return def
		}
		return record[colIdx]
	}

	for _, line := range records[1:] {
		line = p.normalizeRow(line[0])
		var fd = models.FileData{
			Filename: filename,
			ParseData: &models.ParseData{
				Number:    get(line, "n", "0"),
				Mqtt:      get(line, "mqtt", ""),
				InvID:     get(line, "invid", ""),
				UnitGuid:  get(line, "unit_guid", "unknown"),
				MsgID:     get(line, "msg_id", ""),
				MsgText:   get(line, "text", ""),
				Context:   get(line, "context", ""),
				Class:     get(line, "class", ""),
				Level:     get(line, "level", "0"),
				Area:      get(line, "area", ""),
				Addr:      get(line, "addr", ""),
				Block:     get(line, "block", ""),
				Type:      get(line, "type", ""),
				Bit:       get(line, "bit", ""),
				InvertBit: get(line, "invert_bit", ""),
			},
		}
		fds = append(fds, &fd)
	}
	return fds, nil
}

func (p *TSVParser) normalizeRow(s string) []string {
	tabValues := strings.Split(s, "\t")
	res := make([]string, len(tabValues))

	for i, tv := range tabValues {
		trimmed := strings.TrimSpace(tv)
		res[i] = trimmed

	}
	return res
}

// PublishPDF
// Публикует параллельно новые записи в файлы, логгирует ошибки
// numWorkers по умолчанию 5, можно поменять в зависимости от нагрузки
func (p *TSVParser) PublishPDF(data []*models.FileData) {
	UnitGuidMap := map[string]*sync.Mutex{}

	const numWorkers = 5

	for _, fd := range data {
		if _, ok := UnitGuidMap[fd.UnitGuid]; ok {
			continue
		}
		UnitGuidMap[fd.UnitGuid] = &sync.Mutex{}
	}
	errChan := make(chan error, numWorkers)
	pool := make(chan struct{}, numWorkers)

	defer func() {
		close(errChan)
		close(pool)
	}()

	go func() {
		for err := range errChan {
			slog.Warn(err.Error())
		}
	}()

	for _, fd := range data {
		pool <- struct{}{}

		go func(currentFD *models.FileData) {
			mutex := UnitGuidMap[currentFD.UnitGuid]
			mutex.Lock()
			path := p.docxPath(currentFD.UnitGuid)
			if err := pkg.AddRow(path, currentFD.ToRow()); err != nil {
				errChan <- pkg.PackageError(packageName, "add row to file", err)
			}

			defer func() {
				mutex.Unlock()
				<-pool
			}()
		}(fd)
	}
}

func (p *TSVParser) TSVPath(name string) string {
	return fmt.Sprintf("%s/%s", p.in, name)
}

func (p *TSVParser) docxPath(name string) string {
	return fmt.Sprintf("%s/%s.docx", p.out, name)
}

func NewTSVParser(in, out string) *TSVParser {
	return &TSVParser{in: in, out: out}
}
