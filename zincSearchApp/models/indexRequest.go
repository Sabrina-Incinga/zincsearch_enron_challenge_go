package models

type IndexBulkRequest struct {
	Index   string `json:"index"`
	Records []Mail `json:"records"`
}
