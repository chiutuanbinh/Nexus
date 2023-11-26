package common

type OrderedMap interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	LowerBound(key []byte) ([]byte, []byte, error)
	UpperBound(key []byte) ([]byte, []byte, error)
	KeyCount() int
}

type orderedMapImpl struct {
	avlTree AVLTree
}

// Delete implements OrderedMap.
func (o *orderedMapImpl) Delete(key []byte) error {
	return o.avlTree.Delete(key)
}

// Get implements OrderedMap.
func (o *orderedMapImpl) Get(key []byte) ([]byte, error) {
	return o.avlTree.Find(key)
}

// LowerBound implements OrderedMap.
func (o *orderedMapImpl) LowerBound(key []byte) ([]byte, []byte, error) {
	return o.avlTree.LowerBound(key)
}

// Put implements OrderedMap.
func (o *orderedMapImpl) Put(key []byte, value []byte) error {
	return o.avlTree.Insert(key, value)
}

// KeyCount implements OrderedMap.
func (o *orderedMapImpl) KeyCount() int {
	return o.avlTree.NodeCount()
}

// UpperBound implements OrderedMap.
func (o *orderedMapImpl) UpperBound(key []byte) ([]byte, []byte, error) {
	return o.avlTree.UpperBound(key)
}

func CreateOrderedMap() OrderedMap {
	return &orderedMapImpl{
		avlTree: CreateAVLTree(),
	}
}
