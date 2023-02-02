package collectionxclient

import (
	"encoding/json"
	collectionxservice "firebaseapi/collectionx/collectionx_service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func payloadBuilder(p *Payload) (paths []*collectionxservice.PathProto, query *collectionxservice.FilteringProto, err error) {
	paths = make([]*collectionxservice.PathProto, len(p.Path))
	for i := 0; i < len(p.Path); i++ {
		paths[i] = &collectionxservice.PathProto{
			CollectionId:    p.Path[i].CollectionID,
			DocumentId:      p.Path[i].DocumentID,
			NewDocument:     p.Path[i].NewDocument,
			CollectionGroup: p.Path[i].CollectionGroup,
		}
	}

	filters := make([]*collectionxservice.FilterProto, len(p.query.Filter))
	if len(p.query.Filter) != 0 {
		for i := 0; i < len(p.query.Filter); i++ {
			data, err := json.Marshal(p.query.Filter[i].Val)
			if err != nil {
				return nil, nil, err
			}

			filters[i] = &collectionxservice.FilterProto{
				By:  p.query.Filter[i].By,
				Op:  p.query.Filter[i].Op,
				Val: data,
			}
		}
	}

	sorts := make([]*collectionxservice.SortProto, len(p.query.Sort))
	if len(p.query.Sort) != 0 {
		for i := 0; i < len(p.query.Sort); i++ {
			sorts[i] = &collectionxservice.SortProto{
				OrderBy:   p.query.Sort[i].OrderBy,
				OrderType: collectionxservice.OrderTypeProto(p.query.Sort[i].OrderType),
			}
		}
	}

	ranges := new(collectionxservice.DateRangeProto)
	if p.query.DateRange.Field != "" {
		ranges.Field = p.query.DateRange.Field
		ranges.Start = timestamppb.New(p.query.DateRange.Start)
		ranges.End = timestamppb.New(p.query.DateRange.End)
	}

	query = &collectionxservice.FilteringProto{
		Sort:      sorts,
		Filter:    filters,
		DateRange: ranges,
	}

	return paths, query, nil
}
