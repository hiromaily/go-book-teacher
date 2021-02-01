package storage

// Storager interface
type Storager interface {
	Save(string) (bool, error)
	Delete() error
	Close()
}
