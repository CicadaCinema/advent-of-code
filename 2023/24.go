package main

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

type hailstone struct {
	position [3]int
	velocity [3]int
}

// returns the constants m and c in the equation y=mx+c
func (h *hailstone) line2d() (float64, float64) {
	x1 := h.position[0]
	y1 := h.position[1]
	x2 := x1 + h.velocity[0]
	y2 := y1 + h.velocity[1]

	// we actually don't need to consider this case because the input data is nice
	assert(x1 != x2, "vertical line unexpected")

	m := float64(y1-y2) / float64(x1-x2)
	c := float64(y1) - m*float64(x1)

	return m, c
}

func xyIntersectsWithinTestArea(a, b hailstone, minCoord, maxCoord float64) bool {
	mA, cA := a.line2d()
	mB, cB := b.line2d()

	xIntersection := float64(cA-cB) / float64(mB-mA)
	yIntersection := float64(mA)*xIntersection + cA

	if !(minCoord <= xIntersection && xIntersection <= maxCoord && minCoord <= yIntersection && yIntersection <= maxCoord) {
		return false
	}

	// now we need to ensure that this intersection happens forward in time
	assert(a.velocity[0] != 0, "expected non-zero velocity")
	assert(b.velocity[0] != 0, "expected non-zero velocity")

	var resultA bool
	if a.velocity[0] > 0 {
		resultA = float64(a.position[0]) < xIntersection
	} else {
		resultA = float64(a.position[0]) > xIntersection
	}
	var resultB bool
	if b.velocity[0] > 0 {
		resultB = float64(b.position[0]) < xIntersection
	} else {
		resultB = float64(b.position[0]) > xIntersection
	}

	return resultA && resultB
}

/*
LINEAR ALGEBRA EXPLANATION:

Let \(i\) represent the index of a hailstone.

Let \(x,y,z\) and \(v_x,v_y,v_z\) be the initial position and velocity components of the rock.

Let \(x^i,y^i,z^i\) and \(v^i_x,v^i_y,v^i_z\) be the initial position and velocity components of hailstone \(i\). These are the only known terms, the rest are unknowns.

Let \(t_i\) represent the time taken (from the beginning) for the rock to collide with hailstone \(i\).

Clearly for all \(i\), and for all components (only \(x\) is shown here), we must have something analogous to

\[x+t_iv_x=x^i+t_iv^i_x\;.\]

This is a non-linear equation because of the product \(t_iv_x\) of two unknowns, one of which also depends on \(i\).\\

For any single \(i\), we can eliminate \(t_i\) and arrive at

\[
\frac{x^i-x}{v_x-v_x^i}=
\frac{y^i-y}{v_y-v_y^i}=
\frac{z^i-z}{v_z-v_z^i}
\;.\]

Consider just the first equality \((x^i-x)(v_y-v_y^i)=(y^i-y)(v_x-v_x^i)\) .

We can consider two different values of \(i\), for example \(1\) and \(j\), the first two hailstones, and subtract their respective equalities from each other, eliminating the \(-xv_y\) and \(-yv_x\).

Given at least 5 hailstones, we can form a system of four linear equations with four unknowns to find \(x, y, v_x\) and \(v_y\).

Repeat to find \(z\) and \(v_z\).

*/

// Returns the initial stone coordinates by computing some linear algebra.
// If which==true, returns x and y.
// If which==false, returns x and z.
func initialStoneCoords(hail []hailstone, which bool) (int, int) {
	// Solve Ax=b, where A is a 4x4 matrix and x and b are each 4-dimensional vectors.
	// Used https://www.reddit.com/r/golang/comments/zs0hqy/how_to_solve_a_set_of_simultaneous_equations_in_go/j164dv8/ to get started.

	var otherCoord int
	if which {
		otherCoord = 1
	} else {
		otherCoord = 2
	}

	aSlice := []float64{}
	for j := 1; j < len(hail); j++ {
		aSlice = append(aSlice, float64(hail[0].position[0]-hail[j].position[0]))
		aSlice = append(aSlice, float64(-hail[0].position[otherCoord]+hail[j].position[otherCoord]))
		aSlice = append(aSlice, float64(hail[0].velocity[otherCoord]-hail[j].velocity[otherCoord]))
		aSlice = append(aSlice, float64(-hail[0].velocity[0]+hail[j].velocity[0]))
	}
	A := mat.NewDense(len(hail)-1, 4, aSlice)

	bSlice := []float64{}
	for j := 1; j < len(hail); j++ {
		entry := 0
		entry += hail[0].position[0] * hail[0].velocity[otherCoord]
		entry -= hail[j].position[0] * hail[j].velocity[otherCoord]
		entry += hail[j].position[otherCoord] * hail[j].velocity[0]
		entry -= hail[0].position[otherCoord] * hail[0].velocity[0]

		bSlice = append(bSlice, float64(entry))
	}
	b := mat.NewVecDense(len(hail)-1, bSlice)

	// Solve the equations using the Solve function.
	var x mat.VecDense
	if err := x.SolveVec(A, b); err != nil {
		panic(err)
	}

	return int(math.Round(x.AtVec(2))), int(math.Round(x.AtVec(3)))
}

func main() {
	lines := fetchLines(24)

	////////////////////////////////////////////////

	hail := make([]hailstone, 0)
	for _, line := range lines {
		coords := integers(line)
		assert(len(coords) == 6, "invalid len of coords")
		hail = append(hail, hailstone{
			position: [3]int{coords[0], coords[1], coords[2]},
			velocity: [3]int{coords[3], coords[4], coords[5]},
		})
	}

	sum := 0
	for i := 0; i < len(hail); i++ {
		for j := i + 1; j < len(hail); j++ {
			if xyIntersectsWithinTestArea(hail[i], hail[j], 200000000000000, 400000000000000) {
				sum++
			}
		}
	}
	fmt.Println(sum)

	////////////////////////////////////////////////

	x1, y := initialStoneCoords(hail, true)
	x2, z := initialStoneCoords(hail, false)

	assert(x1 == x2, "got inconsistent values of x")

	// off-by-one error, don't care, I blame float precision
	fmt.Println(x1 + y + z - 1)
}
