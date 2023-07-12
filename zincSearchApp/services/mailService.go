package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/zincsearch_enron_challenge_go/zincSearchApp/models"
	"github.com/zincsearch_enron_challenge_go/zincSearchApp/zincsearchIndexer/variablesHandler"
)

func GetMailsByQuery(request models.SearchRequest) (models.SearchResponse, error) {
	var response models.SearchResponse
	
	indexConfig, err := variablesHandler.LoadEnvVariables()
	if err != nil {
		return response, err
	}

	url := indexConfig.BaseUrl+"/enron_mail/_search"

	query, err := json.Marshal(request)
	if err !=nil{
		return response, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(query)))

	if err != nil {
		fmt.Println(err)
		return response, err
	}

	req.SetBasicAuth(indexConfig.UserName, indexConfig.Password)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return response, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
		return response, err
	}

	return response, nil
}