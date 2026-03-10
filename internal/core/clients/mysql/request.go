package mysql

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	LimitMin      = 1
	LimitMax      = 100
	LimitDefault  = 1
	OffsetDefault = 0
)

var (
	ErrInvalidParams = errors.New("invalid request params")
)

type JoinQueryArguments struct {
	Query     string `json:"query"`
	Arguments []any  `json:"arguments"`
}

type QueryRequest struct {
	Tx        *gorm.DB              `json:"-"`          // db transaction
	UseSlave  bool                  `json:"use_slave"`  // use slave
	Alias     string                `json:"alias"`      // table alias
	Fields    string                `json:"fields"`     // fields, default: *
	Query     string                `json:"query"`      // query, like: id = ? and status = ?
	Arguments []any                 `json:"arguments"`  // arguments like: []any{1, 1}
	ForUpdate bool                  `json:"for_update"` // for update
	Joins     []*JoinQueryArguments `json:"joins"`      // join query
}

func (r *QueryRequest) Validate() (err error) {
	if r == nil || (r.Query == "" && len(r.Arguments) > 0) {
		return ErrInvalidParams
	}

	if r.Fields == "" {
		r.Fields = "*"
	}

	return
}

type QueryListRequest struct {
	Tx        *gorm.DB              `json:"-"`         // db transaction
	UseSlave  bool                  `json:"use_slave"` // use slave
	Alias     string                `json:"alias"`     // table alias
	Fields    string                `json:"fields"`    // fields, default: *
	Query     string                `json:"query"`     // query, like: id = ? and status = ?
	Arguments []any                 `json:"arguments"` // arguments like: []any{1, 1}
	Limit     int                   `json:"limit"`
	Offset    int                   `json:"offset"`
	Order     string                `json:"order"`
	Group     string                `json:"group"`
	ForUpdate bool                  `json:"for_update"`
	Joins     []*JoinQueryArguments `json:"joins"` // join query
}

func (r *QueryListRequest) Validate() (err error) {
	if r == nil || (r.Query == "" && len(r.Arguments) > 0) {
		return ErrInvalidParams
	}

	if r.Limit < LimitMin || r.Limit > LimitMax {
		r.Limit = LimitDefault
	}

	if r.Offset < OffsetDefault {
		r.Offset = OffsetDefault
	}

	if r.Fields == "" {
		r.Fields = "*"
	}

	return
}

type QueryCountRequest struct {
	Tx        *gorm.DB              `json:"-"`         // db transaction
	UseSlave  bool                  `json:"use_slave"` // use slave
	Alias     string                `json:"alias"`     // table alias
	Query     string                `json:"query"`     // query, like: id = ? and status = ?
	Arguments []any                 `json:"arguments"` // arguments like: []any{1, 1}
	Joins     []*JoinQueryArguments `json:"joins"`     // join query
	ForUpdate bool                  `json:"for_update"`
}

func (r *QueryCountRequest) Validate() (err error) {
	if r == nil || (r.Query == "" && len(r.Arguments) > 0) {
		return ErrInvalidParams
	}
	return
}

type CreateRequest[T ModelInterface] struct {
	Tx         *gorm.DB          `json:"-"`
	Data       T                 `json:"data"`
	IsDataZero func(data T) bool `json:"-"`
}

func (r *CreateRequest[T]) Validate() (err error) {
	if r == nil {
		return ErrInvalidParams
	}

	if r.IsDataZero != nil && r.IsDataZero(r.Data) {
		return ErrInvalidParams
	}

	return
}

type Updater[T ModelInterface] interface {
	GetTx() *gorm.DB
	GetFields() []string
	GetQuery() string
	GetArguments() []any
	GetLimit() int
	GetData() any
	Validate() (err error)
}

type UpdateRequest[T ModelInterface] struct {
	Tx         *gorm.DB
	Fields     []string
	Query      string
	Arguments  []any
	Data       T
	IsDataZero func(data T) bool
	Limit      int
}

func (r *UpdateRequest[T]) GetTx() *gorm.DB     { return r.Tx }
func (r *UpdateRequest[T]) GetFields() []string { return r.Fields }
func (r *UpdateRequest[T]) GetQuery() string    { return r.Query }
func (r *UpdateRequest[T]) GetArguments() []any { return r.Arguments }
func (r *UpdateRequest[T]) GetLimit() int       { return r.Limit }
func (r *UpdateRequest[T]) GetData() any        { return r.Data }

func (r *UpdateRequest[T]) Validate() (err error) {
	if r == nil || (r.Query == "" && len(r.Arguments) > 0) {
		return ErrInvalidParams
	}

	if r.IsDataZero != nil && r.IsDataZero(r.Data) {
		return ErrInvalidParams
	}

	return
}

type UpdateRequestByMap[T ModelInterface] struct {
	Tx        *gorm.DB
	Fields    []string
	Query     string
	Arguments []any
	Data      map[string]any
	Limit     int
}

func (r *UpdateRequestByMap[T]) GetTx() *gorm.DB     { return r.Tx }
func (r *UpdateRequestByMap[T]) GetFields() []string { return r.Fields }
func (r *UpdateRequestByMap[T]) GetQuery() string    { return r.Query }
func (r *UpdateRequestByMap[T]) GetArguments() []any { return r.Arguments }
func (r *UpdateRequestByMap[T]) GetLimit() int       { return r.Limit }
func (r *UpdateRequestByMap[T]) GetData() any        { return r.Data }

func (r *UpdateRequestByMap[T]) Validate() (err error) {
	if r == nil || (r.Query == "" && len(r.Arguments) > 0) {
		return ErrInvalidParams
	}

	if len(r.Data) == 0 {
		return ErrInvalidParams
	}

	if r.Data == nil {
		return ErrInvalidParams
	}

	return
}

type DeleteRequest[T ModelInterface] struct {
	Tx        *gorm.DB `json:"-"`
	Query     string   `json:"query"`
	Arguments []any    `json:"arguments"`
	Limit     int      `json:"limit"`
}

func (r *DeleteRequest[T]) Validate() (err error) {
	if r == nil || (r.Query == "" && len(r.Arguments) > 0) {
		return ErrInvalidParams
	}
	return
}

type UpsertRequest[T ModelInterface] struct {
	Tx         *gorm.DB
	Data       T
	OnConflict clause.OnConflict
	IsDataZero func(data T) bool
}

func (r *UpsertRequest[T]) Validate() (err error) {
	if r == nil || len(r.OnConflict.Columns) == 0 || len(r.OnConflict.DoUpdates) == 0 {
		return ErrInvalidParams
	}

	if r.IsDataZero != nil && r.IsDataZero(r.Data) {
		return ErrInvalidParams
	}

	return
}
