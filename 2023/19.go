package main

import (
	"fmt"
	"strconv"
	"strings"
)

type workflowRule struct {
	destinationWorkflowName string
	alwaysSends             bool
	attributeIndexToCompare int
	comparisonOperator      string
	comparisonConstant      int
}

// restricts attributes to satisfy and not satisfy this rule
func (rule workflowRule) restrict(x, m, a, s [4001]bool) ([4001]bool, [4001]bool, [4001]bool, [4001]bool, [4001]bool, [4001]bool, [4001]bool, [4001]bool) {
	assert(!rule.alwaysSends, "expected this rule to only sometimes send")

	var xA, mA, aA, sA, xR, mR, aR, sR [4001]bool

	switch rule.attributeIndexToCompare {
	case 0:
		mA, mR = m, m
		aA, aR = a, a
		sA, sR = s, s
	case 1:
		xA, xR = x, x
		aA, aR = a, a
		sA, sR = s, s
	case 2:
		xA, xR = x, x
		mA, mR = m, m
		sA, sR = s, s
	case 3:
		xA, xR = x, x
		mA, mR = m, m
		aA, aR = a, a
	default:
		panic("invalid attribute index")
	}

	if rule.comparisonOperator == ">" {
		switch rule.attributeIndexToCompare {
		case 0:
			// acceptable values
			for i := rule.comparisonConstant + 1; i <= 4000; i++ {
				xA[i] = x[i]
			}
			// rejectable values
			for i := 1; i <= rule.comparisonConstant; i++ {
				xR[i] = x[i]
			}
		case 1:
			for i := rule.comparisonConstant + 1; i <= 4000; i++ {
				mA[i] = m[i]
			}
			for i := 1; i <= rule.comparisonConstant; i++ {
				mR[i] = m[i]
			}
		case 2:
			for i := rule.comparisonConstant + 1; i <= 4000; i++ {
				aA[i] = a[i]
			}
			for i := 1; i <= rule.comparisonConstant; i++ {
				aR[i] = a[i]
			}
		case 3:
			for i := rule.comparisonConstant + 1; i <= 4000; i++ {
				sA[i] = s[i]
			}
			for i := 1; i <= rule.comparisonConstant; i++ {
				sR[i] = s[i]
			}
		default:
			panic("invalid attribute index to compare")
		}
	} else if rule.comparisonOperator == "<" {
		switch rule.attributeIndexToCompare {
		case 0:
			// acceptable values
			for i := 1; i < rule.comparisonConstant; i++ {
				xA[i] = x[i]
			}
			// rejectable values
			for i := rule.comparisonConstant; i <= 4000; i++ {
				xR[i] = x[i]
			}
		case 1:
			for i := 1; i < rule.comparisonConstant; i++ {
				mA[i] = m[i]
			}
			for i := rule.comparisonConstant; i <= 4000; i++ {
				mR[i] = m[i]
			}
		case 2:
			for i := 1; i < rule.comparisonConstant; i++ {
				aA[i] = a[i]
			}
			for i := rule.comparisonConstant; i <= 4000; i++ {
				aR[i] = a[i]
			}
		case 3:
			for i := 1; i < rule.comparisonConstant; i++ {
				sA[i] = s[i]
			}
			for i := rule.comparisonConstant; i <= 4000; i++ {
				sR[i] = s[i]
			}
		default:
			panic("invalid attribute index to compare")
		}
	} else {
		panic("invalid comparison operator")
	}

	return xA, mA, aA, sA, xR, mR, aR, sR
}

// true indexes are the allowed attribute values
// (we expect the first index to be always false)
func acceptableCombinations(x, m, a, s [4001]bool, workflowName string) int {
	// given the restrictions on the input which are passed as arguments, how many acceptable combinations are there when we feed the parts to the given workflow?

	// base case

	if workflowName == "A" {
		xCount := 0
		mCount := 0
		aCount := 0
		sCount := 0
		for _, v := range x {
			if v {
				xCount++
			}
		}
		for _, v := range m {
			if v {
				mCount++
			}
		}
		for _, v := range a {
			if v {
				aCount++
			}
		}
		for _, v := range s {
			if v {
				sCount++
			}
		}
		return xCount * mCount * aCount * sCount
	} else if workflowName == "R" {
		return 0
	}

	result := 0

	var xA, mA, aA, sA [4001]bool

	for _, rule := range workflows[workflowName] {
		if rule.alwaysSends {
			result += acceptableCombinations(x, m, a, s, rule.destinationWorkflowName)
			break
		}

		// now xmas are the rejectable combinations, and the xmasA are the acceptable combinations
		xA, mA, aA, sA, x, m, a, s = rule.restrict(x, m, a, s)

		// see how many combinations we can get by accepting this rule, and otherwise continue and assume this rule has not been accepted
		result += acceptableCombinations(xA, mA, aA, sA, rule.destinationWorkflowName)
	}

	return result

}

