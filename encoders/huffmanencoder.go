package encoders

import (
	"container/heap"
)

// Huffman returns a PrefixEncoder struct with an Alphabet slice and Codewords map.
func Huffman(alphabet []rune, corpus string) PrefixEncoder {
	frequencies := characterFrequencies(corpus)
	tree := makeHuffmanTree(alphabet, frequencies)
	codewords := makeCodewords(tree)

	return PrefixEncoder{
		Alphabet:  alphabet,
		Codewords: codewords,
	}
}

// characterFrequencies returns a map from codepoints to their frequency in text.
func characterFrequencies(text string) map[rune]int {
	frequencies := make(map[rune]int)

	for _, character := range text {
		frequencies[character]++
	}

	return frequencies
}

// Node for constructing Huffman tree.
type Node struct {
	character rune
	frequency int
	left      *Node
	right     *Node
	index     int // Necessary for priority queue
}

func makeCodewords(huffmanTree *Node) map[rune]string {
	codewords := make(map[rune]string)
	huffmanTree.traverseHuffmanTree(codewords, "")
	return codewords
}

func (root *Node) traverseHuffmanTree(codewords map[rune]string, codeword string) {
	if root != nil {
		// Define codewords mapping if a character node is found
		if root.character != -1 {
			codewords[root.character] = codeword
		}
		// Traverse Huffman tree and append 0s and 1s to codeword along the way
		root.left.traverseHuffmanTree(codewords, codeword+"0")
		root.right.traverseHuffmanTree(codewords, codeword+"1")
	}
}

func makeHuffmanTree(alphabet []rune, frequencies map[rune]int) *Node {
	nodes := make([]*Node, len(alphabet))

	// Create list of nodes from characters in alphabet and their frequencies
	for i, character := range alphabet {
		nodes[i] = &Node{
			character: character,
			frequency: frequencies[character],
			left:      nil,
			right:     nil,
		}
	}

	queue := make(PriorityQueue, 0)

	// Add nodes to priority queue sorted by frequency
	for _, node := range nodes {
		heap.Push(&queue, node)
	}

	for queue.Len() > 1 {
		// Pop two least frequent characters
		node1 := heap.Pop(&queue).(*Node)
		node2 := heap.Pop(&queue).(*Node)

		// Add nodes to parent node with combined children frequency
		mergedNode := &Node{
			character: -1,
			frequency: node1.frequency + node2.frequency,
			left:      node1,
			right:     node2,
		}

		heap.Push(&queue, mergedNode)
	}

	// Return root of completed Huffman tree
	return heap.Pop(&queue).(*Node)
}

// PriorityQueue implementation from Go documentation.
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].frequency < pq[j].frequency
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Node)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}
