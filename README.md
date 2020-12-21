# Concurrent-Flight-Search
Parallelized implementation of a flight search engine where it returns the cheapest flights with the given source and destination location for every customer request

# Usage 

There are two ways to run this program:

1) Using python on terminal command line to get the speedup graph 
2) Using go on terminal command line to get the output for an individual dataset 


1) Python 

Go into the 'flights' folder, then run the command:
flights $ ./ python benchmark.py


2) Go 

Go into the 'flights' folder, then run the command:
$ ./ go run flights.go map.txt queries.txt

Usage: flights <number of threads> <maps.txt> <queries.txt> 
        <number of threads> = the number of threads to be part of the parallel version 
					    no input indicates sequential version
        <maps.txt> = file with available flights and their prices 
        <queries.txt> = file with customer requests indicating source and destination 