var workflows map[string][]workflowRule

func main() {
	lines := fetchLines(19)

	attributeToIndex := map[string]int{
		"x": 0,
		"m": 1,
		"a": 2,
		"s": 3,
	}

	////////////////////////////////////////////////

	var workflowsS []string
	var partsS []string

	// look for the empty line
	for i := range lines {
		if lines[i] == "" {
			workflowsS = lines[:i]
			partsS = lines[i+1:]
			break
		}
	}

	parts := make([][4]int, 0)
	for _, p := range partsS {
		uints := uintegers(p)
		assert(len(uints) == 4, "unexpected number of part attributes")
		parts = append(parts, [4]int{uints[0], uints[1], uints[2], uints[3]})
	}

	workflows = make(map[string][]workflowRule)
	for _, w := range workflowsS {
		nameAndFlow := strings.Split(w, "{")
		assert(len(nameAndFlow) == 2, "found not exactly one { in this workflow line")

		rules := make([]workflowRule, 0)
		for _, ruleString := range strings.Split(nameAndFlow[1][:len(nameAndFlow[1])-1], ",") {
			colonSplit := strings.Split(ruleString, ":")

			if len(colonSplit) == 1 {
				// always sends
				rules = append(rules, workflowRule{
					destinationWorkflowName: ruleString,
					alwaysSends:             true,
				})
			} else {
				assert(len(colonSplit) == 2, "what? not exactly 0 or 1 colons")

				splitGreater := strings.Split(colonSplit[0], ">")
				splitLess := strings.Split(colonSplit[0], "<")

				if len(splitGreater) == 2 {
					c, err := strconv.Atoi(splitGreater[1])
					check(err)

					rules = append(rules, workflowRule{
						destinationWorkflowName: colonSplit[1],
						alwaysSends:             false,
						attributeIndexToCompare: attributeToIndex[splitGreater[0]],
						comparisonOperator:      ">",
						comparisonConstant:      c,
					})
				} else if len(splitLess) == 2 {
					c, err := strconv.Atoi(splitLess[1])
					check(err)

					rules = append(rules, workflowRule{
						destinationWorkflowName: colonSplit[1],
						alwaysSends:             false,
						attributeIndexToCompare: attributeToIndex[splitLess[0]],
						comparisonOperator:      "<",
						comparisonConstant:      c,
					})
				} else {
					panic("if there is a colon, then there must be a < or > sign, right?")
				}

			}
		}

		workflows[nameAndFlow[0]] = rules
	}

	sum := 0
outer:
	for _, part := range parts {
		workflowName := "in"
	applyWorkflows:
		for {
			if workflowName == "A" {
				sum += part[0]
				sum += part[1]
				sum += part[2]
				sum += part[3]
				continue outer
			}
			if workflowName == "R" {
				continue outer
			}
			for _, rule := range workflows[workflowName] {
				if rule.alwaysSends {
					workflowName = rule.destinationWorkflowName
					continue applyWorkflows
				}

				switch rule.comparisonOperator {
				case ">":
					if part[rule.attributeIndexToCompare] > rule.comparisonConstant {
						workflowName = rule.destinationWorkflowName
						continue applyWorkflows
					}
				case "<":
					if part[rule.attributeIndexToCompare] < rule.comparisonConstant {
						workflowName = rule.destinationWorkflowName
						continue applyWorkflows
					}
				default:
					panic("invalid comparison operator")
				}
			}
			panic("we should never get here")
		}

	}
	fmt.Println(sum)

	////////////////////////////////////////////////

	var x, m, a, s [4001]bool
	for i := range x {
		x[i] = true
		m[i] = true
		a[i] = true
		s[i] = true
	}
	x[0] = false
	m[0] = false
	a[0] = false
	s[0] = false
	fmt.Println(acceptableCombinations(x, m, a, s, "in"))

}
