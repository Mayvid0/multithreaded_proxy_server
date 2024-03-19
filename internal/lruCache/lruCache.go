package lru

import "fmt"

type Node struct {
	data string
	prev *Node
	next *Node
}

type LRUCache struct {
	cache      map[string]*Node
	head, tail *Node
	size       int
}

func newNode(data string) *Node {
	return &Node{
		data: data,
		prev: nil,
		next: nil,
	}
}

func NewLRUCache(size int) *LRUCache {
	head := newNode("")
	tail := newNode("")
	head.next = tail
	tail.prev = head
	return &LRUCache{
		cache: make(map[string]*Node),
		head:  head,
		tail:  tail,
		size:  size,
	}
}

func (lru *LRUCache) insertNode(value string) *Node {
	curr := newNode(value)
	curr.next = lru.head.next
	lru.head.next = curr
	curr.next.prev = curr
	curr.prev = lru.head
	return curr
}

func (lru *LRUCache) deleteNode(nodeToDelete *Node, value string) *Node {
	if nodeToDelete == nil || nodeToDelete.prev == nil || nodeToDelete.next == nil {
		return nil
	}

	nodeToDelete.prev.next = nodeToDelete.next
	nodeToDelete.next.prev = nodeToDelete.prev

	// Delete the node from memory
	nodeToDelete = nil

	locationStoredAt := lru.insertNode(value)
	return locationStoredAt
}

func (lru *LRUCache) removeFromMap(nodeToBeReplaced *Node) {
	for key, node := range lru.cache {
		if node == nodeToBeReplaced {
			delete(lru.cache, key)
			break
		}
	}
}

func (lru *LRUCache) Get(key string) string {
	if node, ok := lru.cache[key]; ok {
		newLocation := lru.deleteNode(node, node.data)
		lru.cache[key] = newLocation
		return newLocation.data
	}
	return ""
}

func (lru *LRUCache) Put(key, value string) {
	if len(lru.cache) < lru.size {
		if node, ok := lru.cache[key]; ok {
			newLocation := lru.deleteNode(node, value)
			lru.cache[key] = newLocation
		} else {
			locationNodeIsStoredAt := lru.insertNode(value)
			lru.cache[key] = locationNodeIsStoredAt
		}
	} else {
		nodeToBeReplaced := lru.tail.prev
		locationNewNodeIsStoredAt := lru.deleteNode(nodeToBeReplaced, value)
		lru.removeFromMap(nodeToBeReplaced)
		lru.cache[key] = locationNewNodeIsStoredAt
	}
}

func (lru *LRUCache) Display() {
	current := lru.head.next
	for current != lru.tail {
		fmt.Printf("%s ", current.data)
		current = current.next
	}
	fmt.Println()
}
