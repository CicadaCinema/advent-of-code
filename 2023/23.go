package main

import "fmt"

var lines []string

type route struct {
	currentLocation [2]int
	visitedCells    map[[2]int]bool
	routeComplete   bool
	noRemove        bool
}

func allowedNextSteps1(routeSoFar route) [][2]int {
	result := make([][2]int, 0)
	for _, candidateDirection := range [][2]int{[2]int{1, 0}, [2]int{0, 1}, [2]int{-1, 0}, [2]int{0, -1}} {
		newY := routeSoFar.currentLocation[0] + candidateDirection[0]
		newX := routeSoFar.currentLocation[1] + candidateDirection[1]

		// ensure y is in bounds
		if newY < 0 || newY >= len(lines) {
			continue
		}

		// ensure x is in bounds
		if newX < 0 || newX >= len(lines[0]) {
			continue
		}

		// ensure we stay on the path
		switch lines[newY][newX] {
		case '#':
			continue
		case '>':
			if candidateDirection != [2]int{0, 1} {
				continue
			}
		case '<':
			if candidateDirection != [2]int{0, -1} {
				continue
			}
		case '^':
			if candidateDirection != [2]int{-1, 0} {
				continue
			}
		case 'v':
			if candidateDirection != [2]int{1, 0} {
				continue
			}
		}

		// clearly we can't visit the same cell twice
		if routeSoFar.visitedCells[[2]int{newY, newX}] {
			continue
		}

		result = append(result, candidateDirection)

	}
	return result
}

func allowedNextSteps2(currentLocation [2]int, visitedCells map[[2]int]bool) [][2]int {
	result := make([][2]int, 0)
	for _, candidateDirection := range [][2]int{[2]int{1, 0}, [2]int{0, 1}, [2]int{-1, 0}, [2]int{0, -1}} {
		newY := currentLocation[0] + candidateDirection[0]
		newX := currentLocation[1] + candidateDirection[1]

		// ensure y is in bounds
		if newY < 0 || newY >= len(lines) {
			continue
		}

		// ensure x is in bounds
		if newX < 0 || newX >= len(lines[0]) {
			continue
		}

		// ensure we stay on the path
		if lines[newY][newX] == '#' {
			continue
		}

		// clearly we can't visit the same cell twice
		if visitedCells[[2]int{newY, newX}] {
			continue
		}

		result = append(result, candidateDirection)

	}
	return result
}

// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func remove(slice []route, s int) []route {
	return append(slice[:s], slice[s+1:]...)
}

func cloneVisitedCells(m map[[2]int]bool) map[[2]int]bool {
	result := make(map[[2]int]bool)
	for k, v := range m {
		result[k] = v
	}
	return result
}

// a node is either a fork in the road, or an element of the first/last row
func isNode(cell [2]int) bool {
	assert(lines[cell[0]][cell[1]] != '#', "we should not be at a forest tile")

	if cell[0] == 0 || cell[0] == len(lines)-1 {
		return true
	}

	visited := make(map[[2]int]bool)
	visited[cell] = true

	allowedCount := len(allowedNextSteps2(cell, visited))
	assert(allowedCount > 0, "cannot go anywhere")

	// a node is a FORK in the road
	return allowedCount >= 3
}

