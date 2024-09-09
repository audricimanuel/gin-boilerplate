package paramquery

import (
	"context"
	"gin-boilerplate/src/middleware"
)

type (
	BaseParamQuery struct {
		Page    int
		Limit   int
		Offset  int
		Keyword *string
	}
)

func SetBaseParamQuery(ctx context.Context) BaseParamQuery {
	paramQuery := BaseParamQuery{
		Page:   ctx.Value(middleware.ParamQueryPage).(int),
		Limit:  ctx.Value(middleware.ParamQueryLimit).(int),
		Offset: ctx.Value(middleware.ParamQueryOffset).(int),
	}

	keyword := ctx.Value(middleware.ParamQueryKeyword).(string)
	if keyword != "" {
		paramQuery.Keyword = &keyword
	}

	return paramQuery
}
