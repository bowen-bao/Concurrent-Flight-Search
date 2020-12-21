package bookings

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"os"
	"proj3/balancing"
	"proj3/deque"
	"proj3/futures"
	"proj3/graph"
	"proj3/queue"
	"sync"
)

type Config struct {
	Mode    	string        	// "s" = sequential, "p" = parallel
	Threads 	int 			// Represents the number of threads to spawn
	MapsFile 	string 			// Filename of maps.txt used to create graph
	QueryFile 	string 			// Filename of queries.txt representing list of customer requests
}

func NewConfig(m string, c int, maps string, query string) *Config {
	return &Config{m, c, maps, query}
}

//Run starts up based on the configuration information provided
func Run(config Config) {
	switch config.Mode{
	case "s":
		Sequential(config)
	case "p":
		Parallel(config)
	}
}

func Sequential(config Config) {
	//Create graph of flights
	flightMap := GetMap(config)

	//Get index
	flightMap.IndexMap()

	//Convert graph into matrix of MinPaths
	matrix := flightMap.NewMatrix()

	//Perform dijkstra on all nodes and fill the matrix of cheapest routes
	for city, _ := range flightMap.Cities{
		minpaths := flightMap.Dijkstra(city)
		flightMap.FillMatrix(minpaths, matrix)
	}

	requestsQueue := requestProducer(config)

	enc := json.NewEncoder(os.Stdout)

	for _, request := range requestsQueue.Items{
		result := Process(flightMap, request.(Request), matrix)
		if err := enc.Encode(&result); err != nil {
			log.Println(err)
			continue
		}
	}
}

func ParallelAllPairs (cityq *queue.Queue, matrix [][]futures.Future, g *graph.Graph, section int){

	myCities := queue.NewQueue()

	cityq.Lock()
	//Modify number of requests taken if queue is shorter than the original number
	if cityq.Length() < section{
		section = cityq.Length()
	}

	//Append city to own queue
	for i := 0; i < section; i++ {
		myCities.Enqueue(cityq.Dequeue())
	}

	cityq.Unlock()

	for _, city := range myCities.Items{
		if city == nil{
			break
		}
		paths := g.Dijkstra(city.(string))
		g.FillMatrix(paths, matrix)
	}
}


func Parallel(config Config) {
	//Create graph of flights
	flightMap := GetMap(config)

	//Get index
	flightMap.IndexMap()

	//Add cities to shared queue
	cityQueue := queue.NewQueue()
	addCityQueue(flightMap, cityQueue)

	//Convert graph into matrix of MinPaths
	matrix := flightMap.NewMatrix()

	section := int(math.Ceil(float64(cityQueue.Length()) / float64(config.Threads)))

	for i := 0; i < config.Threads; i++ {
		go ParallelAllPairs (cityQueue, matrix, flightMap, section)
	}

	//Get stream of requests
	requestsQueue := requestProducer(config)

	requestsTaken := int(math.Ceil(float64(requestsQueue.Length()) / float64(config.Threads)))

	dequeSlice := make([]deque.Deque, 0)

	for i := 0; i < config.Threads; i++ {
		myDeque := deque.NewDeque()
		dequeSlice = append(dequeSlice, *myDeque)
	}
	balance := balancing.NewWorkSharingThread(dequeSlice, 5)

	var w sync.WaitGroup
	w.Add(config.Threads)

	for i := 0; i < config.Threads; i++ {
		go requestConsumer(i, *balance, &dequeSlice[i], requestsTaken, requestsQueue, flightMap, matrix, &w)
	}
	w.Wait()
}

func requestConsumer(id int, balance balancing.WorkSharingThread, myDeque *deque.Deque, requestsTaken int, requestsQueue *queue.Queue, flightMap *graph.Graph, matrix [][]futures.Future, w *sync.WaitGroup){

	requestsQueue.Lock()

	//Modify number of requests taken if queue is shorter than the original number
	if requestsQueue.Length() < requestsTaken{
		requestsTaken = requestsQueue.Length()
	}

	//Append request to own deque
	for i := 0; i < requestsTaken; i++ {
		qr := requestsQueue.Dequeue()
		if qr != nil{
			myDeque.Lock()
			myDeque.PushFront(qr.(Request))
			myDeque.Unlock()
		}
	}

	requestsQueue.Unlock()

	enc := json.NewEncoder(os.Stdout)

	//Process the requests in the deque
	for{
		myDeque.Lock()
		success, r := myDeque.PopFront()
		myDeque.Unlock()
		if !success {
			break
		}

		request := r.(Request)

		result := Process(flightMap, request, matrix)

		if err := enc.Encode(&result); err != nil {
			log.Println(err)
			continue
		}
		balance.Run(id)
	}
	w.Done()
}


// Represents a flight from the map.txt file
// Used to add an edge to the graph of flights
type Flight struct{
	Origin 		string
	Destination string
	Price 		int
}

// Returns an array of Flight objects based on the info in the MapsFile
func mapToGraph(config Config) []Flight {
	file, err := os.Open(config.MapsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	flightStream := make([]Flight, 0)

	dec := json.NewDecoder(file)
	for {
		var trip map[string]interface{}
		if err := dec.Decode(&trip); err == io.EOF {
			//If reaches end of file
			break
		} else if err != nil {
			log.Println(err)
			continue
		}
		// convert map to json
		jsonString, _ := json.Marshal(trip)

		// convert json to struct
		flight := Flight{}
		json.Unmarshal(jsonString, &flight)
		flightStream = append(flightStream, flight)
	}
	return flightStream
}

// Returns a graph of the flights based on info from MapsFile
func GetMap(config Config) *graph.Graph {

	//Get array of Flight objects based on the info in the MapsFile
	flightStream := mapToGraph(config)
	//Create new empty graph
	flightMap := graph.NewGraph()
	//Fill in the graph using Flight objects as edges
	for _, trip := range flightStream{
		flightMap.AddEdge(trip.Origin, trip.Destination, trip.Price)
	}
	return flightMap
}

// Adds the index of cities to a queue
func addCityQueue(g *graph.Graph, q *queue.Queue){

	for city, _ := range g.Index{
		q.Enqueue(city)
	}
}

// Returns a queue of Request objects based on the info in the QueryFile
func requestProducer(config Config) *queue.Queue {
	file, err := os.Open(config.QueryFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	requestQueue := queue.NewQueue()

	dec := json.NewDecoder(file)
	for {
		var query map[string]interface{}
		if err := dec.Decode(&query); err == io.EOF {
			//If reaches end of file
			break
		} else if err != nil {
			log.Println(err)
			continue
		}
		// convert map to json
		jsonString, _ := json.Marshal(query)

		// convert json to struct
		request := Request{}
		json.Unmarshal(jsonString, &request)
		requestQueue.Enqueue(request)
	}
	return requestQueue
}

type Request struct{
	ID 			int
	Origin 		string
	Destination string
}

type Result struct{
	ID 			int
	Origin 		string
	Destination string
	Price 		int
	Path 		[]string
}

// Return Result based on Request info by retrieving the MinPath info
// from the matrix of cheapest paths
func Process(g *graph.Graph, request Request, matrix [][]futures.Future) Result {

	source := request.Origin
	destination := request.Destination

	sourceID := g.Index[source]
	destinationID := g.Index[destination]

	itinerary := matrix[sourceID][destinationID].Get().(graph.MinPath)

	return Result{request.ID, request.Origin, request.Destination, itinerary.Price, itinerary.Path}
}

