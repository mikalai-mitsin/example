package postgres

import (
	"context"
	"fmt"
	"reflect"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
)

type objectPermissionChecker func(model any, user entities.User) error

var hasObjectPermission = map[entities.PermissionID][]objectPermissionChecker{
	entities.PermissionIDUserCreate: {objectAnybody},
	entities.PermissionIDUserList:   {objectAnybody},
	entities.PermissionIDUserDetail: {objectOwner},
	entities.PermissionIDUserUpdate: {objectOwner},
	entities.PermissionIDUserDelete: {
		objectOwner,
	}, entities.PermissionIDPostList: {objectAnybody}, entities.PermissionIDPostDetail: {objectAnybody}, entities.PermissionIDPostCreate: {objectAnybody}, entities.PermissionIDPostUpdate: {objectAnybody}, entities.PermissionIDPostDelete: {objectAnybody}, entities.PermissionIDCommentList: {objectAnybody}, entities.PermissionIDCommentDetail: {objectAnybody}, entities.PermissionIDCommentCreate: {objectAnybody}, entities.PermissionIDCommentUpdate: {objectAnybody}, entities.PermissionIDCommentDelete: {objectAnybody},
}

type PermissionRepository struct {
	database *sqlx.DB
}

func NewPermissionRepository(database *sqlx.DB) *PermissionRepository {
	return &PermissionRepository{
		database: database,
	}
}

func (r *PermissionRepository) HasPermission(
	ctx context.Context,
	permissionID entities.PermissionID,
	user entities.User,
) error {
	permission := entities.Permission{}
	q := sq.Select("permissions.id", "permissions.name").
		From("public.permissions").
		InnerJoin("group_permissions ON permissions.id = group_permissions.permission_id").
		Where(sq.Eq{"group_permissions.group_id": user.GroupID, "permissions.id": permissionID})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, &permission, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		if e.Code == errs.ErrorCodeNotFound {
			e = errs.NewPermissionDeniedError()
		}
		e.AddParam("user_id", fmt.Sprint(user.ID))
		e.AddParam("permission_id", fmt.Sprint(permissionID))
		return e
	}
	return nil
}

func (r *PermissionRepository) HasObjectPermission(
	_ context.Context,
	permission entities.PermissionID,
	user entities.User,
	obj any,
) error {
	checkers := hasObjectPermission[permission]
	for _, checker := range checkers {
		if err := checker(obj, user); err == nil {
			return nil
		}
	}
	return errs.NewPermissionDeniedError()
}

func objectOwner(model any, user entities.User) error {
	valueOf := reflect.ValueOf(model)
	if valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	if valueOf.Kind() != reflect.Struct {
		return errs.NewPermissionDeniedError()
	}
	modelUserID := valueOf.FieldByName("UserID")
	if modelUserID.Kind() == reflect.Pointer {
		modelUserID = modelUserID.Elem()
	}
	modelID := valueOf.FieldByName("ID")
	if modelID.Kind() == reflect.Pointer {
		modelID = modelID.Elem()
	}
	if modelID.String() == string(user.ID) || modelUserID.String() == string(user.ID) {
		return nil
	}
	return errs.NewPermissionDeniedError()
}

func objectOwnerOrAll(model any, user entities.User) error {
	if model == nil {
		return nil
	}
	valueOf := reflect.ValueOf(model)
	if valueOf.Kind() == reflect.Pointer {
		valueOf = valueOf.Elem()
	}
	if valueOf.Kind() != reflect.Struct {
		return errs.NewPermissionDeniedError()
	}
	modelUserID := valueOf.FieldByName("UserID")
	if modelUserID.Kind() == reflect.Pointer {
		if modelUserID.IsNil() {
			return nil
		}
		modelUserID = modelUserID.Elem()
	}
	modelID := valueOf.FieldByName("ID")
	if modelID.Kind() == reflect.Pointer {
		modelID = modelID.Elem()
	}
	if modelID.String() == string(user.ID) || modelUserID.String() == string(user.ID) {
		return nil
	}
	return errs.NewPermissionDeniedError()
}

func objectNobody(_ any, _ entities.User) error {
	return errs.NewPermissionDeniedError()
}

func objectAnybody(_ any, _ entities.User) error {
	return nil
}
