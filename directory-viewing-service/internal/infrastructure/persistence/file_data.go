package persistence

import (
	"context"
	"directory-viewing-service/internal/domain/models"
	"directory-viewing-service/pkg"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileDataRepository struct {
	db *pgxpool.Pool
}

func NewFileDataRepository(db *pgxpool.Pool) *FileDataRepository {
	return &FileDataRepository{db: db}
}

func (f *FileDataRepository) Save(ctx context.Context, fd []*models.FileData) error {
	tx, err := f.db.Begin(ctx)
	if err != nil {
		return pkg.PackageError(packageName, "begin tx to save", err)
	}

	columns := []string{
		"filename", "number", "mqtt", "inv_id", "unit_guid",
		"msg_id", "msg_text", "context", "class", "level",
		"area", "addr", "block", "type", "bit", "invert_bit",
	}
	defer tx.Rollback(ctx)
	cfs := pgx.CopyFromSlice(len(fd), func(i int) ([]any, error) {
		return []any{
			fd[i].Filename, fd[i].Number, fd[i].Mqtt,
			fd[i].InvID, fd[i].UnitGuid, fd[i].MsgID,
			fd[i].MsgText, fd[i].Context, fd[i].Class,
			fd[i].Level, fd[i].Area, fd[i].Addr, fd[i].Block,
			fd[i].Type, fd[i].Bit, fd[i].InvertBit,
		}, nil
	})

	idn := pgx.Identifier([]string{"files_data"})
	_, err = tx.CopyFrom(
		ctx,
		idn,
		columns,
		cfs,
	)
	if err != nil {
		return pkg.PackageError(packageName, "save parse data", err)
	}
	filenames := make([]string, len(fd))
	for i, item := range fd {
		filenames[i] = item.Filename
	}

	_, err = tx.Exec(ctx, `
		UPDATE files_tasks 
		SET status = 'completed'
		WHERE filename = ANY($1)
	`, filenames)

	if err != nil {
		return pkg.PackageError(packageName, "update status by data", err)
	}
	return tx.Commit(ctx)
}

func (f *FileDataRepository) RecordsByUID(ctx context.Context, uid string, limit int) ([]*models.FileData, error) {
	var res []*models.FileData
	query := `SELECT filename, number, mqtt, inv_id, unit_guid, 
    msg_id, msg_text, context, class, level, addr, area, block, type, bit, invert_bit 
    FROM files_data as fd WHERE unit_guid = $1`

	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}

	rows, err := f.db.Query(ctx, query, uid)
	if err != nil {
		return nil, pkg.PackageError(packageName, "records by unit_guid", err)
	}
	if limit > 0 {
		res = make([]*models.FileData, 0, limit)

	}
	for rows.Next() {
		fd := models.FileData{
			Filename:  "",
			ParseData: &models.ParseData{},
		}
		err = rows.Scan(
			&fd.Filename,
			&fd.Number,
			&fd.Mqtt,
			&fd.InvID,
			&fd.UnitGuid,
			&fd.MsgID,
			&fd.MsgText,
			&fd.Context,
			&fd.Class,
			&fd.Level,
			&fd.Addr,
			&fd.Area,
			&fd.Block,
			&fd.Type,
			&fd.Bit,
			&fd.InvertBit,
		)

		if err != nil {
			return nil, pkg.PackageError(packageName, "scan records", err)
		}

		res = append(res, &fd)
	}

	if err = rows.Err(); err != nil {
		return nil, pkg.PackageError(packageName, "scan records", err)
	}
	return res, nil
}
