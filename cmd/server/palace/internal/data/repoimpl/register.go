package repoimpl

import (
	"github.com/google/wire"
)

var ProviderSetRepoImpl = wire.NewSet(
	NewUserRepository,
	NewCaptchaRepository,
	NewTeamRepository,
	NewCacheRepository,
	NewResourceRepository,
	NewTeamRoleRepository,
	NewTeamMenuRepository,
	NewDatasourceRepository,
	NewDatasourceMetricRepository,
	NewLockRepository,
)