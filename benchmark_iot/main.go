package main

import (
	"bufio"
	"fmt"
	"os"

	"benchmarkiot.com/streamtest"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var text string
	for text != "q" { // break the loop if text == "q"
		fmt.Print("Enter your text: ")
		scanner.Scan()
		text = scanner.Text()
		if text != "q" {
			switch text {
			case "ClientUpload":
				fmt.Println("Client Upload invoked")
				streamtest.ClientUpload()
			case "ClientDownload":
				fmt.Println("Client Download invoked")
				streamtest.ClientDownload()
			default:
				fmt.Println("function not found")
			}
		}
	}

}

func printHelp() {
	fmt.Println("The tool will help you benchmark the app.")
	fmt.Println("Enter q to exit")
	fmt.Print("List of functions you can choose from, please enter the correct name.\n\n")
	fmt.Print("1. ClientUpload\n\n")
}
