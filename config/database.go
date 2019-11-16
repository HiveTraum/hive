package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InitPool() *pgxpool.Pool {

	env := InitEnv()

	databaseName := env.DatabaseName

	if isTest() {
		databaseName = databaseName + "_test"
	}

	configString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?pool_max_conns=5",
		env.DatabaseUser,
		env.DatabasePass,
		env.DatabaseHost,
		env.DatabasePort,
		databaseName)

	config, err := pgxpool.ParseConfig(configString)

	//config.ConnConfig.Logger = log15adapter.NewLogger(log15.New("module", "pgx"))

	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		panic(err)
	}

	return pool
}
