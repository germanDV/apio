package tags

// Repository defines the interface that a Tag Repository must implement.
type Repository interface {
	// Save persists a Tag.
	Save(tag TagAggregate) error
	// List fetches all Tags stored in the database.
	List() ([]TagAggregate, error)
}
