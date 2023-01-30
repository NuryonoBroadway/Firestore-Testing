package collectionxclient

import (
	collectionxservice "firebaseapi/collectionx/collectionx_service"
	"reflect"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func payloadBuilder(p *Payload) ([]*collectionxservice.PathProto, *collectionxservice.FilteringProto, *collectionxservice.PaginationProto) {
	pathProto := make([]*collectionxservice.PathProto, len(p.Path))
	for i := 0; i < len(p.Path); i++ {
		pathProto[i] = &collectionxservice.PathProto{
			CollectionId: p.Path[i].CollectionID,
			DocumentId:   p.Path[i].DocumentID,
			NewDocument:  p.Path[i].NewDocument,
		}
	}

	filters := make([]*collectionxservice.FilterProto, len(p.query.Filter))
	if len(filters) != 0 {
		for i := 0; i < len(p.query.Filter); i++ {
			filters[i] = &collectionxservice.FilterProto{
				By: p.query.Filter[i].By,
				Op: p.query.Filter[i].Op,
			}

			xtyp := reflect.TypeOf(p.query.Filter[i].Val)
			switch xtyp.Kind() {
			case reflect.Bool:
				filters[i].Val = &collectionxservice.FilterProto_ValBool{
					ValBool: p.query.Filter[i].Val.(bool),
				}
			case reflect.String:
				filters[i].Val = &collectionxservice.FilterProto_ValString{
					ValString: p.query.Filter[i].Val.(string),
				}
			case reflect.Int64:
				filters[i].Val = &collectionxservice.FilterProto_ValInt{
					ValInt: p.query.Filter[i].Val.(int64),
				}

			}

		}
	}

	sorts := make([]*collectionxservice.SortProto, len(p.query.Sort))
	if len(sorts) == 0 {
		sorts = append(sorts, &collectionxservice.SortProto{
			OrderBy:   "created_at",
			OrderType: collectionxservice.OrderTypeProto_ORDER_TYPE_ASC,
		})
	} else {
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

	query := &collectionxservice.FilteringProto{
		Sort:      sorts,
		Filter:    filters,
		DateRange: ranges,
	}

	pagination := new(collectionxservice.PaginationProto)
	if p.isPagination {
		if p.limit == 0 {
			p.limit = 2 // default limit
		}
		pagination.Page = p.pagination.Page
	}

	return pathProto, query, pagination
}
