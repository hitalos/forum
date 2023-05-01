package db

type Topic struct {
	ID          int     `db:"id" json:"id"`
	ParentID    int     `db:"parent_id" json:"parent_id,omitempty"`
	Level       int     `db:"level" json:"level,omitempty"`
	ZOrder      int     `db:"zorder" json:"zorder,omitempty"`
	Subtopic    []Topic `db:"subtopic" json:"subtopic,omitempty"`
	Title       string  `db:"title" json:"title"`
	Description string  `db:"description" json:"description,omitempty"`
	CreatedAt   string  `db:"created_at" json:"created_at,omitempty"`
}

type DB interface {
	Open() error
	Close() error
	ListTopics(parent_id int) ([]Topic, error)
}
