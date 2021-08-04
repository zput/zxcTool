package ztOrm

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"testing"
)

func TestParsePgError(t *testing.T) {
	pgErr := &pq.Error{Code: "23503", Message: "TestErr"}

	err := fmt.Errorf("Test Err:[%w]. ", pgErr)
	t.Log(ParsePgError(err))
	err = ParsePgError(fmt.Errorf("SeriesAll.GetById err: %w", DaoErrNotExist))
	t.Log(CheckDaoErrNew(err))

}

func TestNewCustomizeErrr(t *testing.T) {
	err := NewCustomizeErr(1000, "test")
	var err1 DaoCustomizeErrr
	t.Log(errors.Is(err1, err))

	t.Log(errors.As(err, &DaoErrNotId))
}
