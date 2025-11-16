package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SebiGabor/personal-finance-cli/internal/db"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	fmt.Println("Database connected & migrations applied!")

	if len(os.Args) < 2 {
		fmt.Println("Usage: finance <command>")
		return
	}

	switch os.Args[1] {
	case "import":
		fmt.Println("Import placeholder")
	case "add":
		fmt.Println("Add placeholder")
	default:
		fmt.Println("Unknown command")
	}
}
