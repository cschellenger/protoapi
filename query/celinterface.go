package query

import (
	"darkport.net/protoapi/model"
)

type QueryBuilder interface {
	BuildQuery(searchRequest *model.SearchRequest) (string, []any, error)
}
