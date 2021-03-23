package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rancher/wrangler/pkg/slice"

	"github.com/rancher/wrangler/pkg/schemas/validation"
)

type ResponseAPIError struct {
	validation.ErrorCode
	Message string
}

type APIErrorCode interface {
	ErrorCode() *validation.ErrorCode
}

type ResponseError struct {
	RespCode         int
	RespBody         []byte
	ResponseAPIError ResponseAPIError
}

func (e *ResponseError) ErrorCode() *validation.ErrorCode {
	return &validation.ErrorCode{
		Code:   e.ResponseAPIError.Code,
		Status: e.ResponseAPIError.Status,
	}
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("respCode： %d, respBody: %s", e.RespCode, string(e.RespBody))
}

func NewResponseError(respCode int, respBody []byte) error {
	var responseAPIError ResponseAPIError
	if err := json.Unmarshal(respBody, &responseAPIError); err != nil {
		return err
	}
	return &ResponseError{
		RespCode:         respCode,
		RespBody:         respBody,
		ResponseAPIError: responseAPIError,
	}
}

func CodeForError(err error) *validation.ErrorCode {
	if errorCode := APIErrorCode(nil); errors.As(err, &errorCode) {
		return errorCode.ErrorCode()
	}
	return nil
}

func IsNotFound(err error) bool {
	return CodeForError(err).Code == validation.NotFound.Code
}

func IsConflict(err error) bool {
	return CodeForError(err).Code == validation.Conflict.Code
}

func RetryOnCodes(retryNum, retryInterval int64, process func() error, codes ...string) error {
	for {
		if err := process(); err != nil {
			if retryNum == 0 {
				return err
			}

			if !slice.ContainsString(codes, CodeForError(err).Code) {
				return err
			}
			retryNum--
			if retryInterval > 0 {
				time.Sleep(time.Duration(retryInterval) * time.Second)
			}
			continue
		}
		return nil
	}
}
