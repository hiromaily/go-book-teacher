package storages

// Saver is to save
type Saver interface {
	Save(string) bool
}

// Deleter is to Notice
type Deleter interface {
	Delete() error
}

// SaveDeleter is to Save and Delete
type SaveDeleter interface {
	Saver
	Deleter
}
