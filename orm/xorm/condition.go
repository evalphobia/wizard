package xorm

type FindCondition struct {
	Table   interface{}
	Where   []Where
	WhereIn []Where
	OrderBy []Order
	Limit   int
	Offset  int
}

type Where struct {
	Statement string
	Args      []interface{}
}

type Order struct {
	Name        string
	OrderByDesc bool
}

func NewFindCondition(table interface{}) FindCondition {
	return FindCondition{
		Table: table,
	}
}

func (c *FindCondition) And(s string, args ...interface{}) {
	w := Where{
		Statement: s,
		Args:      args,
	}
	c.Where = append(c.Where, w)
}

func (c *FindCondition) In(s string, args ...interface{}) {
	w := Where{
		Statement: s,
		Args:      args,
	}
	c.WhereIn = append(c.WhereIn, w)
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
