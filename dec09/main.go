package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"
)

//go:embed input.txt
var input string

const emptyVal = -1

func makeInitialLayout(s string) []int {
	var layout []int

	empty := false
	fileID := 0
	for _, ch := range s {
		val := int(ch - '0')
		appendVal := emptyVal
		if !empty {
			appendVal = fileID
			fileID++
		}
		for i := 0; i < val; i++ {
			layout = append(layout, appendVal)
		}
		empty = !empty
	}
	return layout
}

func part1(layout []int) int {
	start := -1
	end := len(layout)

	moveStart := func() {
		for {
			start++
			if start == len(layout) {
				break
			}
			if layout[start] == emptyVal {
				break
			}
		}
	}
	moveEnd := func() {
		for {
			end--
			if end < 0 {
				break
			}
			if layout[end] != emptyVal {
				break
			}
		}
	}

	for {
		moveStart()
		moveEnd()
		if start > end {
			break
		}
		layout[start], layout[end] = layout[end], layout[start]
	}
	sum := 0
	for i := 0; i < len(layout); i++ {
		if layout[i] == emptyVal {
			continue
		}
		sum += i * layout[i]
	}
	return sum
}

type block struct {
	fileID     int
	start, end int
	next       *block
	prev       *block
}

func (b *block) isEmpty() bool {
	return b.fileID == emptyVal
}

func (b *block) length() int {
	return b.end - b.start + 1
}

type blockList struct {
	head *block
	tail *block
}

func (bl *blockList) dump() {
	node := bl.head
	for node != nil {
		if node.isEmpty() {
			fmt.Printf(" E:%d", node.length())
		} else {
			fmt.Printf(" %d:%d", node.fileID, node.length())
		}
		node = node.next
	}
	fmt.Println()
}

func (bl *blockList) add(b *block) {
	if bl.head == nil {
		bl.head = b
		bl.tail = b
		return
	}
	bl.tail.next = b
	b.prev = bl.tail
	bl.tail = b
}

func (bl *blockList) fill(b *block, fileID int, byLength int) {
	if !b.isEmpty() {
		panic(fmt.Sprintf("fill non-empty block: %+v", b))
	}
	if byLength > b.length() {
		panic("block too short")
	}
	b2 := &block{fileID: emptyVal, start: b.start + byLength, end: b.end}
	b.fileID = fileID
	b.end = b2.start - 1
	bl.insertAfter(b, b2)
}

func (bl *blockList) insertAfter(b *block, newBlock *block) {
	newBlock.prev = b
	newBlock.next = b.next
	if b.next == nil {
		bl.tail = newBlock
	}
	b.next = newBlock
}

func part2(layout []int) int {
	list := &blockList{}
	b := &block{start: 0, end: 0, fileID: layout[0]}
	for i := 1; i < len(layout); i++ {
		changed := b.fileID != layout[i]
		if changed {
			list.add(b)
			b = &block{fileID: layout[i], start: i, end: i}
			continue
		}
		b.end = i
	}
	// add the last block
	list.add(b)

	// start from right
	candidate := list.tail
	for candidate != nil {
		if candidate.isEmpty() {
			candidate = candidate.prev
			continue
		}
		// look for potential free space from left
		potential := list.head
		for potential != nil {
			if potential == candidate {
				break
			}
			if !potential.isEmpty() {
				potential = potential.next
				continue
			}
			if potential.length() >= candidate.length() {
				list.fill(potential, candidate.fileID, candidate.length())
				candidate.fileID = emptyVal
				break
			}
			potential = potential.next
		}
		candidate = candidate.prev
	}

	node := list.head
	sum := 0
	for node != nil {
		if !node.isEmpty() {
			for i := node.start; i <= node.end; i++ {
				sum += i * node.fileID
			}
		}
		node = node.next
	}
	return sum
}

func main() {
	input = strings.TrimSpace(input)
	log.Println("PART1 SUM:", part1(makeInitialLayout(input)))
	log.Println("PART2 SUM:", part2(makeInitialLayout(input)))
}
