package user

import (
	"strconv"

	"github.com/nqxcode/platform_common/pagination"
)

func buildListCacheKeyByLimit(limit pagination.Limit) string {
	return buildListCacheKey(strconv.Itoa(int(limit.Offset)) + "-" + strconv.Itoa(int(limit.Limit)))
}

func buildListCacheKey(value string) string {
	return listCacheKey + ":" + value
}
