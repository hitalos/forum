package pg

import (
	"crg.eti.br/go/forum/config"
	"crg.eti.br/go/forum/db"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Pg struct {
	db *sqlx.DB
}

func New() db.DB {
	db := &Pg{}
	return db
}

func (p *Pg) Open() error {
	var err error
	p.db, err = sqlx.Open("postgres", config.DBURL)
	if err != nil {
		return err
	}

	return p.db.Ping()
}

func (p *Pg) Close() error {
	return p.db.Close()
}
