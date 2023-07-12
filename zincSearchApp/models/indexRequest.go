package models

type IndexRequest struct {
	Index   string `json:"index"`
	Records []Mail `json:"records"`
}
