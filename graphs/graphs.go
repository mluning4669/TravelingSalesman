package graphs

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//Graph an adjacency list with an associated map to map vertex names to indeces in adjlist
type Graph struct {
	dict      map[string]int
	AdjList   []List
	vertCount int
	directed  bool
	weighted  bool
}

//NewGraph a constructor for type Graph
func NewGraph(directed bool, weighted bool) *Graph {
	g := Graph{directed: directed, vertCount: 0, weighted: weighted}
	g.dict = make(map[string]int)
	g.AdjList = make([]List, 0)
	return &g
}

//Node node in homebrew doubly linked list
type Node struct {
	Val    int
	Weight *float64
	parent *Node
	prev   *Node
	next   *Node
}

//List tail, current and head of doubly linked list. An array of Lists will serve as the adjacency list. Use a map[string]int to map the vertices to the array
type List struct {
	tail    *Node
	Head    *Node
	current *Node
}

//InsertVertex inserts a vertex with no neigbors. If v is found to exist in the Graph's dictionary then it already exists so return
func (g *Graph) InsertVertex(v string) {
	_, ok := g.dict[v]
	if ok {
		return
	}
	//if v is not in the graph dictionary then add it to the adjacency list
	g.vertCount++
	e := g.vertCount - 1 //because of zero-based indexing
	g.dict[v] = e
	g.AdjList = append(g.AdjList, List{nil, nil, nil})
}

//InsertEdge inserts two vertices as an edge into graph. If the graph is directed then v1 is the head and v2 is the tail i.e v1->v2
func (g *Graph) InsertEdge(v1 string, v2 string, weight *float64) {
	e1, ok1 := g.dict[v1]
	if !ok1 {
		g.vertCount++
		e1 = g.vertCount - 1 //because of zero-based indexing
		g.dict[v1] = e1
		g.AdjList = append(g.AdjList, List{nil, nil, nil})
	}

	e2, ok2 := g.dict[v2]
	if !ok2 {
		g.vertCount++
		e2 = g.vertCount - 1 //because of zero-based indexing
		g.dict[v2] = e2
		g.AdjList = append(g.AdjList, List{nil, nil, nil})
	}

	//Insert verteces into adjacency list
	newNode1 := &Node{Val: e2, Weight: weight, prev: nil, next: nil}
	g.AdjList[e1].insertNode(newNode1)

	//check to see if g is directed or if v1 = v2 (meaning you have a edge with the same vertex).
	//If so then return
	if g.directed || v1 == v2 {
		return
	}

	//if not then insert e1 into e2's list
	newNode2 := &Node{Val: e1, Weight: weight, prev: nil, next: nil}
	g.AdjList[e2].insertNode(newNode2)
}

//InsertNode inserts new node at end of linked list
func (l *List) insertNode(newNode *Node) error {
	if l.Head == nil {
		l.Head = newNode
		l.tail = newNode
	} else {
		currentNode := l.tail
		currentNode.next = newNode
		newNode.prev = l.tail
	}
	l.tail = newNode
	return nil
}

//check for use with ReadFile
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//buildWeightedGraph builds weighted graph, either directed or undirected, and returns the graph output
func buildWeightedGraph(dir bool, lines []string) *Graph {
	edges := make([][]string, 0)
	for _, line := range lines {
		edges = append(edges, strings.Split(line, "="))
	}

	//instantiate Graph here
	graph := NewGraph(dir, true)
	//insert edges into graph
	for _, es := range edges {
		if len(strings.TrimSpace(es[1])) > 0 {
			es2, err := strconv.ParseFloat(strings.TrimSpace(es[2]), 64)
			if err != nil { //Assume 0 if no weight is provided
				es2 = 0
			}
			graph.InsertEdge(strings.TrimSpace(es[0]), strings.TrimSpace(es[1]), &es2)
		} else {
			graph.InsertVertex(strings.TrimSpace(es[0]))
		}

	}

	//return graph
	return graph
}

func buildUnweightedGraph(dir bool, lines []string) *Graph {
	edges := make([][]string, 0)
	for _, line := range lines {
		edges = append(edges, strings.Split(line, "="))
	}

	//instantiate Graph here
	graph := NewGraph(dir, false)

	//insert edges into graph
	for _, es := range edges {
		if len(strings.TrimSpace(es[1])) > 0 {
			graph.InsertEdge(strings.TrimSpace(es[0]), strings.TrimSpace(es[1]), nil)
		} else {
			graph.InsertVertex(strings.TrimSpace(es[0]))
		}

	}

	//return graph
	return graph
}

func printWeightedGraph(graph *Graph) {
	//create a map that swaps the keys and values of the graph's dictionary
	rdict := make(map[int]string)
	for k, v := range graph.dict {
		rdict[v] = k
	}

	for i, v := range graph.AdjList {
		v.current = v.Head
		//print the current vertex with the current index of the adjacency list
		fmt.Print(rdict[i], ": ")
		//check is v.current is nil to test for isolated nodes
		if v.current == nil {
			fmt.Println()
			continue
		}
		//iterate over the neighbors of v
		for {
			fmt.Print("(", rdict[v.current.Val], ", ", *v.current.Weight, ")")
			if v.current == v.tail {
				break
			}
			fmt.Print(", ")
			v.current = v.current.next
		}
		//print newline and reset current to head. It seems like the polite thing to do.
		fmt.Println()
		v.current = v.Head
	}
}

func printUnweightedGraph(graph *Graph) {
	//create a map that swaps the keys and values of the graph's dictionary
	rdict := make(map[int]string)
	for k, v := range graph.dict {
		rdict[v] = k
	}

	for i, v := range graph.AdjList {
		v.current = v.Head
		//print the current vertex with the current index of the adjacency list
		fmt.Print(rdict[i], ": ")
		//check is v.current is nil to test for isolated nodes
		if v.current == nil {
			fmt.Println()
			continue
		}
		//iterate over the neighbors of v
		for {
			fmt.Print(rdict[v.current.Val])
			if v.current == v.tail {
				break
			}
			fmt.Print(", ")
			v.current = v.current.next
		}
		//print newline and reset current to head. It seems like the polite thing to do.
		fmt.Println()
		v.current = v.Head
	}
}

//ReadFile reads in graph langauge input from a text file
func ReadFile(fileName string) *Graph {
	dat, err := ioutil.ReadFile(fileName)
	check(err)

	//Split on newline
	lines := strings.Split(string(dat), "\n")

	//Split on header
	headers := strings.Split(lines[0], " ")

	//Switch on header[1] because directed and undirected graphs are input exactly the same way
	switch strings.TrimSpace(headers[1]) {
	case "weighted":
		return buildWeightedGraph(headers[0] == "directed", lines[1:])
	case "unweighted":
		return buildUnweightedGraph(headers[0] == "directed", lines[1:])
	default:
		return nil
	}
}

//PrintGraph switches on graph input and executes appropriate output function
func PrintGraph(g *Graph) {
	if g.weighted {
		printWeightedGraph(g)
	} else {
		printUnweightedGraph(g)
	}
}
