package converter

import (
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	"github.com/nqxcode/platform_common/pagination"
)

// ToLimitFromDesc convert to pagination limit
func ToLimitFromDesc(message *desc.Limit) *pagination.Limit {
	return &pagination.Limit{
		Offset: message.Offset,
		Limit:  message.Limit,
	}
}
