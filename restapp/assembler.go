package restapp

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"strconv"
	"time"
)

type RestAssembler struct {
}

type Options struct {
	allowNull *bool
}

func NewOptions() *Options {
	return &Options{
		allowNull: nil,
	}
}

func (a *RestAssembler) NewOptions() *Options {
	return &Options{
		allowNull: nil,
	}
}

func (a *RestAssembler) AssFindByIdRequest(ictx iris.Context) (*ddd_query.FindByIdQuery, error) {
	tenantId, err := a.GetTenantId(ictx)
	if err != nil {
		return nil, err
	}
	id, err := a.GetId(ictx)
	if err != nil {
		return nil, err
	}
	return &ddd_query.FindByIdQuery{
		TenantId: tenantId,
		Id:       id,
	}, nil
}

func (a *RestAssembler) AssFindByIdsRequest(ictx iris.Context) (*ddd_query.FindByIdsQuery, error) {
	tenantId, err := a.GetTenantId(ictx)
	if err != nil {
		return nil, err
	}
	ids, err := a.GetIds(ictx)
	if err != nil {
		return nil, err
	}
	return &ddd_query.FindByIdsQuery{
		TenantId: tenantId,
		Ids:      ids,
	}, nil
}

func (a *RestAssembler) AssFindAllRequest(ictx iris.Context) (*ddd_query.FindAllQuery, error) {
	tenantId, err := a.GetTenantId(ictx)
	if err != nil {
		return nil, err
	}
	return &ddd_query.FindAllQuery{
		TenantId: tenantId,
	}, nil
}

func (a *RestAssembler) AssFindAutoCompleteRequest(ictx iris.Context) (*ddd_query.FindAutocompleteQuery, error) {
	fqry, err := a.AssFindPagingRequest(ictx)
	if err != nil {
		return nil, err
	}
	caseId := ictx.URLParamDefault("caseId", "")
	field := ictx.URLParamDefault("field", "")
	value := ictx.URLParamDefault("value", "")
	qry := &ddd_query.FindAutocompleteQuery{
		FindPagingQuery: fqry,
		CaseId:          caseId,
		Field:           field,
		Value:           value,
	}
	return qry, err
}

func (a *RestAssembler) AssFindPagingRequest(ictx iris.Context) (*ddd_query.FindPagingQuery, error) {
	tenantId, err := a.GetTenantId(ictx)
	if err != nil {
		return nil, err
	}
	pageNum := ictx.URLParamInt64Default("page-num", 0)
	pageSize := ictx.URLParamInt64Default("page-size", 20)
	filter := ictx.URLParamDefault("filter", "")
	sort := ictx.URLParamDefault("sort", "")
	fields := ictx.URLParamDefault("fields", "")
	isTotalRows := true
	if val := ictx.URLParamDefault("is-total-rows", "true"); val == "false" {
		isTotalRows = false
	}
	groupCols := ictx.URLParamDefault("group-cols", "")
	groupKeys := ictx.URLParamDefault("group-keys", "")
	valueCols := ictx.URLParamDefault("value-cols", "")
	mustFilter := ictx.URLParamDefault("must-filter", "")

	req := ddd_repository.FindPagingQueryDTO{
		TenantId:    tenantId,
		PageNum:     pageNum,
		PageSize:    pageSize,
		Filter:      filter,
		MustFilter:  mustFilter,
		Sort:        sort,
		Fields:      fields,
		IsTotalRows: isTotalRows,
		GroupCols:   groupCols,
		GroupKeys:   groupKeys,
		ValueCols:   valueCols,
	}

	return req.NewFindPagingQueryRequest(), nil
}

func (a *RestAssembler) GetTenantId(ictx iris.Context) (string, error) {
	return a.GetIdParam(ictx, "tenantId")
}

func (a *RestAssembler) GetCaseId(ictx iris.Context) (string, error) {
	return a.GetIdParam(ictx, "caseId")
}

func (a *RestAssembler) GetId(ictx iris.Context) (string, error) {
	return a.GetIdParam(ictx, "id")
}

func (a *RestAssembler) GetIds(ictx iris.Context) ([]string, error) {
	ids := ictx.URLParamSlice("id")
	return ids, nil
}

func (a *RestAssembler) GetIdParam(ictx iris.Context, name string) (string, error) {
	id := ictx.Params().GetStringDefault(name, "")
	if id == "" {
		return "", errors.New(name + " is empty")
	}
	return id, nil
}

//
// GetValueByUrlPath
// @Description: 从URL路径中中获取string变量
// @receiver a
// @param ictx iris上下文
// @param name 变量名
// @param opts 可选项
// @return string 返回值
// @return error 错误
//
func (a *RestAssembler) GetValueByUrlPath(ictx iris.Context, name string, opts ...*Options) (string, error) {
	options := NewOptions().Merge(opts)
	id := ictx.Params().GetStringDefault(name, "")
	if id == "" && !options.GetAllowNull() {
		return "", errors.New(name + " cannot be empty")
	}
	return id, nil
}

func (a *RestAssembler) URLParam(ictx iris.Context, name string, opts ...*Options) (string, error) {
	options := NewOptions().Merge(opts)
	id := ictx.URLParamDefault(name, "")
	if id == "" && !options.GetAllowNull() {
		return "", errors.New(name + " cannot be empty")
	}
	return id, nil
}

func (a *RestAssembler) URLParamBool(ictx iris.Context, name string, opts ...*Options) (bool, error) {
	options := NewOptions().Merge(opts)
	val := ictx.URLParamDefault(name, "false")
	if val == "" && !options.GetAllowNull() {
		return false, errors.New(name + " cannot be empty")
	}
	return strconv.ParseBool(val)
}

func (a *RestAssembler) URLParamTime(ictx iris.Context, name string, opts ...*Options) (*time.Time, error) {
	options := NewOptions().Merge(opts)
	id := ictx.URLParamDefault(name, "")
	if id == "" && !options.GetAllowNull() {
		return nil, errors.New(name + " cannot be empty")
	}
	if id == "" && options.GetAllowNull() {
		return nil, nil
	}
	v, err := time.Parse("2006-01-02", id)
	return &v, err
}

func (a *RestAssembler) URLParamFloat(ictx iris.Context, name string, opts ...*Options) (*float64, error) {
	options := NewOptions().Merge(opts)
	str := ictx.URLParamDefault(name, "")
	if str == "" && !options.GetAllowNull() {
		return nil, errors.New(name + " cannot be empty")
	}
	if str == "" && options.GetAllowNull() {
		return nil, nil
	}
	// 将字符串转换为float64类型
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil, err
	}
	return &f, err
}

func (o *Options) Merge(opts []*Options) *Options {
	for _, i := range opts {
		if i.allowNull != nil {
			o.allowNull = i.allowNull
		}
	}
	return o
}

func (o *Options) SetAllowNull(v bool) *Options {
	o.allowNull = &v
	return o
}

func (o *Options) GetAllowNull() bool {
	if o.allowNull != nil {
		return *o.allowNull
	}
	return false
}