// given a node position which is at an intersection, return a slice containing all the nodes it is connected to
func findAdjacentNodes(node [2]int) ([][2]int, []int) {
	result := make([][2]int, 0)
	lengths := make([]int, 0)

	visited := make(map[[2]int]bool)
	visited[node] = true

	possibleFirstSteps := allowedNextSteps2(node, visited)
	assert(node[0] == 0 || node[0] == len(lines)-1 || len(possibleFirstSteps) > 1, "bad node")

	// take the first step in each direction
	routes := []route{}
	for _, firstStep := range possibleFirstSteps {
		newPos := [2]int{node[0] + firstStep[0], node[1] + firstStep[1]}
		newVisited := cloneVisitedCells(visited)
		newVisited[newPos] = true
		routes = append(routes, route{
			currentLocation: newPos,
			visitedCells:    newVisited,
		})
	}

	// travel in the direction of each route until we hit a location where != 1 step is possible
	// and assert that we have hit a node
	for _, route := range routes {
		// repeatedly take steps
		for {
			possibleNextSteps := allowedNextSteps2(route.currentLocation, route.visitedCells)
			if len(possibleNextSteps) != 1 {
				break
			}
			route.currentLocation = [2]int{route.currentLocation[0] + possibleNextSteps[0][0], route.currentLocation[1] + possibleNextSteps[0][1]}
			route.visitedCells[route.currentLocation] = true
		}

		// so the current cell is a node
		assert(isNode(route.currentLocation), "expected a node here")
		result = append(result, route.currentLocation)
		lengths = append(lengths, len(route.visitedCells)-2)
	}

	assert(len(result) == len(lengths), "wrong len")
	return result, lengths
}

type destNodeAndLength struct {
	dest   [2]int
	length int
}

type graph struct {
	// the lengths are the lengths of the passages between the nodes, not counting the cells which are the nodes themselves
	edgeLengths map[[2][2]int]int

	nodeToDest map[[2]int][]destNodeAndLength
}

func makeGraph() graph {
	graph := graph{
		edgeLengths: map[[2][2]int]int{},
		nodeToDest:  map[[2]int][]destNodeAndLength{},
	}

	nodesFrontier := [][2]int{[2]int{0, 1}}
	visitedNodes := make(map[[2]int]bool)

	for len(nodesFrontier) > 0 {
		// pop first node
		node := nodesFrontier[0]
		nodesFrontier = nodesFrontier[1:]

		// if we have already visited this node, there is nothing to do
		if visitedNodes[node] {
			continue
		}
		visitedNodes[node] = true

		// see which are the adjacent nodes
		adjacentNodes, adjacentLengths := findAdjacentNodes(node)
		nodesFrontier = append(nodesFrontier, adjacentNodes...)

		// iterate over every node, adjacentNodes[i] pair
		for i := range adjacentNodes {
			pair := [2][2]int{node, adjacentNodes[i]}
			pairInv := [2][2]int{adjacentNodes[i], node}
			length, exists := graph.edgeLengths[pair]
			_, existsInv := graph.edgeLengths[pairInv]
			if exists {
				assert(existsInv, "inv doesn't exist")
				// assert that the length matches
				assert(adjacentLengths[i] == length, "wrong length")
			} else {
				assert(!existsInv, "inv exists")

				// here is where we actually make the assignments to the graph object!
				graph.edgeLengths[pair] = adjacentLengths[i]
				graph.edgeLengths[pairInv] = adjacentLengths[i]

				if _, exists1 := graph.nodeToDest[node]; !exists1 {
					graph.nodeToDest[node] = make([]destNodeAndLength, 0)
				}
				if _, exists2 := graph.nodeToDest[adjacentNodes[i]]; !exists2 {
					graph.nodeToDest[adjacentNodes[i]] = make([]destNodeAndLength, 0)
				}

				graph.nodeToDest[node] = append(graph.nodeToDest[node], destNodeAndLength{
					dest:   adjacentNodes[i],
					length: adjacentLengths[i],
				})
				graph.nodeToDest[adjacentNodes[i]] = append(graph.nodeToDest[adjacentNodes[i]], destNodeAndLength{
					dest:   node,
					length: adjacentLengths[i],
				})
			}
		}
	}

	return graph
}

type graphPath struct {
	currentNode  [2]int
	visitedNodes map[[2]int]bool
	length       int
	complete     bool
}

