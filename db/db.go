package db

type DB interface {
	Open() error
	Close() error
}
