package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zincsearch_enron_challenge_go/zincSearchApp/models"
)

func processFilesInFolder(mailArray *[]models.Mail, folderPath string) error {
	var wg sync.WaitGroup
	ch := make(chan models.Mail)

	err := filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if !d.IsDir() {
			wg.Add(1)
			go processFile(path, ch, &wg)
		}
		return nil
	})
	if err != nil {
		return err
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	counter := 0
	for mail := range ch {
		counter += 1
		mail.ID = counter

		*mailArray = append(*mailArray, mail)
	}

	return nil
}

func processFile(filePath string, ch chan models.Mail, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return
	}

	fileContent := string(file)

	messageID := extractValue(fileContent, "Message-ID:", "Date:")
	date := extractValue(fileContent, "Date:", "From:")
	from := extractValue(fileContent, "From:", "To:")
	to := extractValue(fileContent, "To:", "Subject:")
	subject := extractValue(fileContent, "Subject:", "Cc:")
	cc := extractValue(fileContent, "Cc:", "Mime-Version:")
	mimeVersion := extractValue(fileContent, "Mime-Version:", "Content-Type:")
	contentType := extractValue(fileContent, "Content-Type:", "Content-Transfer-Encoding:")
	contentTransferEncoding := extractValue(fileContent, "Content-Transfer-Encoding:", "Bcc:")
	bcc := extractValue(fileContent, "Bcc:", "X-From:")
	xFrom := extractValue(fileContent, "X-From:", "X-To:")
	xTo := extractValue(fileContent, "X-To:", "X-cc:")
	xcc := extractValue(fileContent, "X-cc:", "X-bcc:")
	xbcc := extractValue(fileContent, "X-bcc:", "X-Folder:")
	xFolder := extractValue(fileContent, "X-Folder:", "X-Origin:")
	xOrigin := extractValue(fileContent, "X-Origin:", "X-FileName:")
	xFileName := extractValue(fileContent, "X-FileName:", "\n")
	body := extractValue(fileContent, xFileName, "")

	mail := models.Mail{
		MessageID: messageID,
		Date:      date,
		From:      from,
		To:        to,
		Subject:   subject,
		Cc: cc,
		MimeVersion: mimeVersion,
		ContentType: contentType,
		ContentTransferEncoding: contentTransferEncoding,
		Bcc: bcc,
		XFrom: xFrom,
		XTo: xTo,
		Xcc: xcc,
		Xbcc: xbcc,
		XFolder: xFolder,
		XOrigin: xOrigin,
		XFileName: xFileName,
		Body: body,
	}

	ch <- mail
}

func extractValue(content, key string, possibleNextKey string) string {
	startIndex := strings.Index(content, key)
	if startIndex == -1 {
		return ""
	}

	startIndex += len(key)

	if possibleNextKey == "" {
		return strings.TrimSpace(content[startIndex:])
	}

	endIndex := strings.Index(content[startIndex:], possibleNextKey)
	if endIndex == -1 {
		nextKeyIndex := strings.Index(content[startIndex:], ":")
		endIndex = strings.LastIndex(content[startIndex:startIndex+nextKeyIndex], "\n") 
		if endIndex == -1 {
			return ""
		}
	}
	return strings.TrimSpace(content[startIndex:endIndex+startIndex])
}

func createJsonFile(mailsArray []models.Mail, fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// El archivo no existe, se puede crear
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		encoder := json.NewEncoder(writer)
		encoder.Encode(mailsArray)

		// Vaciar el búfer en el archivo subyacente
		err = writer.Flush()
		if err != nil {
			return err
		}
		fmt.Println("archivo enron_email.json creado correctamente")
	} 
	
	return nil
}
