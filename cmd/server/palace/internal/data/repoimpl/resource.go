package repoimpl

import (
	"context"

	"github.com/aide-family/moon/pkg/merr"
	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/aide-family/moon/cmd/server/palace/internal/biz/bo"
	"github.com/aide-family/moon/cmd/server/palace/internal/biz/repository"
	"github.com/aide-family/moon/cmd/server/palace/internal/data"
	"github.com/aide-family/moon/pkg/palace/model"
	"github.com/aide-family/moon/pkg/palace/model/query"
	"github.com/aide-family/moon/pkg/util/types"
	"github.com/aide-family/moon/pkg/vobj"
)

// NewResourceRepository 创建资源仓库
func NewResourceRepository(data *data.Data) repository.Resource {
	return &resourceRepositoryImpl{
		data: data,
	}
}

type resourceRepositoryImpl struct {
	data *data.Data
}

func (l *resourceRepositoryImpl) CheckPath(ctx context.Context, s string) error {
	// 1. 检查主库API是否存在且开启
	mainQuery := query.Use(l.data.GetMainDB(ctx))
	mainApiDo, err := mainQuery.SysAPI.WithContext(ctx).Where(mainQuery.SysAPI.Path.Eq(s)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merr.ErrorI18nForbidden(ctx)
		}
		return err
	}
	if !mainApiDo.Status.IsEnable() {
		return merr.ErrorI18nForbidden(ctx)
	}

	// 2. 检查业务库API是否存在且开启
	bizQuery, err := getBizQuery(ctx, l.data)
	if !types.IsNil(err) {
		return err
	}

	bizApiDo, err := bizQuery.SysTeamAPI.WithContext(ctx).Where(bizQuery.SysTeamAPI.Path.Eq(s)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merr.ErrorI18nForbidden(ctx)
		}
		return err
	}
	if !bizApiDo.Status.IsEnable() {
		return merr.ErrorI18nForbidden(ctx)
	}
	return nil
}

func (l *resourceRepositoryImpl) GetByID(ctx context.Context, id uint32) (*model.SysAPI, error) {
	mainQuery := query.Use(l.data.GetMainDB(ctx))
	return mainQuery.SysAPI.WithContext(ctx).Where(mainQuery.SysAPI.ID.Eq(id)).First()
}

func (l *resourceRepositoryImpl) FindByPage(ctx context.Context, params *bo.QueryResourceListParams) ([]*model.SysAPI, error) {
	mainQuery := query.Use(l.data.GetMainDB(ctx))
	apiQuery := mainQuery.SysAPI.WithContext(ctx)

	var wheres []gen.Condition

	if !types.TextIsNull(params.Keyword) {
		wheres = append(wheres, apiQuery.Or(mainQuery.SysAPI.Name.Like(params.Keyword), mainQuery.SysAPI.Path.Like(params.Keyword)))
	}
	apiQuery = apiQuery.Where(wheres...)
	if !params.Status.IsUnknown() {
		apiQuery = apiQuery.Where(mainQuery.SysAPI.Status.Eq(params.Status.GetValue()))
	}
	if !params.IsAll {
		if err := types.WithPageQuery[query.ISysAPIDo](apiQuery, params.Page); err != nil {
			return nil, err
		}
	}

	return apiQuery.Order(mainQuery.SysAPI.ID.Desc()).Find()
}

func (l *resourceRepositoryImpl) UpdateStatus(ctx context.Context, status vobj.Status, ids ...uint32) error {
	mainQuery := query.Use(l.data.GetMainDB(ctx))
	_, err := mainQuery.SysAPI.WithContext(ctx).Where(mainQuery.SysAPI.ID.In(ids...)).Update(mainQuery.SysAPI.Status, status)
	return err
}

func (l *resourceRepositoryImpl) FindSelectByPage(ctx context.Context, params *bo.QueryResourceListParams) ([]*model.SysAPI, error) {
	mainQuery := query.Use(l.data.GetMainDB(ctx))
	apiQuery := mainQuery.SysAPI.WithContext(ctx)

	if !types.TextIsNull(params.Keyword) {
		apiQuery = apiQuery.Or(mainQuery.SysAPI.Name.Like(params.Keyword), mainQuery.SysAPI.Path.Like(params.Keyword))
	}
	if !params.Status.IsUnknown() {
		apiQuery = apiQuery.Where(mainQuery.SysAPI.Status.Eq(params.Status.GetValue()))
	}
	if !params.IsAll {
		if err := types.WithPageQuery[query.ISysAPIDo](apiQuery, params.Page); err != nil {
			return nil, err
		}
	}

	return apiQuery.Select(mainQuery.SysAPI.ID, query.SysAPI.Name, mainQuery.SysAPI.Status, mainQuery.SysAPI.DeletedAt).Order(mainQuery.SysAPI.ID.Desc()).Find()
}
