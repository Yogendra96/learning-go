package stl4go

import (
	"math/bits"
	"math/rand"
	"time"
)

const (
	skipListMaxLevel = 40
)

type SkipList[K any, V any] struct {
	level      int
	len        int
	head       skipListNode[K, V]
	prevsCache []*skipListNode[K, V]
	rander     *rand.Rand
	impl       skipListImpl[K, V]
}

func NewSkipList[K Ordered, V any]() *SkipList[K, V] {
	sl := skipListOrdered[K, V]{}
	sl.inti()
	sl.impl = (skipListimpl[K, V])(&sl)
	return &sl.SkipList
}

func NewSkipListFromMap[K Ordered, V any](m map[K]V) *SkipList[K, V] {
	sl := NewSkipList[K, V]()
	for k, v := range m {
		sl.Insert(k, v)
	}
	return sl
}

func NewSkipListFunc[K any, V any](keyCmp CompareFn[K]) *SkipList[K, V] {
	sl := skipListFunc[K, V]{}
	sl.init()
	sl.keyCmp = keyCmp
	sl.impl = skipListImpl[K, V](&sl)
	return &sl.SkipList
}

// IsEmpty implements the Container interface.
func (sl *SkipList[K, V]) IsEmpty() bool {
	return sl.len == 0
}

// Len implements the Container interface.
func (sl *SkipList[K, V]) Len() int {
	return sl.len
}

func (sl *SkipList[K, V]) Clear() {
	for i := range sl.head.next {
		sl.head.next[i] = nil
	}
	sl.level = 1
	sl.len = 0
}

func (sl *SkipList[K, V]) Iterate() MutableMapIterator[K, V] {
	return &skipListIterator[K, V]{sl.head.next[0], nil}
}

func (sl *SkipList[K, V]) Insert(key K, value V) {
	node, prevs := sl.impl.findInsertPoint(key)

	if node != nil {
		node.value = value
		return
	}

	level := sl.randomLevel()
	node = newSkipListNode(level, key, value)

	for i := 0; i < Min(level, sl.level); i++ { // 将节点插入到位置
		node.next[i] = prevs[i].next[i]
		prevs[i].next[i] = node
	}

	if level > sl.level {
		for i := sl.level; i < level; i++ {
			sl.head.next[i] = node
		}
		sl.level = level
	}
	sl.len++
}

func (sl *SkipList[K, V]) Find(key K) *V {
	node := sl.impl.findNode(key)
	if node != nil {
		return &node.value
	}
	return nil
}

func (sl *SkipList[K, V]) Has(key K) bool {
	return sl.impl.findNode(key) != nil
}

// LowerBound returns an iterator to the first element in the skiplist that
// does not satisfy element < value (i.e. greater or equal to),
// or a end iterator if no such element is found.
func (sl *SkipList[K, V]) LowerBound(key K) MutableMapIterator[K, V] {
	return &skipListIterator[K, V]{sl.impl.lowerBound(key), nil}
}

// UpperBound returns an iterator to the first element in the skiplist that
// does not satisfy value < element (i.e. strictly greater),
// or a end iterator if no such element is found.
func (sl *SkipList[K, V]) UpperBound(key K) MutableMapIterator[K, V] {
	return &skipListIterator[K, V]{sl.impl.upperBound(key), nil}
}

// FindRange returns an iterator in range [first, last) (last is not included).
func (sl *SkipList[K, V]) FindRange(first, last K) MutableMapIterator[K, V] {
	return &skipListIterator[K, V]{sl.impl.lowerBound(first), sl.impl.upperBound(last)}
}

func (sl *SkipList[K, V]) Remove(key K) bool {
	node, prevs := sl.impl.findRemovePoint(key)
	if node == nil {
		return false
	}
	for i, v := range node.next {
		prevs[i].next[i] = v
	}
	for sl.level > 1 && sl.head.next[sl.level-1] == nil {
		sl.level--
	}
	sl.len--
	return true
}

func (sl *SkipList[K, V]) ForEach(op func(K, V)) {
	for e := sl.head.next[0]; e != nil; e = e.next[0] {
		op(e.key, e.value)
	}
}

