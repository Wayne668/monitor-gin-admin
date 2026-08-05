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
	wire.Struct(new(dal.HostRule), "*"),
	wire.Struct(new(biz.HostRule), "*"),
	wire.Struct(new(api.HostRule), "*"),
	wire.Struct(new(dal.HostField), "*"),
	wire.Struct(new(biz.HostField), "*"),
	wire.Struct(new(api.HostField), "*"),
)
