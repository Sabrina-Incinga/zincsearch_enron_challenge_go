package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/zincsearch_enron_challenge_go/zincSearchApp/models"
)

func main() {
	rootFolder := "C:/Users/USUARIO/Desktop/Sabrina/Go/src/PruebaTecnica/enron_mail_20110402/maildir/bailey-s"
	var mailArray []models.Mail = make([]models.Mail, 0)

	err := processFilesInFolder(&mailArray, rootFolder)
	if err != nil {
		fmt.Println("Error al procesar los archivos:", err)
	}
	
	filepath := "enron_email.json"
	err = createJsonFile(mailArray, filepath)
	if err != nil {
		fmt.Println("Error al crear archivo json:", err)
	}
	err = uploadDataToIndex(filepath)
	if err != nil {
		fmt.Println("Error al subir archivos a index:", err)
	}

}

func uploadDataToIndex(filepath string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Println(err)
	}

	// Build the request body.
	req, err := http.NewRequest("POST", "http://localhost:5080/api/default/test7/_json",
		strings.NewReader(string(file)))
	if err != nil {
		return err
	}
	req.SetBasicAuth("root@example.com", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println(resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}
