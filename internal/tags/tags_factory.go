package tags

import "github.com/germandv/apio/internal/id"

// FromReq produces a TagAggregate from a Request.
func FromReq(reqName string) (TagAggregate, error) {
	name, err := ParseName(reqName)
	if err != nil {
		return TagAggregate{}, err
	}

	uid := id.New()
	t := TagAggregate{
		TagEntity: TagEntity{ID: uid, Name: name},
		NoteCount: 0,
	}

	return t, nil
}

// FromDB produces a TagAggregate from data in the DB format.
func FromDB(dbID string, dbName string) (TagAggregate, error) {
	uid, err := id.Parse(dbID)
	if err != nil {
		return TagAggregate{}, err
	}

	name, err := ParseName(dbName)
	if err != nil {
		return TagAggregate{}, err
	}

	t := TagAggregate{
		TagEntity: TagEntity{ID: uid, Name: name},
		NoteCount: 0, // TODO: implement
	}

	return t, nil
}
