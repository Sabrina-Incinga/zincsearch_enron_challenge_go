package zincsearchIndexer

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/pprof"

	"github.com/zincsearch_enron_challenge_go/zincSearchApp/zincsearchIndexer/variablesHandler"
)

//go:embed data/*
var data embed.FS

// Function that calls the creation of the json file to upload the data to the index, calls data upload to index and creates profiles to record performance
func RunFilesIndexer() error{
	indexConfig, err := variablesHandler.LoadEnvVariables()

	//Create cpu profile file
	cpuProfileFile, err := os.Create("zincsearchIndexer/data/cpuProfile.pprof")
	if err != nil {
		return fmt.Errorf("Error al crear archivo de perfil de CPI: %v", err)
	}
	defer cpuProfileFile.Close()

	if err := pprof.StartCPUProfile(cpuProfileFile); err != nil {
		return fmt.Errorf("Error al iniciar el perfil de CPU: %v", err)
	}
	defer pprof.StopCPUProfile()

	//Default root folder where enron email files are located
	rootFolder := "C:/Users/USUARIO/Desktop/Sabrina/Go/src/PruebaTecnica/enron_mail_20110402/maildir"

	//name of the .json file that contains all emails files information
	filepath := fmt.Sprintf("zincsearchIndexer/data/%s.json", indexConfig.IndexName)
	err = CreateJsonFile(rootFolder, filepath, indexConfig.IndexName)
	if err != nil {
		return fmt.Errorf("Error al crear archivo json: %v", err)
	}

	indexExists, err := validateIndexExistence(indexConfig)
	if !indexExists {
		fmt.Printf("El índice %s está creándose...\n", indexConfig.IndexName)
		//Upload .json file data to index if it doesn't exist
		uploadDataToIndex(filepath, indexConfig)
	} else{ fmt.Printf("El índice %s ya existe\n", indexConfig.IndexName)}


	//Create memory profile file
	memoryProfileFile, err := os.Create("zincsearchIndexer/data/memoryProfile.pprof")
	if err != nil {
		return fmt.Errorf("Error al crear archivo de perfil de memoria: %v", err)
	}
	defer memoryProfileFile.Close()
	if err := pprof.WriteHeapProfile(memoryProfileFile); err != nil {
		return fmt.Errorf("Error al escribir el perfil de memoria: %v", err)
	}

	return nil
}

// Function that sends the request to zincsearch to bulk the data to the index
func uploadDataToIndex(filepath string, indexConfig variablesHandler.IdexConfig) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	url := fmt.Sprintf("%s/_bulkv2", indexConfig.BaseUrl)

	// Build the request body.
	req, err := http.NewRequest(http.MethodPost, url, file)
	if err != nil {
		log.Println(err)
	}

	req.SetBasicAuth(indexConfig.UserName, indexConfig.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	log.Println(resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(body))
}

// Function that validates if the index exists
func validateIndexExistence(indexConfig variablesHandler.IdexConfig) (bool, error) {
	indexConfig, err := variablesHandler.LoadEnvVariables()
	var url string = fmt.Sprintf("%s/index/%s", indexConfig.BaseUrl, indexConfig.IndexName)

	validateIndexExistenceReq, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return false, err
	}

	validateIndexExistenceReq.SetBasicAuth(indexConfig.UserName, indexConfig.Password)
	validateIndexExistenceReq.Header.Set("Content-Type", "application/json")
	validateIndexExistenceReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(validateIndexExistenceReq)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	responseStatusCode := resp.StatusCode

	if responseStatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}
