package main

import (
	"fmt"

	"github.com/zincsearch_enron_challenge_go/zincSearchApp/zincsearchIndexer"
)

func main(){
	err := zincsearchIndexer.RunFilesIndexer()
	if err != nil {
		fmt.Printf("Error al crear el índice: %v\n", err)
	}

}