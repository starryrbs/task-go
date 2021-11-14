package data

import (
	"github.com/task-go/week04/internal/data/ent"
)
import "github.com/task-go/week04/internal/conf"

type Data struct {
	db *ent.Client
}

func NewData(database *database.Database) (*Data, error) {
	client, err := ent.Open(database.Driver, database.Source)
	if err != nil {
		return nil, err
	}
	return &Data{db: client}, err
}
