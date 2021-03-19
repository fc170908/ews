package main

import (
	"fmt"
	"log"

	"github.com/mhewedy/ews"
)

func main() {

	c := ews.NewClient(
		"https://outlook.office365.com/EWS/Exchange.asmx",
		"example@example.com",
		"examplepassword",
		&ews.Config{Dump: true, NTLM: false},
	)

	folders, err := ews.FindFolders(c, "")

	if err != nil {
		log.Fatal("err>: ", err.Error())
	}

	for _, folder := range folders {
		fmt.Printf("%s %s\n", folder.DisplayName, folder.FolderId.Id)
	}

	fmt.Println("--- success ---")
}
