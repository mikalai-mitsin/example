package errs

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

const sqlConflictCode = "23505"

func FromPostgresError(err error) *Error {
	e := &Error{Code: ErrorCodeInternal, Message: "Unexpected behavior.", Params: nil, Err: err}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		e.AddParam("details", pqErr.Detail)
		e.AddParam("message", pqErr.Message)
		e.AddParam("postgres_code", fmt.Sprint(pqErr.Code))
		if pqErr.Code == sqlConflictCode {
			e = NewInvalidFormError().WithCause(err)
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		e = NewEntityNotFoundError()
	}
	return e
}
