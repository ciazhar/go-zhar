package mongo

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

func CreatePagingAndSortingOptions(page int, size int, sort string) (*options.FindOptions, error) {

	sortFieldsBson := bson.D{}
	if sort != "" {
		sort2 := strings.TrimSpace(sort)
		sort3 := strings.Split(sort2, ";")
		for i := range sort3 {

			sort4 := strings.Split(sort3[i], ",")
			if len(sort4) != 2 {
				return nil, errors.New("sort format invalid")
			}

			if sort4[1] == "asc" {
				sortFieldsBson = append(sortFieldsBson, bson.E{Key: sort4[0], Value: 1})
			} else if sort4[1] == "desc" {
				sortFieldsBson = append(sortFieldsBson, bson.E{Key: sort4[0], Value: -1})
			}
		}
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((page - 1) * size))
	findOptions.SetLimit(int64(size))
	findOptions.SetSort(sortFieldsBson)

	return findOptions, nil
}
