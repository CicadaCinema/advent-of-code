package main

import (
	"cmp"
	"fmt"
	"regexp"
	"slices"

	"gonum.org/v1/gonum/stat/combin"
)

type graph struct {
	nodeToNeighbours map[string][]string
}

// walks the graph starting from an arbitrary node and returns the size of the connected subgraph conncted to that node
func (g *graph) sizeOfArbitraryConnectedComponent(bannedEdges map[[2]string]bool) int {
	visited := make(map[string]bool)
	queue := make([]string, 0)

	// populate the queue with an initial node
	for k := range g.nodeToNeighbours {
		queue = append(queue, k)
		visited[k] = true
		break
	}

	for len(queue) > 0 {
		// pop from the queue
		node := queue[0]
		queue = queue[1:]

		for _, connectedNode := range g.nodeToNeighbours[node] {
			if bannedEdges[normaliseEdge([2]string{node, connectedNode})] {
				continue
			}
			if !visited[connectedNode] {
				visited[connectedNode] = true
				queue = append(queue, connectedNode)
			}
		}
	}

	return len(visited)
}

func normaliseEdge(edge [2]string) [2]string {
	// the smallest index goes first
	if edge[1] < edge[0] {
		return [2]string{edge[1], edge[0]}
	} else {
		return edge
	}
}

func (g *graph) calculateEdgeCentralityFromNode(start string, edgeCentrality map[[2]string]int) {
	type incompletePath struct {
		currentNode  string
		visitedEdges [][2]string
	}

	// we want to visit every node, starting at `start`
	visitedNodes := make(map[string]bool)
	visitedNodes[start] = true

	// these are the nodes for which we have processed the edge centrality
	processedNodes := make(map[string]bool)
	processedNodes[start] = true

	frontier := []incompletePath{
		{
			currentNode:  start,
			visitedEdges: make([][2]string, 0),
		},
	}

	for len(visitedNodes) < len(g.nodeToNeighbours) {
		// pop from frontier
		thisPath := frontier[0]
		frontier = frontier[1:]

		// discover all the possible ways to extend thisPath
		for _, nextStep := range g.nodeToNeighbours[thisPath.currentNode] {
			// if we have visited this next node already, continue
			if visitedNodes[nextStep] {
				continue
			}
			visitedNodes[nextStep] = true

			if !processedNodes[nextStep] {
				// update edge centrality for all the edges on this path
				for _, edge := range thisPath.visitedEdges {
					edgeCentrality[edge] += 1
				}

				processedNodes[nextStep] = true
			}

			nextStepEdge := normaliseEdge([2]string{thisPath.currentNode, nextStep})
			newVisitedEdges := append(thisPath.visitedEdges, nextStepEdge)

			// we have not visited this next node, so we can create a new incomplete path
			newIncompletePath := incompletePath{
				currentNode:  nextStep,
				visitedEdges: newVisitedEdges,
			}

			frontier = append(frontier, newIncompletePath)
		}
	}

}

func cloneVisitedNodes(m map[string]bool) map[string]bool {
	result := make(map[string]bool)
	for k, v := range m {
		result[k] = v
	}
	return result
}

type incompletePath struct {
	visitedNodes map[string]bool
}

// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func remove(slice [][2]string, s int) [][2]string {
	return append(slice[:s], slice[s+1:]...)
}

func main() {
	lines := fetchLines(25)

	////////////////////////////////////////////////

	graph := graph{nodeToNeighbours: map[string][]string{}}
	nodes := make([]string, 0)
	edges := make([][2]string, 0)
	rWord := regexp.MustCompile("[a-z]+")

	// build the graph
	for _, line := range lines {
		elements := rWord.FindAllString(line, -1)
		assert(len(elements) >= 2, "bad elements len")

		origin := elements[0]
		if !slices.Contains(nodes, origin) {
			nodes = append(nodes, origin)
		}
		_, e1 := graph.nodeToNeighbours[origin]
		if !e1 {
			graph.nodeToNeighbours[origin] = make([]string, 0)
		}

		for _, dest := range elements[1:] {
			_, e2 := graph.nodeToNeighbours[dest]
			if !e2 {
				graph.nodeToNeighbours[dest] = make([]string, 0)
			}
			graph.nodeToNeighbours[origin] = append(graph.nodeToNeighbours[origin], dest)
			graph.nodeToNeighbours[dest] = append(graph.nodeToNeighbours[dest], origin)
			edges = append(edges, normaliseEdge([2]string{origin, dest}))
			if !slices.Contains(nodes, dest) {
				nodes = append(nodes, dest)
			}
		}
	}

	assert(len(nodes) == len(graph.nodeToNeighbours), "invalid graph contsruction")

	// calculate edge betweenness centrality for every edge
	edgeCentrality := make(map[[2]string]int)
	for _, e := range edges {
		edgeCentrality[e] = 0
	}
	for i := range nodes {
		graph.calculateEdgeCentralityFromNode(nodes[i], edgeCentrality)
	}

	// sort entries in reverse order
	type EdgeCentralityEntry struct {
		edge       [2]string
		centrality int
	}
	edgeCentralityEntries := make([]EdgeCentralityEntry, 0)
	for k, v := range edgeCentrality {
		edgeCentralityEntries = append(edgeCentralityEntries, EdgeCentralityEntry{edge: k, centrality: v})
	}
	slices.SortFunc(edgeCentralityEntries, func(a, b EdgeCentralityEntry) int { return -1 * cmp.Compare(a.centrality, b.centrality) })

	// find a combination of three edges, from some small number of edges with high centrality,
	// that partition the graph into disconnected components

	for _, combination := range combin.Combinations(10, 3) {
		// ban these three edges
		bannedEdges := make(map[[2]string]bool)

		for _, edgeIndex := range combination {
			bannedEdges[edgeCentralityEntries[edgeIndex].edge] = true
		}

		assert(len(bannedEdges) == 3, "bad combination")

		oneComponentSize := graph.sizeOfArbitraryConnectedComponent(bannedEdges)
		otherComponentSize := len(nodes) - oneComponentSize

		if 0 != oneComponentSize*otherComponentSize {
			fmt.Println(oneComponentSize * otherComponentSize)

		}

	}

	////////////////////////////////////////////////

}
