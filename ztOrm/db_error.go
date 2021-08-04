package ztOrm

import (
	"errors"
	"github.com/lib/pq"
)

var (
	DaoErrNotExist             = errors.New("not found")
	DaoErrFailedToInsert       = errors.New("Failed to insert. ")
	DaoErrForeignKeyConstraint = errors.New("foreign_key_violation. ") //外键约束
	DaoErrUniqueViolation      = errors.New("unique_violation. ")      //唯一键约束
	DaoErrNotNullViolation     = errors.New("not_null_violation. ")    //not null约束
	DaoErrNotId                = errors.New("Id is 0. ")
	DaoParamErr                = errors.New("parameter exception. ")
)

type DaoCustomizeErrr struct {
	ErrMsg  string
	ErrCode int
}

func (e DaoCustomizeErrr) Error() string {
	return e.ErrMsg
}

var (
	ErrImportDifferentClassificationCourse = errors.New("Different classes cannot be imported into the same series. ")
)

func NewCustomizeErr(errCode int, errMsg string) DaoCustomizeErrr {
	return DaoCustomizeErrr{ErrCode: errCode, ErrMsg: errMsg}
}

func CheckDaoErr(err error) (errNo int, errMsg string, hasErr bool) {
	if err != nil {
		hasErr = true
		switch err {
		case DaoErrNotExist:
			errMsg = errno.EM_ResourceNotFound
			errNo = errno.EN_ResourceNotFound
		case DaoErrForeignKeyConstraint:
			errNo = errno.EN_InvalidArgs
			errMsg = err.Error()
		default:
			errMsg = errno.EM_InnerServer
			errNo = errno.EN_InnerServer
		}
	}
	return
}

func ParsePgError(err error) error {
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if !ok {
			unwrapErr := errors.Unwrap(err)
			if unwrapErr != nil {
				return ParsePgError(unwrapErr)
			}
			return err
		}
		switch pqErr.Code {
		case "23503":
			//违反外键约束
			return DaoErrForeignKeyConstraint
		case "23505":
			//唯一键约束
			return DaoErrUniqueViolation
		case "23502":
			//not null违规
			return DaoErrNotNullViolation
		default:
			return err
		}
	}
	return nil
}

func CheckDaoErrNew(err error) (errNo int, errMsg string) {
	if err != nil {
		switch err {
		case DaoErrNotExist:
			errMsg = errno.EM_ResourceNotFound
			errNo = errno.EN_ResourceNotFound
		case DaoErrForeignKeyConstraint, DaoErrNotNullViolation:
			errNo = errno.EN_InvalidArgs
			errMsg = err.Error()
		case DaoErrUniqueViolation:
			errNo = errno.EN_ResourceExists
			errMsg = errno.EM_ResourceExists
		case ErrImportDifferentClassificationCourse:
			errNo = errno.EN_InvalidArgs
			errMsg = err.Error()
		case DaoParamErr:
			errNo = errno.EN_InvalidArgs
			errMsg = errno.EM_InvalidOperation
		default:
			switch err.(type) {
			case DaoCustomizeErrr:
				errNo = err.(DaoCustomizeErrr).ErrCode
				errMsg = err.Error()
			default:
				errMsg = errno.EM_InnerServer
				errNo = errno.EN_InnerServer
			}
		}
	}
	return
}
