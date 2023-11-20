package hash

type Hasher interface {
	Hash(in string) string
	GetHashSize() int
}
