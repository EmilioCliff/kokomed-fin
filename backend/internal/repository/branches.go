package repository

import "context"

type Branch struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type BranchRepository interface {
	CreateBranch(ctx context.Context, branch *Branch) (Branch, error)
	ListBranches(ctx context.Context) ([]Branch, error)
	GetBranchByID(ctx context.Context, id uint32) (Branch, error)
	UpdateBranch(ctx context.Context, name string, id uint32) (Branch, error)
}
