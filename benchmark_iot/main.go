package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"benchmarkiot.com/streamtest"
)

func main() {
	printHelp()
	scanner := bufio.NewScanner(os.Stdin)

	var text string
	for text != "q" { // break the loop if text == "q"
		fmt.Print("Enter fun number: ")
		scanner.Scan()
		text = scanner.Text()
		if text != "q" {
			switch text {
			case "1":
				fmt.Println("Client Upload invoked")
				clients, size, count, filename := scanClientUpload()
				fmt.Println("clients: ", clients)
				fmt.Println("messagesize: ", size)
				fmt.Println("messagecount: ", count)
				streamtest.ClientUpload(clients, size, count, filename)
			case "2":
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
	fmt.Print("List of functions you can choose from, please enter the correct number.\n\n")
	fmt.Print("1. ClientUpload\n\n")
	fmt.Print("2. ClientDownload\n\n")
}

func scanClientUpload() (int, int, int, string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter number of clients: ")
	scanner.Scan()
	text := scanner.Text()
	clients, _ := strconv.Atoi(text)

	fmt.Print("Enter message size in bytes: ")
	scanner.Scan()
	text = scanner.Text()
	messagesize, _ := strconv.Atoi(text)

	fmt.Print("Enter message count per client: ")
	scanner.Scan()
	text = scanner.Text()
	messagecount, _ := strconv.Atoi(text)

	fmt.Print("Enter filename to write results to: ")
	scanner.Scan()
	filename := scanner.Text()

	return clients, messagesize, messagecount, filename
}
