package dhttp

import (
	"errors"
	"time"
)

type GoodResp struct {
	Id          int       `json:"id"`
	ProjectId   int       `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
}

type (
	Meta struct {
		Total   int `json:"total"`
		Removed int `json:"removed"`
		Limit   int `json:"limit"`
		Offset  int `json:"offset"`
	}
	GetGoodsResp struct {
		Meta  Meta        `json:"meta"`
		Goods []*GoodResp `json:"goods"`
	}
)

type RemoveGoodResp struct {
	Id        int  `json:"id"`
	ProjectId int  `json:"project_id"`
	Removed   bool `json:"removed"`
}

type CreateGoodReq struct {
	Name string `json:"name"`
}

func (cgr *CreateGoodReq) Validate() error {
	return validateString(&cgr.Name, "name")
}

type UpdateGoodReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (ugr *UpdateGoodReq) Validate() error {
	return errors.Join(validateString(&ugr.Name, "name"), validateString(&ugr.Description, "description"))
}
