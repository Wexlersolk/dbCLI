package main

import (
	"dbCLI/cmd"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:MyNewPass@tcp(127.0.0.1:3306)/Lab1?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	cmd.SetDB(db)

	for {
		fmt.Println("1. Show Tables")
		fmt.Println("2. Show Structure")
		fmt.Println("3. Call HeadCount")
		fmt.Println("4. Calculate Total Price Of All Books")
		fmt.Println("5. Get Books After Year")
		fmt.Println("5. Something about price")
		fmt.Println("6. Show books based on price and number of pages")
		fmt.Println("7. Sort and show books based on price")
		fmt.Println("8. Change publisher")
		fmt.Println("9. Show table values")
		fmt.Println("10. Exit")

		var choice string
		fmt.Print("Enter your choice (1-5): ")
		fmt.Scanln(&choice)

		switch choice {
		case "1":

			fmt.Println("trying to run a command")
			err := runCommand("dbCLI show-tables")
			if err != nil {
				fmt.Println("Error executing command:", err)
			}

		case "2":

			fmt.Println("enter the name of the table, bitch")
			fmt.Scanln(&choice)
			err := runCommand("dbCLI show-structure " + choice)
			if err != nil {
				fmt.Println("Error executing command:", err)
			}

		case "3":
			fmt.Println("enter your number of letters")
			fmt.Scanln(&choice)
			err := runCommand("dbCLI call-HeadCount " + choice)
			if err != nil {
				fmt.Println("Error executing command:", err)
			}

		case "4":

			err := runCommand("dbCLI calculate-total-price")
			if err != nil {
				fmt.Println("Error executing command:", err)
			}

		case "5":
			fmt.Println("enter your year")
			fmt.Scanln(&choice)
			err := runCommand("dbCLI get-books-after-year " + choice)
			if err != nil {
				fmt.Println("Error executing command:", err)
			}

		case "6":

			fmt.Println("enter max price")
			fmt.Scanln(&choice)
			var maxpages string
			fmt.Println("enter max mages")
			fmt.Scanln(&maxpages)
			err := runCommand("dbCLI show-books-Mpages-Mprice " + maxpages + " " + choice)
			if err != nil {
				fmt.Println("Error executing command:", err)
			}

		case "7":

			err := runCommand("dbCLI sort-books")
			if err != nil {
				fmt.Println("Error executing command:", err)
			}

		case "8":

			fmt.Println("enter book library code")
			fmt.Scanln(&choice)
			var newpublisher string
			fmt.Println("enter new publisher")
			fmt.Scanln(&newpublisher)
			err := runCommand("dbCLI change-publisher " + choice + " " + newpublisher)
			if err != nil {
				fmt.Println("Error executing command:", err)
			}

		case "9":

			fmt.Println("enter the name of the table")
			fmt.Scanln(&choice)
			err := runCommand("dbCLI show-values " + choice)
			if err != nil {
				fmt.Println("Error executing command:", err)
			}

		case "10":
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please enter a number from 1 to 5.")
		}

		fmt.Print("Do you want to continue? (yes/no): ")
		var continueChoice string
		fmt.Scanln(&continueChoice)

		if strings.ToLower(continueChoice) != "yes" && strings.ToLower(continueChoice) != "y" {
			os.Exit(0)
		}
	}
}

func runCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)

	// Capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	// Print the output
	fmt.Println(string(output))

	return nil
}
