package xorm

type Where struct {
	Statement string
	Args      []interface{}
}

func NewWhere(s string, args ...interface{}) Where {
	return Where{
		Statement: s,
		Args:      args,
	}
}

type Order struct {
	Name        string
	OrderByDesc bool
}

// FindCondition is conditions for FindParallel
type FindCondition struct {
	Table   interface{}
	Where   []Where
	WhereIn []Where
	OrderBy []Order
	Limit   int
	Offset  int
}

func NewFindCondition(table interface{}) FindCondition {
	return FindCondition{
		Table: table,
	}
}

func (c *FindCondition) And(s string, args ...interface{}) {
	c.Where = append(c.Where, NewWhere(s, args...))
}

func (c *FindCondition) In(s string, args ...interface{}) {
	c.WhereIn = append(c.WhereIn, NewWhere(s, args...))
}

func (c *FindCondition) OrderByAsc(s string) {
	o := Order{
		Name: s,
	}
	c.OrderBy = append(c.OrderBy, o)
}

func (c *FindCondition) OrderByDesc(s string) {
	o := Order{
		Name:        s,
		OrderByDesc: true,
	}
	c.OrderBy = append(c.OrderBy, o)
}

func (c *FindCondition) SetLimit(i int) {
	c.Limit = i
}

func (c *FindCondition) SetOffset(i int) {
	c.Offset = i
}

// UpdateCondition is conditions for UpdateParallel
type UpdateCondition struct {
	Table           interface{}
	Where           []Where
	WhereIn         []Where
	AllColumns      bool
	Columns         []string
	MustColumns     []string
	OmitColumns     []string
	NullableColumns []string
	Increments      []Where
	Decrements      []Where
}

func NewUpdateCondition(table interface{}) UpdateCondition {
	return UpdateCondition{
		Table: table,
	}
}

func (c *UpdateCondition) And(s string, args ...interface{}) {
	c.Where = append(c.Where, NewWhere(s, args...))
}

func (c *UpdateCondition) In(s string, args ...interface{}) {
	c.WhereIn = append(c.WhereIn, NewWhere(s, args...))
}

func (c *UpdateCondition) AllCols() {
	c.AllColumns = true
}

func (c *UpdateCondition) Cols(cols ...string) {
	c.Columns = append(c.Columns, cols...)
}

func (c *UpdateCondition) MustCols(cols ...string) {
	c.MustColumns = append(c.MustColumns, cols...)
}

func (c *UpdateCondition) Omit(cols ...string) {
	c.OmitColumns = append(c.OmitColumns, cols...)
}

func (c *UpdateCondition) Nullable(cols ...string) {
	c.NullableColumns = append(c.NullableColumns, cols...)
}

func (c *UpdateCondition) Incr(s string, args ...interface{}) {
	c.Increments = append(c.Increments, NewWhere(s, args...))
}

func (c *UpdateCondition) Decr(s string, args ...interface{}) {
	c.Decrements = append(c.Decrements, NewWhere(s, args...))
}
