package stuntdb

import (
	"embed"

	"github.com/boringsql/queries"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

var sqlStore = loadQueries()

func loadQueries() *queries.QueryStore {
	store := queries.NewQueryStore()
	if err := store.LoadFromEmbed(sqlFiles, "sql"); err != nil {
		panic(err)
	}
	return store
}

func query(name string) string {
	return sqlStore.MustHaveQuery(name).RawQuery()
}