func cloneVisitedNodes(m map[[2]int]bool) map[[2]int]bool {
	result := make(map[[2]int]bool)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func (g *graph) extendPath(path graphPath) []graphPath {
	result := make([]graphPath, 0)

	for _, dest := range g.nodeToDest[path.currentNode] {
		// if we have already visited dest, can't proceed
		if path.visitedNodes[dest.dest] {
			continue
		}

		// now we know that we can visit this next node
		newVisitedNodes := cloneVisitedNodes(path.visitedNodes)
		newVisitedNodes[dest.dest] = true
		newGraphPath := graphPath{
			currentNode:  dest.dest,
			visitedNodes: newVisitedNodes,
			length:       path.length + dest.length + 1,
			complete:     true,
		}

		// if at least one of the neighbours of the dest is not yet visited, it is not a complete path
		for _, furtherDest := range g.nodeToDest[dest.dest] {
			if !newGraphPath.visitedNodes[furtherDest.dest] {
				newGraphPath.complete = false
				break
			}
		}

		result = append(result, newGraphPath)
	}

	return result
}

func main() {
	lines = fetchLines(23)

	////////////////////////////////////////////////

	// maze traversal cell by cell

	incompleteRoutes := []route{route{
		currentLocation: [2]int{0, 1},
		visitedCells:    make(map[[2]int]bool),
	}}
	incompleteRoutes[0].visitedCells[incompleteRoutes[0].currentLocation] = true

	var possibleNextSteps int
	for {
		possibleNextSteps = 0

		maxIndex := len(incompleteRoutes) - 1
		for i := range incompleteRoutes {
			if incompleteRoutes[i].routeComplete {
				continue
			}
			nextDirections := allowedNextSteps1(incompleteRoutes[i])
			possibleNextSteps += len(nextDirections)
			if len(nextDirections) == 0 {
				incompleteRoutes[i].routeComplete = true
			}
			for j := range nextDirections {
				nextPosition := [2]int{incompleteRoutes[i].currentLocation[0] + nextDirections[j][0], incompleteRoutes[i].currentLocation[1] + nextDirections[j][1]}
				visitedCellsClone := cloneVisitedCells(incompleteRoutes[i].visitedCells)
				visitedCellsClone[nextPosition] = true
				incompleteRoutes = append(incompleteRoutes, route{
					currentLocation: nextPosition,
					visitedCells:    visitedCellsClone,
				})
			}
		}
		for i := maxIndex; i >= 0; i-- {
			if incompleteRoutes[i].routeComplete {
				continue
			}
			incompleteRoutes = remove(incompleteRoutes, i)
		}

		if possibleNextSteps == 0 {
			break
		}
	}

	longestRoute := -1
	for _, route := range incompleteRoutes {
		if route.currentLocation[0] == len(lines)-1 && len(route.visitedCells) > longestRoute {
			longestRoute = len(route.visitedCells)
		}
	}
	fmt.Println(longestRoute - 1)

	////////////////////////////////////////////////

	// graph traversal
	// this works okay, but takes a little bit of time and requires a lot of memory
	// a better approach would be to do depth first search instead of breadth first search
	// but ultimately this works fine and produces the correct answer

	graph := makeGraph()

	incompleteGraphRoutes := []graphPath{graphPath{
		currentNode:  [2]int{0, 1},
		visitedNodes: map[[2]int]bool{},
		length:       0,
		complete:     false,
	}}
	incompleteGraphRoutes[0].visitedNodes[incompleteGraphRoutes[0].currentNode] = true

	maxLength := -1
	for len(incompleteGraphRoutes) > 0 {
		// pop off the first incomplete route
		route := incompleteGraphRoutes[0]
		incompleteGraphRoutes = incompleteGraphRoutes[1:]

		continuations := graph.extendPath(route)
		for _, c := range continuations {
			if c.complete && c.length > maxLength {
				maxLength = c.length
			}
			if !c.complete {
				incompleteGraphRoutes = append(incompleteGraphRoutes, c)
			}
		}
	}

	fmt.Println(maxLength)
}
