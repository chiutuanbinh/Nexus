package common

type OrderedMap interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	//Find the first key-value pair which key is equal or larger than key
	LowerBound(key []byte) ([]byte, []byte, error)
	//Find the first key-value pair which key is larger than key
	UpperBound(key []byte) ([]byte, []byte, error)
	KeyCount() int
}

type orderedMapImpl struct {
	bstTree BstTree
}

// Delete implements OrderedMap.
func (o *orderedMapImpl) Delete(key []byte) error {
	return o.bstTree.Delete(key)
}

// Get implements OrderedMap.
func (o *orderedMapImpl) Get(key []byte) ([]byte, error) {
	return o.bstTree.Find(key)
}

// LowerBound implements OrderedMap.
func (o *orderedMapImpl) LowerBound(key []byte) ([]byte, []byte, error) {
	return o.bstTree.LowerBound(key)
}

// Put implements OrderedMap.
func (o *orderedMapImpl) Put(key []byte, value []byte) error {
	return o.bstTree.Insert(key, value)
}

// KeyCount implements OrderedMap.
func (o *orderedMapImpl) KeyCount() int {
	return o.bstTree.NodeCount()
}

// UpperBound implements OrderedMap.
func (o *orderedMapImpl) UpperBound(key []byte) ([]byte, []byte, error) {
	return o.bstTree.UpperBound(key)
}

func CreateOrderedMap() OrderedMap {
	return &orderedMapImpl{
		bstTree: CreateAVLTree(),
	}
}
