package main

import (
	"fmt"
	"os"
	"proj3/bookings"
	"strconv"
)

func printUsage() {
	fmt.Printf("Usage: flights <number of threads> <maps.txt> <queries.txt> \n" +
		"\t<number of threads> = the number of threads to be part of the parallel version \n" +
		"\t\t\t\t no input indicates sequential version \n" +
		"\t<maps.txt> = file with available flights and their prices \n" +
		"\t<queries.txt> = file with customer requests indicating source and destination locations \n")
}


func main() {

	mode := ""
	threads := 0
	mapFile := ""
	queryFile := ""

	if len(os.Args) == 3 {
		mode = "s"
		mapFile = os.Args[1]
		queryFile = os.Args[2]
	} else if len(os.Args) == 4 {
		mode = "p"
		threads, _ = strconv.Atoi(os.Args[1])
		mapFile = os.Args[2]
		queryFile = os.Args[3]
	} else{
		printUsage()
	}

	config := bookings.NewConfig(mode, threads, mapFile, queryFile)

	bookings.Run(*config)

}
