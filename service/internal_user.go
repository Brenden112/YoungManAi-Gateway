package service

import (
	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
)

// GroupInternal is the special group name that grants access to experimental_proxy channels.
const GroupInternal = "internal"

// IsInternalUser returns true when the caller has internal-tier access:
//   - their user group is "internal", OR
//   - they are an admin (role >= RoleAdminUser).
//
// Internal users may see and call experimental_proxy channels that are explicitly enabled.
func IsInternalUser(c *gin.Context) bool {
	userGroup := common.GetContextKeyString(c, constant.ContextKeyUserGroup)
	if userGroup == GroupInternal {
		return true
	}
	userId := c.GetInt("id")
	if userId > 0 && model.IsAdmin(userId) {
		return true
	}
	return false
}
