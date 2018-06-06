package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("===============This is tsts1.go=========-")

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//log.Fatalf("This is my log: %v", "lalala")
	log.Println("This is Sparta")
	log.Printf("This is Rome.%v", "Rome")
	log.Println("Nana")
}
