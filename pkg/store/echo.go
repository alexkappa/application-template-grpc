package store

// EchoStore abstracts the implementation details of the store needed for
type EchoStore interface {
	Echo(m string) (string, error)
}

// NewInMemoryEchoStore is an implementation of EchoStore that doesn't persist
// any data on disk. It's created for illustration purposes
func NewInMemoryEchoStore() EchoStore { return &inMemoryEchoStore{} }

type inMemoryEchoStore struct{}

func (s inMemoryEchoStore) Echo(m string) (string, error) { return m, nil }
