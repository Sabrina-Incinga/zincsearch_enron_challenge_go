package zincsearchIndexer

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

// Function that walks all the files in the specified directory and processes its content
func processFilesInFolder(folderPath string) ([]models.Mail, error) {
	var mailArray []models.Mail = make([]models.Mail, 0)
	var wg sync.WaitGroup
	ch := make(chan models.Mail)

	err := filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			wg.Add(1)
			go processFile(path, ch, &wg)
		}
		return nil
	})
	if err != nil {
		return mailArray, err
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	counter := 0
	for mail := range ch {
		counter += 1
		mail.ID = counter

		mailArray = append(mailArray, mail)
	}

	return mailArray, nil
}

// Function that opens and reads the file and creates a Mail object for each one
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
		MessageID:               messageID,
		Date:                    date,
		From:                    from,
		To:                      to,
		Subject:                 subject,
		Cc:                      cc,
		MimeVersion:             mimeVersion,
		ContentType:             contentType,
		ContentTransferEncoding: contentTransferEncoding,
		Bcc:                     bcc,
		XFrom:                   xFrom,
		XTo:                     xTo,
		Xcc:                     xcc,
		Xbcc:                    xbcc,
		XFolder:                 xFolder,
		XOrigin:                 xOrigin,
		XFileName:               xFileName,
		Body:                    body,
	}

	ch <- mail
}

// Function that returns the content associated with the key
func extractValue(content, key string, possibleNextKey string) string {
	finalIndex := len(content)
	//validate if the file contains a thread of mails
	if strings.Contains(content, "-----Original Message-----") {
		finalIndex = strings.Index(content, "-----Original Message-----")
	}

	//finds the index of the key between the begining of the file and the end of it, considering the end of the file the beginning of the thread if there's one
	startIndex := strings.Index(content[:finalIndex], key)
	if startIndex == -1 {
		return ""
	}

	startIndex += len(key)

	//If there's no next key, it means it's the body of the mail, hence it takes the whole content from startIndex
	if possibleNextKey == "" {
		return strings.TrimSpace(content[startIndex:])
	}

	endIndex := strings.Index(content[startIndex:finalIndex], possibleNextKey)

	if endIndex == -1 {
		//if the file doesn't contain the specfied next key, it finds the next existent key by looking for the next ":" from startIndex
		nextKeyIndex := strings.Index(content[startIndex:], ":") + 1
		//it looks for the last line break from startIndex to nextKeyIndex
		endIndex = strings.LastIndex(content[startIndex:startIndex+nextKeyIndex], "\n")
		if endIndex == -1 {
			//in case the ":" found is in the same line of the key, there won't be any break line, so it looks for the next ":" from the previous one
			nextKeyIndex2 := strings.Index(content[startIndex+nextKeyIndex:], ":") + 1
			endIndex = strings.LastIndex(content[startIndex:startIndex+nextKeyIndex+nextKeyIndex2], "\n")

			//if there's still no endIndex, return empty string
			if endIndex == -1 {
				return ""
			}
		}
	}
	return strings.TrimSpace(content[startIndex : endIndex+startIndex])
}

// Function that calls files processing and creates the json file that holds the information of every email
func CreateJsonFile(rootFolder string, fileName string, indexName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// El archivo no existe, se puede crear
		mailsArray, err := processFilesInFolder(rootFolder)
		if err != nil {
			return fmt.Errorf("Error al procesar los archivos: %v", err)
		}
		var request models.IndexBulkRequest
		request.Index = indexName
		request.Records = mailsArray

		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		encoder := json.NewEncoder(writer)
		encoder.Encode(request)

		// Vaciar el búfer en el archivo subyacente
		err = writer.Flush()
		if err != nil {
			return err
		}
		fmt.Printf("archivo %s.json creado correctamente\n", indexName)
		return nil
	}
	fmt.Printf("archivo %s.json ya disponible\n", indexName)
	return nil
}