func (sl *SkipList[K, V]) ForEachMutable(op func(K, *V)) {
	for e := sl.head.next[0]; e != nil; e = e.next[0] {
		op(e.key, e.value)
	}
}

// ForEachIf implements the Map interface.
func (sl *SkipList[K, V]) ForEachIf(op func(K, V) bool) {
	for e := sl.head.next[0]; e != nil; e = e.next[0] {
		if !op(e.key, e.value) {
			return
		}
	}
}

// ForEachMutableIf implements the Map interface.
func (sl *SkipList[K, V]) ForEachMutableIf(op func(K, *V) bool) {
	for e := sl.head.next[0]; e != nil; e = e.next[0] {
		if !op(e.key, &e.value) {
			return
		}
	}
}

/// SkipList implementation part.

type skipListNode[K any, V any] struct {
	key   K
	value V
	next  []*skipListNode[K, V]
}

type skipListIterator[K any, V any] struct {
	node, end *skipListNode[K, V]
}

func (it *skipListIterator[K, V]) IsNotEnd() bool {
	return it.node != it.end
}

func (it *skipListIterator[K, V]) MoveToNext() {
	it.node = it.node.next[0]
}

func (it *skipListIterator[K, V]) Key() K {
	return it.node.key
}

func (it *skipListIterator[K, V]) Value() V {
	return it.node.value
}

func (it *skipListIterator[K, V]) Pointer() *V {
	return &it.node.value
}

// skipListImpl is an interface to provide different implementation for Ordered key or CompareFn.
//
// We can use CompareFn to compare Ordered keys, but a separated implementation is much faster.
// We don't make the whole skip list an interface, in order to share the type independented method.
// And because these methods are called directly without going through the interface, they are also
// much faster.
type skipListImpl[K any, V any] interface {
	findNode(key K) *skipListNode[K, V]
	lowerBound(key K) *skipListNode[K, V]
	upperBound(key K) *skipListNode[K, V]
	findInsertPoint(key K) (*skipListNode[K, V], []*skipListNode[K, V])
	findRemovePoint(key K) (*skipListNode[K, V], []*skipListNode[K, V])
}

func (sl *SkipList[K, V]) init() {
	sl.level = 1
	sl.rander = rand.New(rand.NewSource(time.Now().Unix()))
	sl.prevsCache = make([]*skipListNode[K, V], skipListMaxLevel)
	sl.head.next = make([]*skipListNode[K, V], skipListMaxLevel)
}

func (sl *SkipList[K, V]) randomLevel() int {
	total := uint64(1)<<uint64(skipListMaxLevel) - 1 // 2^n-1
	k := sl.rander.Uint64() & total
	level := skipListMaxLevel - bits.Len64(k) + 1
	for level > 3 && 1<<(level-3) > sl.len { //  找到一个合适的level
		level--
	}
	return level
}

// skipListOrdered is the skip list implementation for Ordered types.
type skipListOrdered[K Ordered, V any] struct {
	SkipList[K, V]
}

func (sl *skipListOrdered[K, V]) findNode(key K) *skipListNode[K, V] {
	return sl.doFindNode(key, true)
}

func (sl *skipListOrdered[K, V]) doFindNode(key K, eq bool) *skipListNode[K, V] {
	prev := &sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur := prev.next[i]; cur != nil; cur = cur.next[i] {
			if cur.key == key {
				return cur
			}
			if cur.key > key {
				break
			}
			prev = cur
		}
	}
	if eq {
		return nil
	}
	return prev.next[0]
}

func (sl *skipListOrdered[K, V]) lowerBound(key K) *skipListNode[K, V] {
	return sl.doFindNode(key, false)
}

func (sl *skipListOrdered[K, V]) upperBound(key K) *skipListNode[K, V] {
	node := sl.lowerBound(key)
	if node != nil && node.key == key {
		return node.next[0]
	}
	return node
}

