package errs

import (
	"errors"

	"github.com/IBM/sarama"
)

func FromKafkaError(err error) *Error {
	e := &Error{Code: ErrorCodeInternal, Message: "Unexpected behavior.", Params: nil, Err: err}
	var prErr *sarama.ProducerError
	if errors.As(err, &prErr) {
		e.AddParam("error", prErr.Error())
	}
	e.AddParam("error", err.Error())
	return e
}
