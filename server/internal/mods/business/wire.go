package business

import (
	"monitor-gin-admin/internal/mods/business/api"
	"monitor-gin-admin/internal/mods/business/biz"
	"monitor-gin-admin/internal/mods/business/dal"
	"github.com/google/wire"
)

// Collection of wire providers
var Set = wire.NewSet(
	wire.Struct(new(Business), "*"),
	wire.Struct(new(dal.Category), "*"),
	wire.Struct(new(biz.Category), "*"),
	wire.Struct(new(api.Category), "*"),
)
