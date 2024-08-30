package notes

// Repository defines the interface that a Note Repository must implement.
type Repository interface {
	// Save persists a Note.
	Save(NoteAggregate) error
	// List fetches all Notes stored in the database.
	List() ([]NoteAggregate, error)
	// TagsExist checks whether all provided tags exist in the database.
	TagsExist([]NoteTagEntity) (bool, error)
}
