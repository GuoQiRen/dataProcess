package local

type AtomicId interface {
	GetAtomicId(atomType int) (id int32, err error)
}
