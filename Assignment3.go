package main

import (
	"TravelingSalesman/graphs"
	"fmt"
	"math"
)

var globalMin float64
var zeroEdges map[int]float64

func main() {
	fmt.Println("Traveling Salesman")
	graph := graphs.ReadFile("salesman.gl")

	//graphs.PrintGraph(graph)

	currentPath := graphs.Path{PathList: &graphs.List{}, PathCount: 1}
	visited := make(map[int]bool)

	for i := 0; i < graph.VertCount; i++ { //front load visited map with false
		visited[i] = false
	}

	zeroEdges := make(map[int]float64) //all the edges to the first element for O(1) lookup time

	var edges = graph.AdjList[0].Head

	for edges.Next != nil {
		zeroEdges[edges.Val] = *edges.Weight
		edges = edges.Next
	}

	tsp(graph, 0, currentPath, 0, visited)
}

func tsp(graph *graphs.Graph, currentVert int, currentPath graphs.Path, pathCost float64, visited map[int]bool) {
	if currentPath.PathCount == graph.VertCount && currentVert > 0 {
		currentPath.AppendToPath(&graphs.Node{Val: 0})
		pathCost = pathCost + zeroEdges[currentVert]
		globalMin = math.Min(globalMin, pathCost)
		fmt.Print(pathCost, " - ")
		currentPath.PrintPath()

	} else {
		visited[currentVert] = true
		var edge = graph.AdjList[currentVert].Head

		for edge.Next != nil {
			if !visited[edge.Val] {
				visited[edge.Val] = true
				currentPath.AppendToPath(&graphs.Node{Val: edge.Val})
				tsp(graph, edge.Val, currentPath, pathCost+*edge.Weight, visited)
				currentPath.RemoveFromPath()
			}
			edge = edge.Next
		}
	}
}
