package main

import (
	"bufio"
	"fmt"
	"github.com/SSSaaS/sssa-golang"
	"log"
	"os"
	"strings"
)

type resp struct {
	Name         string
	TicketNumber int
}

const (
	splitMode   = "split"
	recoverMode = "recover"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify working mode")
	}
	var workMode = os.Args[1]
	// parse flag's options
	fmt.Printf("Program work in %s mode\n", strings.ToUpper(workMode))
	switch workMode {
	case splitMode:
		splitmode()
	case recoverMode:
		recovermode()
	default:
		log.Fatal("Please specify mode: splitmode or recover")
	}

}
func splitmode() {
	var privateKey string
	var N int
	var T int
	fmt.Println("Please enter your private key")
	_, _ = fmt.Scanf("%s\n", &privateKey)
	fmt.Println("Enter count of all part and count of access part")
	_, _ = fmt.Scanf("%d %d \n", &N, &T)
	if privateKey == "" || T < 2 || N > 100 || T > N {
		log.Fatal("Incorrect input data")
	}
	keyParts, err := sssa.Create(T, N, privateKey)
	if err != nil {
		log.Fatalf("Cann't devide key to parts %s", err.Error())
	}
	for _, v := range keyParts {
		fmt.Println(v)
	}
}
func recovermode() {
	input := bufio.NewScanner(os.Stdin) //Creating a Scanner that will read the input from the console
	arrSecretParts := make([]string, 0)
	for input.Scan() {
		if input.Text() == "" {
			break
		}
		arrSecretParts = append(arrSecretParts, input.Text())

	}
	fmt.Println("Full secret key")
	fullPrivateKey, err := sssa.Combine(arrSecretParts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fullPrivateKey)

}
