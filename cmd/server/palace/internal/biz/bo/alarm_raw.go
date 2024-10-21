package bo

type (

	// CreateAlarmRawParams 创建告警原始数据参数
	CreateAlarmRawParams struct {
		Fingerprint string `json:"fingerprint"`
		RawInfo     string `json:"rawInfo"`
	}

	// GetTeamStrategyParams 获取团队策略参数
	GetTeamStrategyParams struct {
		TeamID     uint32 `json:"teamId"`
		StrategyID uint32 `json:"strategyId"`
	}

	// GetTeamStrategyLevelParams 获取团队策略等级参数
	GetTeamStrategyLevelParams struct {
		TeamID  uint32 `json:"teamId"`
		LevelID uint32 `json:"level"`
	}

	//GetTeamDatasourceParams  获取团队数据源信息参数
	GetTeamDatasourceParams struct {
		TeamID        uint32   `json:"teamId"`
		DatasourceIds []uint32 `json:"datasourceIds"`
	}
)