package conv

import (
	dclickhouse "github.com/0LuigiCode0/hezzl/internal/domain/clickhouse"
	dpostgres "github.com/0LuigiCode0/hezzl/internal/domain/postgres"
	dusecase "github.com/0LuigiCode0/hezzl/internal/domain/usecase"
	"github.com/0LuigiCode0/hezzl/internal/utils"
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
		EventTime:   utils.TimeNow(),
	}
}

func GoodPgToResp(in *dpostgres.Good) *dusecase.GoodResp {
	return &dusecase.GoodResp{
		Id:          in.Id,
		ProjectId:   in.ProjectId,
		Name:        in.Name,
		Description: in.Description,
		Priority:    in.Priority,
		Removed:     in.Removed,
		CreatedAt:   in.CreatedAt,
	}
}

func GoodPgToRemoveResp(in *dpostgres.Good) *dusecase.RemoveGoodResp {
	return &dusecase.RemoveGoodResp{
		Id:        in.Id,
		ProjectId: in.ProjectId,
		Removed:   in.Removed,
	}
}

func GoodsPgToRespMeta(meta *dpostgres.Meta, goods []*dpostgres.Good, limit, offset int) *dusecase.GetGoodsResp {
	return &dusecase.GetGoodsResp{
		Meta: &dusecase.Meta{
			Total:   meta.Total,
			Removed: meta.Removed,
			Limit:   limit,
			Offset:  offset,
		},
		Goods: convArr(goods, GoodPgToResp),
	}
}
