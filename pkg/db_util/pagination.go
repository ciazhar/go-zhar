package db_util

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func PageToLimitOffset(size, page int) (limit int, offset int, err error) {
	if size <= 0 {
		return 0, 0, errors.New("page size must be greater than zero")
	}
	if page < 1 {
		return 0, 0, errors.New("current page must be at least 1")
	}

	limit = size
	offset = (page - 1) * size

	return limit, offset, nil
}

func CountPageSize(dataLength int, pageSize int) int {
	d := float64(dataLength) / float64(pageSize)
	return int(math.Ceil(d))
}

func ParseCursor(cursor string) (nextPrev, id string, currPage int, err error) {

	if cursor == "" {
		return "", "", 1, nil
	}

	cursors := strings.Split(cursor, ",")
	if len(cursors) != 3 {
		return "", "", 0, errors.New("cursor must be in next,cursor,page format")
	}

	page, err := strconv.Atoi(cursors[2])
	if err != nil {
		return "", "", 0, err
	}

	return cursors[0], cursors[1], page, nil
}
