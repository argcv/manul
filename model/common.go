package model

import (
	"github.com/argcv/go-argcvapis/status/errcodes"
	"github.com/argcv/go-argcvapis/status/status"
	"github.com/golang/protobuf/ptypes/any"
)

/* Status is a wrapper to pb-status
 */
type Status struct {
	Code    errcodes.Code
	Message string
	Details []*any.Any
}

func (st *Status) ToPbStatus() (pbSt *status.Status) {
	pbSt = &status.Status{
		Code:    st.Code,
		Message: st.Message,
		Details: st.Details,
	}
	return
}

func FromPbStatus(st *status.Status) *Status {
	return &Status{
		Code:    st.Code,
		Message: st.Message,
		Details: st.Details,
	}
}

/* OrderType is an option type
 */
type OrderType int

const (
	OrderTypeUnset OrderType = iota
	OrderTypeAsc
	OrderTypeDesc
)
