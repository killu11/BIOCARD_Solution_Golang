package dto

import "directory-viewing-service/internal/domain/models"

type FileDataMessage struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
}

func TaskToMessage(task *models.FileTask) *FileDataMessage {
	return &FileDataMessage{
		ID:       task.ID,
		Filename: task.Filename,
	}
}
