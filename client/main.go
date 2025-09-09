package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// LOGIN - REGISTER page
	fmt.Println("> Wizard Duel <")
}

func Menu(title string, args ...string) int {
	for {
		fmt.Println(title)
		for id, op := range args {
			fmt.Println(id+1, "-", op+";")
		}
		input := Input("Select an option:\n> ")
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("[error] - Invalid Option!!")
		}
		if choice > 0 && choice <= len(args) {
			return choice - 1
		}
	}

}

func Input(title string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(title)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
