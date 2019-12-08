package graphs

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//Graph an adjacency list with an associated map to map vertex names to indeces in adjlist
type Graph struct {
	Dict      map[string]int
	Idict     map[int]string //inverse of dict
	AdjList   []List
	VertCount int
	directed  bool
	weighted  bool
}

//NewGraph a constructor for type Graph
func NewGraph(directed bool, weighted bool) *Graph {
	g := Graph{directed: directed, VertCount: 0, weighted: weighted}
	g.Dict = make(map[string]int)
	g.Idict = make(map[int]string)
	g.AdjList = make([]List, 0)
	return &g
}

//Path a representation of a path through a graph with the assumption that no node is visited more than once
type Path struct {
	PathList  *List
	PathCount int
}

//Node node in homebrew doubly linked list
type Node struct {
	Val       int
	HeapLabel string
	Weight    *float64
	AttCost   float64
	Label     string
	Prev      *Node
	Next      *Node
	Parent    *Node
}

//List tail, current and head of doubly linked list. An array of Lists will serve as the adjacency list. Use a map[string]int to map the vertices to the array
type List struct {
	Tail    *Node
	Head    *Node
	Current *Node
}

//AppendToPath add node to end of PathList and bump PathCount by 1
func (p *Path) AppendToPath(node *Node) {
	p.PathList.insertNode(node)
	p.PathCount++
}

//RemoveFromPath remove node from end of PathList and shrink PathCount by 1
func (p *Path) RemoveFromPath() {
	p.PathList.Tail = p.PathList.Tail.Prev
	p.PathList.Tail.Next = nil
	p.PathCount--
}

//InsertVertex inserts a vertex with no neigbors. If v is found to exist in the Graph's dictionary then it already exists so return
func (g *Graph) InsertVertex(v string) {
	_, ok := g.Dict[v]
	if ok {
		return
	}
	//if v is not in the graph dictionary then add it to the adjacency list
	g.VertCount++
	e := g.VertCount - 1 //because of zero-based indexing
	g.Dict[v] = e
	g.Idict[e] = v
	g.AdjList = append(g.AdjList, List{nil, nil, nil})
}

//InsertEdge inserts two vertices as an edge into graph. If the graph is directed then v1 is the head and v2 is the tail i.e v1->v2
func (g *Graph) InsertEdge(v1 string, v2 string, weight *float64) {
	e1, ok1 := g.Dict[v1]
	if !ok1 {
		g.VertCount++
		e1 = g.VertCount - 1 //because of zero-based indexing
		g.Dict[v1] = e1
		g.Idict[e1] = v1
		g.AdjList = append(g.AdjList, List{nil, nil, nil})
	}

	e2, ok2 := g.Dict[v2]
	if !ok2 {
		g.VertCount++
		e2 = g.VertCount - 1 //because of zero-based indexing
		g.Dict[v2] = e2
		g.Idict[e2] = v2
		g.AdjList = append(g.AdjList, List{nil, nil, nil})
	}

	//Insert verteces into adjacency list
	newNode1 := &Node{Val: e2, Weight: weight, Prev: nil, Next: nil}
	g.AdjList[e1].insertNode(newNode1)

	//check to see if g is directed or if v1 = v2 (meaning you have a edge with the same vertex).
	//If so then return
	if g.directed || v1 == v2 {
		return
	}

	//if not then insert e1 into e2's list
	newNode2 := &Node{Val: e1, Weight: weight, Prev: nil, Next: nil}
	g.AdjList[e2].insertNode(newNode2)
}

//InsertNode inserts new node at end of linked list
func (l *List) insertNode(newNode *Node) error {
	if l.Head == nil {
		l.Head = newNode
		l.Tail = newNode
	} else {
		currentNode := l.Tail
		currentNode.Next = newNode
		newNode.Prev = l.Tail
	}
	l.Tail = newNode
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
	for i, v := range graph.AdjList {
		v.Current = v.Head
		//print the current vertex with the current index of the adjacency list
		fmt.Print(graph.Idict[i], ": ")
		//check is v.Current is nil to test for isolated nodes
		if v.Current == nil {
			fmt.Println()
			continue
		}
		//iterate over the neighbors of v
		for {
			fmt.Print("(", graph.Idict[v.Current.Val], ", ", *v.Current.Weight, ")")
			if v.Current == v.Tail {
				break
			}
			fmt.Print(", ")
			v.Current = v.Current.Next
		}
		//print newline and reset current to head. It seems like the polite thing to do.
		fmt.Println()
		v.Current = v.Head
	}
}

func printUnweightedGraph(graph *Graph) {

	for i, v := range graph.AdjList {
		v.Current = v.Head
		//print the current vertex with the current index of the adjacency list
		fmt.Print(graph.Idict[i], ": ")
		//check is v.Current is nil to test for isolated nodes
		if v.Current == nil {
			fmt.Println()
			continue
		}
		//iterate over the neighbors of v
		for {
			fmt.Print(graph.Idict[v.Current.Val])
			if v.Current == v.Tail {
				break
			}
			fmt.Print(", ")
			v.Current = v.Current.Next
		}
		//print newline and reset current to head. It seems like the polite thing to do.
		fmt.Println()
		v.Current = v.Head
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
