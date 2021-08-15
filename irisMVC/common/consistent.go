package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type units []uint32

func (x units) Len() int {
	return len(x)
}

func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

//swap values
func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

//if no data on circle
var errEmpty = errors.New("hash has no data")

// Consistent struct to store data
type Consistent struct {
	//hash circle，key is hash code，value to store node value
	circle map[uint32]string
	//hash slice
	sortedHashes units
	//v node for node balancing
	VirtualNode int
	//rw lock
	sync.RWMutex
}

//default nodes
func NewConsistent() *Consistent {
	return &Consistent{
		circle: make(map[uint32]string),
		VirtualNode: 20,
	}
}

func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

//get hash location
func (c *Consistent) hashkey(key string) uint32 {
	if len(key) < 64 {
		var srcatch [64]byte
		copy(srcatch[:], key)
		//return CRC-32 sum
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashes = nil
	}

	//add hashes
	for k := range c.circle {
		hashes = append(hashes, k)
	}

	//sort hash for binary search
	sort.Sort(hashes)
	c.sortedHashes = hashes
}

// Add add node to hash circle
func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

//add nodes
func (c *Consistent) add(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		//add to hash circle
		c.circle[c.hashkey(c.generateKey(element, i))] = element
	}
	//update sorting
	c.updateSortedHashes()
}

//remove removing node
func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashkey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

//search nearest node
func (c *Consistent) search(key uint32) int {
	//search
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	//binary search closest node
	i := sort.Search(len(c.sortedHashes), f)
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

// Get get nearest node info by data string
func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	if len(c.circle) == 0 {
		return "", errEmpty
	}
	//get hash key
	key := c.hashkey(name)
	//search node by hash key
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}