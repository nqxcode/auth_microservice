package auth

import (
	"github.com/nqxcode/auth_microservice/internal/converter"
	"github.com/nqxcode/auth_microservice/internal/model"
	"github.com/nqxcode/platform_common/helper/gob"
)

func MakeAuditCreatePayload(user *model.User) any {
	if user == nil {
		return nil
	}

	var u *model.User
	gob.DeepClone(user, &u)
	u.Password = HiddenPassword

	return converter.ToLogUserMessageFromService(u)
}

func MakeAuditUpdatePayload(userID int64, info *model.UpdateUserInfo) any {
	return converter.ToLogUpdateUserMessageFromService(userID, info)
}
