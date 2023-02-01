package collectionxclient

import (
	"firebaseapi/collectionx/collection_core_service"
	"reflect"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func payloadBuilder(p *Payload) ([]*collection_core_service.PathProto, *collection_core_service.FilteringProto, *collection_core_service.PaginationProto) {
	pathProto := make([]*collection_core_service.PathProto, len(p.Path))
	for i := 0; i < len(p.Path); i++ {
		pathProto[i] = &collection_core_service.PathProto{
			CollectionId: p.Path[i].CollectionID,
			DocumentId:   p.Path[i].DocumentID,
			NewDocument:  p.Path[i].NewDocument,
		}
	}

	filters := make([]*collection_core_service.FilterProto, len(p.query.Filter))
	if len(filters) != 0 {
		for i := 0; i < len(p.query.Filter); i++ {
			filters[i] = &collection_core_service.FilterProto{
				By: p.query.Filter[i].By,
				Op: p.query.Filter[i].Op,
			}

			xtyp := reflect.TypeOf(p.query.Filter[i].Val)
			switch xtyp.Kind() {
			case reflect.Bool:
				filters[i].Val = &collection_core_service.FilterProto_ValBool{
					ValBool: p.query.Filter[i].Val.(bool),
				}
			case reflect.String:
				filters[i].Val = &collection_core_service.FilterProto_ValString{
					ValString: p.query.Filter[i].Val.(string),
				}
			case reflect.Int64:
				filters[i].Val = &collection_core_service.FilterProto_ValInt{
					ValInt: p.query.Filter[i].Val.(int64),
				}

			}

		}
	}

	sorts := make([]*collection_core_service.SortProto, len(p.query.Sort))
	if len(sorts) == 0 {
		sorts = append(sorts, &collection_core_service.SortProto{
			OrderBy:   "created_at",
			OrderType: collection_core_service.OrderTypeProto_ORDER_TYPE_ASC,
		})
	} else {
		for i := 0; i < len(p.query.Sort); i++ {
			sorts[i] = &collection_core_service.SortProto{
				OrderBy:   p.query.Sort[i].OrderBy,
				OrderType: collection_core_service.OrderTypeProto(p.query.Sort[i].OrderType),
			}
		}
	}

	ranges := new(collection_core_service.DateRangeProto)
	if p.query.DateRange.Field != "" {
		ranges.Field = p.query.DateRange.Field
		ranges.Start = timestamppb.New(p.query.DateRange.Start)
		ranges.End = timestamppb.New(p.query.DateRange.End)
	}

	query := &collection_core_service.FilteringProto{
		Sort:      sorts,
		Filter:    filters,
		DateRange: ranges,
	}

	pagination := new(collection_core_service.PaginationProto)
	if p.isPagination {
		if p.limit == 0 {
			p.limit = 2 // default limit
		}
		pagination.Page = p.pagination.Page
	}

	return pathProto, query, pagination
}
