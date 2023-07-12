package models

type SearchResponse struct{
	Hits	Hits	`json:"hits"`
}

type Hits struct{
	Total  	Total	`json:"total"`
	Hits	[]Data	`json:"hits"`
}

type Total struct{
	Value int	`json:"value"`
} 

type Data struct{
	Source Mail	`json:"_source"`
}