package main

import (
	"TravelingSalesman/graphs"
	"fmt"
)

func main() {
	fmt.Println("Traveling Salesman")
	graph := graphs.ReadFile("salesman.gl")

	graphs.PrintGraph(graph)
}
