package conv

import (
	"time"

	dclickhouse "github.com/0LuigiCode0/hezzl/internal/domain/clickhouse"
	dhttp "github.com/0LuigiCode0/hezzl/internal/domain/http"
	dpostgres "github.com/0LuigiCode0/hezzl/internal/domain/postgres"
)

func convArr[In, Out any](in []In, f func(In) Out) []Out {
	out := make([]Out, 0, len(in))
	for _, v := range in {
		out = append(out, f(v))
	}

	return out
}

func GoodPgToCh(in *dpostgres.Good) *dclickhouse.LogEventGood {
	return &dclickhouse.LogEventGood{
		Id:          in.Id,
		ProjectId:   in.ProjectId,
		Name:        in.Name,
		Description: in.Description,
		Priority:    in.Priority,
		Removed:     in.Removed,
		EventTime:   time.Now(),
	}
}

func GoodPgToResp(in *dpostgres.Good) *dhttp.GoodResp {
	return &dhttp.GoodResp{
		Id:          in.Id,
		ProjectId:   in.ProjectId,
		Name:        in.Name,
		Description: in.Description,
		Priority:    in.Priority,
		Removed:     in.Removed,
		CreatedAt:   in.CreatedAt,
	}
}

func GoodPgToRemoveResp(in *dpostgres.Good) *dhttp.RemoveGoodResp {
	return &dhttp.RemoveGoodResp{
		Id:        in.Id,
		ProjectId: in.ProjectId,
		Removed:   in.Removed,
	}
}

func GoodsPgToRespMeta(meta *dpostgres.Meta, goods []*dpostgres.Good, limit, offset int) *dhttp.GetGoodsResp {
	return &dhttp.GetGoodsResp{
		Meta: dhttp.Meta{
			Total:   meta.Total,
			Removed: meta.Removed,
			Limit:   limit,
			Offset:  offset,
		},
		Goods: convArr(goods, GoodPgToResp),
	}
}
