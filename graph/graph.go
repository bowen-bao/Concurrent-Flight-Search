package graph

import (
	"fmt"
	"math"
	"proj3/futures"
	"proj3/heapq"
	"sort"
)

// Holds the adjacent city and the weight of that edge
type AdjCity struct{
	City string
	Weight int
}

// Holds an dictionary adjacency list where the key is the
// city and the value is a list of adjCitys which contains an adjacent
// city's name and the weight of that edge
// Index is a city: int dictionary representing each city with an integer
type Graph struct{
	Cities 	map[string][]AdjCity
	Index	map[string] int
}

// Returns a pointer to a graph, represented as an empty adjacency list
func NewGraph() *Graph{
	g := make(map[string][]AdjCity)
	i := make(map[string] int)
	return &Graph{g, i}
}

// Adds edge to adjacency list
func (g *Graph) AddEdge(origin string, destination string, weight int){
	g.Cities[origin] = append(g.Cities[origin], AdjCity{destination, weight})
	//g.Cities[destination] = append(g.Cities[destination], AdjCity{})
}

// Get list of neighbors of a given node
func (g *Graph) GetAdjNodes(origin string) []AdjCity {
	return g.Cities[origin]
}

// Print the adjacency list representation of the graph
func (g *Graph) PrintMap() {
	for city := range g.Cities {
		fmt.Print(city)
		fmt.Println(g.GetAdjNodes(city))
	}
}

// Holds the adjacent city and the weight of that edge
type MinPath struct{
	Source		string
	Destination	string
	Price 		int
	Path  		[]string
}

// Perform Dijkstra's algorithm on the graph with a given source node
// Returns an array of MinPaths where each MinPath contains the info of
// the cheapest flight from the origin to that destination for all destinations
func (g *Graph) Dijkstra(origin string) [] MinPath {

	cost := make(map[string] int)
	shortestPath := make(map[string] []string)

	for city := range(g.Cities){
		cost[city] = math.MaxUint32
		shortestPath[city] = []string{origin}
	}

	cost[origin] = 0

	visited := make(map[string]bool)

	hq := heapq.NewHeapQ()

	for city := range(g.Cities) {
		hq.Push(heapq.Path{city, cost[city], shortestPath[city]})
	}

	for len(*hq.PQ) > 0 {
		cheapest := hq.Pop()
		visited[cheapest.City] = true

		for _, neighbor := range g.GetAdjNodes(cheapest.City){
			if visited[neighbor.City] == true{
				continue
			}

			if cost[cheapest.City] + neighbor.Weight < cost[neighbor.City]{
				cost[neighbor.City] = cost[cheapest.City] + neighbor.Weight
				shortestPath[neighbor.City] = append([]string{}, append(cheapest.Cities, neighbor.City)...)

				hq.Push(heapq.Path{neighbor.City, cost[neighbor.City], append([]string{}, append(cheapest.Cities, neighbor.City)...)})
			}
		}
	}

	//Create an array of shortest path flights from origin and
	//cheapest := make(map[string] MinPath)
	cheapest := make([]MinPath, 0)

	for city := range g.Cities{
		cheapest = append(cheapest, MinPath{origin, city, cost[city], shortestPath[city]})
		//cheapest[city] = MinPath{origin, city, cost[city], shortestPath[city]}
	}
	return cheapest
}

// Creates a dictionary mapping {city: index} with city in chronological order
// Used to find source and destination flights in the flight matrix using index
// instead of string names
func (g *Graph) IndexMap() {

	//Order cities in adjacency list chronologically
	keys := make([]string, 0, len(g.Cities))
	for k := range g.Cities {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	//Create dictionary that maps city: index
	cityIndex := make(map[string] int)
	for i, city := range keys{
		cityIndex[city] = i
	}
	g.Index = cityIndex
}

// Create new empty nxn matrix of MinPaths where n = number of cities
func (g *Graph) NewMatrix() [][]futures.Future {

	n := len(g.Cities)
	matrix := make([][]futures.Future, n)

	for i := 0; i < n; i++ {
		matrix[i] = make([]futures.Future, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			f := futures.NewFuture()
			matrix[i][j] = *f
		}
	}
	return matrix
}


// Fills the 2D matrix of MinPaths where the index for row and column
// are the source and destination locations represented by their indices from
// the IndexMap
func (g *Graph) FillMatrix(itineraryStream []MinPath, matrix [][]futures.Future) {

	for _, itinerary := range itineraryStream {
		source := itinerary.Source
		destination := itinerary.Destination
		sourceIndex := g.Index[source]
		destinationIndex := g.Index[destination]
		matrix[sourceIndex][destinationIndex].Set(itinerary)
	}
}

func (g *Graph) PrintMatrix(matrix [][]MinPath) {
	for i := 0; i < len(matrix); i++ {
		fmt.Println(matrix[i])
	}
}

