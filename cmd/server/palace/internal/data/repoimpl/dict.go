package repoimpl

import (
	"context"

	"github.com/aide-family/moon/cmd/server/palace/internal/biz/bo"
	"github.com/aide-family/moon/cmd/server/palace/internal/biz/repository"
	"github.com/aide-family/moon/cmd/server/palace/internal/data"
	"github.com/aide-family/moon/pkg/palace/imodel"
	"github.com/aide-family/moon/pkg/palace/model"
	"github.com/aide-family/moon/pkg/palace/model/query"
	"github.com/aide-family/moon/pkg/util/types"
	"gorm.io/gen"
)

// NewDictRepository 创建数据库字典操作
func NewDictRepository(data *data.Data) repository.Dict {
	return &dictRepositoryImpl{
		data: data,
	}
}

type dictRepositoryImpl struct {
	data *data.Data
}

func (l *dictRepositoryImpl) UpdateStatusByIds(ctx context.Context, params *bo.UpdateDictStatusParams) error {
	ids := params.IDs
	mainQuery := query.Use(l.data.GetMainDB(ctx))
	_, err := mainQuery.WithContext(ctx).SysDict.Where(mainQuery.SysDict.ID.In(ids...)).Update(mainQuery.SysDict.Status, params.Status)
	return err
}

func (l *dictRepositoryImpl) DeleteByID(ctx context.Context, id uint32) error {
	mainQuery := query.Use(l.data.GetMainDB(ctx))
	_, err := mainQuery.WithContext(ctx).SysDict.Where(mainQuery.SysDict.ID.Eq(id)).Delete()
	return err
}

func (l *dictRepositoryImpl) Create(ctx context.Context, dict *bo.CreateDictParams) (imodel.IDict, error) {
	// system creation
	dictModel := createDictParamsToModel(ctx, dict)
	mainQuery := query.Use(l.data.GetMainDB(ctx))
	if err := mainQuery.WithContext(ctx).SysDict.Create(dictModel); !types.IsNil(err) {
		return nil, err
	}
	return dictModel, nil
}

func (l *dictRepositoryImpl) FindByPage(ctx context.Context, params *bo.QueryDictListParams) ([]imodel.IDict, error) {
	return l.listDictModel(ctx, params)
}

func (l *dictRepositoryImpl) GetByID(ctx context.Context, id uint32) (imodel.IDict, error) {
	mainQuery := query.Use(l.data.GetMainDB(ctx))
	return mainQuery.SysDict.WithContext(ctx).Where(mainQuery.SysDict.ID.Eq(id)).First()
}

func (l *dictRepositoryImpl) UpdateByID(ctx context.Context, dict *bo.UpdateDictParams) error {
	return l.updateDictModel(ctx, dict)
}

func (l *dictRepositoryImpl) listDictModel(ctx context.Context, params *bo.QueryDictListParams) ([]imodel.IDict, error) {
	dictQuery := query.Use(l.data.GetMainDB(ctx)).SysDict
	queryWrapper := dictQuery.WithContext(ctx)

	var wheres []gen.Condition
	if !params.Status.IsUnknown() {
		wheres = append(wheres, dictQuery.Status.Eq(params.Status.GetValue()))
	}

	if !params.DictType.IsUnknown() {
		wheres = append(wheres, dictQuery.DictType.Eq(params.DictType.GetValue()))
	}

	if !types.TextIsNull(params.Keyword) {
		queryWrapper = queryWrapper.Or(
			dictQuery.Name.Like(params.Keyword),
			dictQuery.Value.Like(params.Keyword),
			dictQuery.Remark.Like(params.Keyword),
		)
	}
	var err error
	queryWrapper = queryWrapper.Where(wheres...)
	if queryWrapper, err = types.WithPageQuery(queryWrapper, params.Page); err != nil {
		return nil, err
	}
	dbDictList, err := queryWrapper.Order(dictQuery.ID.Desc()).Find()
	if !types.IsNil(err) {
		return nil, err
	}
	dictList := types.SliceTo(dbDictList, func(dict *model.SysDict) imodel.IDict {
		return dict
	})
	return dictList, nil
}

func (l *dictRepositoryImpl) updateDictModel(ctx context.Context, params *bo.UpdateDictParams) error {
	id := params.ID
	updateParam := params.UpdateParam
	mainQuery := query.Use(l.data.GetMainDB(ctx))
	_, err := mainQuery.SysDict.WithContext(ctx).Where(mainQuery.SysDict.ID.Eq(id)).UpdateSimple(
		mainQuery.SysDict.Name.Value(updateParam.Name),
		mainQuery.SysDict.Value.Value(updateParam.Value),
		mainQuery.SysDict.CSSClass.Value(updateParam.CSSClass),
		mainQuery.SysDict.ColorType.Value(updateParam.ColorType),
		mainQuery.SysDict.Remark.Value(updateParam.Remark),
		mainQuery.SysDict.ImageURL.Value(updateParam.ImageURL),
		mainQuery.SysDict.Icon.Value(updateParam.Icon),
	)
	return err
}

// createDictParamsToModel create dict params to model
func createDictParamsToModel(ctx context.Context, dict *bo.CreateDictParams) *model.SysDict {
	if types.IsNil(dict) {
		return nil
	}
	dictModel := &model.SysDict{
		Name:         dict.Name,
		Value:        dict.Value,
		DictType:     dict.DictType,
		ColorType:    dict.ColorType,
		CSSClass:     dict.CSSClass,
		Status:       dict.Status,
		Remark:       dict.Remark,
		Icon:         dict.Icon,
		ImageURL:     dict.ImageURL,
		LanguageCode: dict.LanguageCode,
	}
	dictModel.WithContext(ctx)
	return dictModel
}