// findInsertPoint returns (*node, nil) to the existed node if the key exists,
// or (nil, []*node) to the previous nodes if the key doesn't exist
func (sl *skipListOrdered[K, V]) findInsertPoint(key K) (*skipListNode[K, V], []*skipListNode[K, V]) {
	prevs := sl.prevsCache[0:sl.level]
	prev := &sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for next := prev.next[i]; next != nil; next = next.next[i] {
			if next.key == key {
				// The key is already existed, prevs are useless because no new node insertion.
				// stop searching.
				return next, nil
			}
			if next.key > key {
				// All other node in this level must be greater than the key,
				// search the next level.
				break
			}
			prev = next
		}
		prevs[i] = prev
	}
	return nil, prevs
}

// findRemovePoint finds the node which match the key and it's previous nodes.
func (sl *skipListOrdered[K, V]) findRemovePoint(key K) (*skipListNode[K, V], []*skipListNode[K, V]) {
	prevs := sl.findPrevNodes(key)
	node := prevs[0].next[0]
	if node == nil || node.key != key {
		return nil, nil
	}
	return node, prevs
}

func (sl *skipListOrdered[K, V]) findPrevNodes(key K) []*skipListNode[K, V] {
	prevs := sl.prevsCache[0:sl.level]
	prev := &sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for next := prev.next[i]; next != nil; next = next.next[i] {
			if next.key >= key {
				break
			}
			prev = next
		}
		prevs[i] = prev
	}
	return prevs
}

/// skipListFunc part

// skipListFunc is the skip list implementation which compare keys with func.
type skipListFunc[K any, V any] struct {
	SkipList[K, V]
	keyCmp CompareFn[K]
}

func (sl *skipListFunc[K, V]) findNode(key K) *skipListNode[K, V] {
	node := sl.lowerBound(key)
	if node != nil && sl.keyCmp(node.key, key) == 0 {
		return node
	}
	return nil
}

func (sl *skipListFunc[K, V]) lowerBound(key K) *skipListNode[K, V] {
	var prev = &sl.head
	for i := sl.level - 1; i >= 0; i-- {
		cur := prev.next[i]
		for ; cur != nil; cur = cur.next[i] {
			cmpRet := sl.keyCmp(cur.key, key)
			if cmpRet == 0 {
				return cur
			}
			if cmpRet > 0 {
				break
			}
			prev = cur
		}
	}
	return prev.next[0]
}

func (sl *skipListFunc[K, V]) upperBound(key K) *skipListNode[K, V] {
	node := sl.lowerBound(key)
	if node != nil && sl.keyCmp(node.key, key) == 0 {
		return node.next[0]
	}
	return node
}

// findInsertPoint returns (*node, nil) to the existed node if the key exists,
// or (nil, []*node) to the previous nodes if the key doesn't exist
func (sl *skipListFunc[K, V]) findInsertPoint(key K) (*skipListNode[K, V], []*skipListNode[K, V]) {
	prevs := sl.prevsCache[0:sl.level]
	prev := &sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur := prev.next[i]; cur != nil; cur = cur.next[i] {
			r := sl.keyCmp(cur.key, key)
			if r == 0 {
				// The key is already existed, prevs are useless because no new node insertion.
				// stop searching.
				return cur, nil
			}
			if r > 0 {
				// All other node in this level must be greater than the key,
				// search the next level.
				break
			}
			prev = cur
		}
		prevs[i] = prev
	}
	return nil, prevs
}

// findRemovePoint finds the node which match the key and it's previous nodes.
func (sl *skipListFunc[K, V]) findRemovePoint(key K) (*skipListNode[K, V], []*skipListNode[K, V]) {
	prevs := sl.findPrevNodes(key)
	node := prevs[0].next[0]
	if node == nil || sl.keyCmp(node.key, key) != 0 {
		return nil, nil
	}
	return node, prevs
}

func (sl *skipListFunc[K, V]) findPrevNodes(key K) []*skipListNode[K, V] {
	prevs := sl.prevsCache[0:sl.level]
	prev := &sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for next := prev.next[i]; next != nil; next = next.next[i] {
			if sl.keyCmp(next.key, key) >= 0 {
				break
			}
			prev = next
		}
		prevs[i] = prev
	}
	return prevs
}
