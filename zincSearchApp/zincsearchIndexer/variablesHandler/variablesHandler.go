package variablesHandler

import (
	"os"
	"github.com/joho/godotenv"
)

type IdexConfig struct {
	BaseUrl		string
	UserName    string
	Password 	string
	IndexName  	string
}

func LoadEnvVariables() (IdexConfig, error) {
	err := godotenv.Load("./config/global.env")
	if err != nil {
		return IdexConfig{}, err
	}
	return IdexConfig{
		BaseUrl: os.Getenv("baseUrl"),
		UserName: os.Getenv("userName"),
		Password: os.Getenv("password"),
		IndexName: os.Getenv("indexName"),
	}, nil
}
