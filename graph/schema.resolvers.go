package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"testbed-monitor/graph/generated"
	"testbed-monitor/graph/model"
)

func (r *queryResolver) Hosts(ctx context.Context) ([]*model.HostStatus, error) {
	return r.GetHosts()
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
