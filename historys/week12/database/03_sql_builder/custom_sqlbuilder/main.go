package main

/*
 - 自行实现一个函数链式的sql构建器
 - 一个小demo，还有诸多不完善的地方
*/

import (
	"fmt"
	"strconv"
	"strings"
)

// Builder 根据一个函数生成一小段sql
type Builder interface { // select、where、limit、orderby这些都是Builder
	toString() string
	getPrev() Builder
}

type LimitBuilder struct {
	sb   strings.Builder
	prev Builder // 前面的Builder
}

func newLimitBuilder(offset, n int) *LimitBuilder {
	builder := &LimitBuilder{}
	// 通过strings.Builder实现高效的字符串连接
	builder.sb.WriteString(" limit ")
	builder.sb.WriteString(strconv.Itoa(offset))
	builder.sb.WriteString(",")
	builder.sb.WriteString(strconv.Itoa(n))
	return builder
}

func (self *LimitBuilder) toString() string {
	return self.sb.String()
}

func (self *LimitBuilder) getPrev() Builder {
	return self.prev
}

func (self *LimitBuilder) ToString() string {
	var root Builder
	root = self
	for root.getPrev() != nil { // 找到最前面的root Builder
		root = root.getPrev()
	}
	return root.toString()
}

type OrderByBuilder struct {
	sb    strings.Builder
	limit *LimitBuilder
	prev  Builder
}

func newOrderByBuilder(column string) *OrderByBuilder {
	builder := &OrderByBuilder{}
	builder.sb.WriteString(" order by ")
	builder.sb.WriteString(column)
	return builder
}

func (self *OrderByBuilder) getPrev() Builder {
	return self.prev
}

func (self *OrderByBuilder) toString() string {
	if self.limit != nil {
		self.sb.WriteString(self.limit.toString())
	}
	return self.sb.String()
}

func (self *OrderByBuilder) ToString() string {
	var root Builder
	root = self
	for root.getPrev() != nil {
		root = root.getPrev()
	}
	return root.toString()
}

func (self *OrderByBuilder) Asc() *OrderByBuilder {
	self.sb.WriteString(" asc")
	return self
}

func (self *OrderByBuilder) Desc() *OrderByBuilder {
	self.sb.WriteString(" desc")
	return self
}

// orderby后面可以接limit
func (self *OrderByBuilder) Limit(offset, n int) *LimitBuilder {
	limit := newLimitBuilder(offset, n)
	limit.prev = self
	self.limit = limit
	return limit
}

type WhereBuilder struct {
	sb      strings.Builder
	orderby *OrderByBuilder
	limit   *LimitBuilder
	prev    Builder
}

func newWhereBuilder(condition string) *WhereBuilder {
	builder := &WhereBuilder{}
	builder.sb.WriteString(" where ")
	builder.sb.WriteString(condition)
	return builder
}

func (self *WhereBuilder) getPrev() Builder {
	return self.prev
}

func (self *WhereBuilder) toString() string {
	// 递归调用后续Builder的ToString()
	if self.orderby != nil {
		self.sb.Write([]byte(self.orderby.toString()))
	}
	if self.limit != nil {
		self.sb.Write([]byte(self.limit.toString()))
	}
	return self.sb.String()
}

func (self *WhereBuilder) ToString() string {
	var root Builder
	root = self
	for root.getPrev() != nil {
		root = root.getPrev()
	}
	return root.toString()
}

// And和Or都是where里的可选部分，它们的地位平等，都返回WhereBuilder
func (self *WhereBuilder) And(condition string) *WhereBuilder {
	self.sb.WriteString(" and ")
	self.sb.WriteString(condition)
	return self
}

// And和Or都是where里的可选部分，它们的地位平等，都返回WhereBuilder
func (self *WhereBuilder) Or(condition string) *WhereBuilder {
	self.sb.WriteString(" or ")
	self.sb.WriteString(condition)
	return self
}

// where后面可以接order by
func (self *WhereBuilder) OrderBy(column string) *OrderByBuilder {
	orderby := newOrderByBuilder(column)
	self.orderby = orderby
	orderby.prev = self
	return orderby
}

// where后面可以接limit
func (self *WhereBuilder) Limit(offset, n int) *LimitBuilder {
	limit := newLimitBuilder(offset, n)
	limit.prev = self
	self.limit = limit
	return limit
}

type SelectBuilder struct {
	sb    strings.Builder
	table string
	where *WhereBuilder
}

func NewSelectBuilder(table string) *SelectBuilder {
	builder := &SelectBuilder{
		table: table,
	}
	builder.sb.WriteString("select ")
	return builder
}

func (self *SelectBuilder) getPrev() Builder {
	return nil
}

func (self *SelectBuilder) toString() string {
	if self.where != nil {
		self.sb.Write([]byte(self.where.toString()))
	}
	return self.sb.String()
}

func (self *SelectBuilder) ToString() string {
	var root Builder
	root = self
	for root.getPrev() != nil {
		root = root.getPrev()
	}
	return root.toString()
}

// 通过select查询哪几列
func (self *SelectBuilder) Column(columns string) *SelectBuilder {
	self.sb.WriteString(columns)
	self.sb.WriteString(" from ")
	self.sb.WriteString(self.table)
	return self
}

func (self *SelectBuilder) Where(condition string) *WhereBuilder {
	where := newWhereBuilder(condition)
	self.where = where
	where.prev = self
	return where
}

func main() {
	// Where、OrderBy、Limit有没有都不影响调用ToString()；Where里的And和Or有没有都不影响调用ToString()
	sql := NewSelectBuilder("student").Column("id,name,city").
		Where("id>0").
		And("city='郑州'").
		Or("city='北京'").
		OrderBy("score").Desc().
		Limit(0, 10).ToString()
	fmt.Println(sql)
}
