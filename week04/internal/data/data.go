package data

import (
	"context"
	"github.com/google/wire"
	_ "github.com/mattn/go-sqlite3"
	"github.com/task-go/week04/internal/conf"
	"github.com/task-go/week04/internal/data/ent"
	"log"
)

type Data struct {
	db *ent.Client
}

func NewData(ctx context.Context, conf *conf.Conf) (*Data, error) {
	client, err := ent.Open(conf.Database.Driver, conf.Database.Source)
	if err != nil {
		return nil, err
	}

	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return &Data{db: client}, err
}

var ProviderSet = wire.NewSet(NewData, NewUserRepo)
