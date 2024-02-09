package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type moduleFlipflop struct {
	name         string
	state        bool
	destinations []string
}

type moduleConjunction struct {
	name                         string
	latestPulseReceivedFromInput map[string]bool
	destinations                 []string
}

var flipflops map[string]*moduleFlipflop
var conjunctions map[string]*moduleConjunction

type pulse struct {
	input       string
	bit         bool
	destination string
}

var pulseQueue []pulse

func (f *moduleFlipflop) pulse(pulseType bool) {
	if pulseType {
		// if a flip-flop module receives a high pulse, it is ignored and nothing happens
		highPulses++
		return
	}
	lowPulses++

	// if a flip-flop module receives a low pulse, it flips between on and off
	f.state = !f.state

	// If it was off, it turns on and sends a high pulse. If it was on, it turns off and sends a low pulse.
	for _, d := range f.destinations {
		pulseQueue = append(pulseQueue, pulse{
			input:       f.name,
			bit:         f.state,
			destination: d,
		})
	}
}

func (c *moduleConjunction) pulse(pulseType bool) {
	if pulseType {
		highPulses++
	} else {
		lowPulses++
	}

	for _, v := range c.latestPulseReceivedFromInput {
		// if at least one is low, send a high pulse
		if !v {
			for _, d := range c.destinations {
				pulseQueue = append(pulseQueue, pulse{
					input:       c.name,
					bit:         true,
					destination: d,
				})
			}
			return
		}
	}

	// if it remembers high pulses for all inputs, it sends a low pulse
	for _, d := range c.destinations {
		pulseQueue = append(pulseQueue, pulse{
			input:       c.name,
			bit:         false,
			destination: d,
		})
	}
}

var lowPulses, highPulses int

func main() {
	lines := fetchLines(20)
	pulseQueue = make([]pulse, 0)
	flipflops = make(map[string]*moduleFlipflop)
	conjunctions = make(map[string]*moduleConjunction)

	////////////////////////////////////////////////

	lowPulses = 1000
	highPulses = 0
	rWords := regexp.MustCompile("[a-z]+")

	starter := make([]pulse, 0)

	names1 := make([]string, 0)
	names2 := make([]string, 0)
	names3 := make([]string, 0)

	foundBroadcaster := false
	for _, line := range lines {
		if line[:5] == "broad" {
			assert(!foundBroadcaster, "found two broadcasters")
			foundBroadcaster = true
			for _, dest := range strings.Split(line[15:], ", ") {
				starter = append(starter, pulse{
					input:       "broadcaster",
					bit:         false,
					destination: dest,
				})
			}
			continue
		}

		words := rWords.FindAllString(line, -1)
		names3 = append(names3, words[0])

		switch line[0] {
		case '%':
			names1 = append(names1, words[0])
			flipflops[words[0]] = &moduleFlipflop{
				name:         words[0],
				state:        false,
				destinations: words[1:],
			}
		case '&':
			names2 = append(names2, words[0])
			conjunctions[words[0]] = &moduleConjunction{
				name:                         words[0],
				latestPulseReceivedFromInput: map[string]bool{},
				destinations:                 words[1:],
			}
		default:
			panic("unexpected module")
		}
	}
	assert(foundBroadcaster, "did not find broadcaster")

	// populate conjunction inputs
	for name, c := range conjunctions {
		for otherName, otherC := range conjunctions {
			// if this conjunction is a destination of other, then otherName is an input
			if slices.Contains(otherC.destinations, name) {
				c.latestPulseReceivedFromInput[otherName] = false
			}
		}

		for otherName, otherF := range flipflops {
			if slices.Contains(otherF.destinations, name) {
				c.latestPulseReceivedFromInput[otherName] = false
			}
		}
	}

	// zp is the only input of rx
	inputToPeriod := make(map[string]int)
	for k, v := range conjunctions {
		if slices.Contains(v.destinations, "zp") {
			inputToPeriod[k] = -1
		}
	}
	mapEntriesToPopulate := len(inputToPeriod)

	for i := 1; i < 100000000000; i++ {
		assert(len(pulseQueue) == 0, "queue not empty")
		for _, p := range starter {
			pulseQueue = append(pulseQueue, p)
		}
		for len(pulseQueue) > 0 {
			// pop from the queue
			thisPulse := pulseQueue[0]
			pulseQueue = pulseQueue[1:]
			// fmt.Println(thisPulse)

			flipFlop, existsF := flipflops[thisPulse.destination]
			conjunction, existsC := conjunctions[thisPulse.destination]
			if existsF {
				flipFlop.pulse(thisPulse.bit)
			} else if existsC {
				// when a pulse is received, the conjunction module first updates its memory for that input
				conjunction.latestPulseReceivedFromInput[thisPulse.input] = thisPulse.bit
				conjunction.pulse(thisPulse.bit)

				if inputToPeriod[thisPulse.destination] == -1 && !thisPulse.bit {
					inputToPeriod[thisPulse.destination] = i
					mapEntriesToPopulate--
				}
			} else {
				if thisPulse.bit {
					highPulses++
				} else {
					lowPulses++
				}
			}
		}
		if i == 1000 {
			fmt.Println(lowPulses * highPulses)
		}
		if mapEntriesToPopulate == 0 && i > 1000 {
			break
		}
	}

	bigPeriod := 1
	for _, v := range inputToPeriod {
		bigPeriod = lcm(bigPeriod, v)
	}

	fmt.Println(bigPeriod)

	////////////////////////////////////////////////

}
