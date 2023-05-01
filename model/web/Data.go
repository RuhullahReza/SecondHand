package web

import "github.com/google/uuid"

type CreateDataRequest struct {
	Name string `db:"name" json:"name" binding:"required" conform:"name"`
}

type DataResponse struct {
	Id   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
}

type GetByIdRequest struct {
	ID string `uri:"id"`
}

type GetByPath struct {
	Path string `uri:"path"`
}