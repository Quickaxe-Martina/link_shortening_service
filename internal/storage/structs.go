package storage

// Storage variables
type Storage struct {
	URLData map[string]string
}

// NewStorage create Storage
func NewStorage() *Storage {
	var storage = Storage{
		URLData: make(map[string]string),
	}
	return &storage
}
