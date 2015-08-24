package goqu

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/c2fo/testify/assert"
	"github.com/c2fo/testify/suite"
)

type datasetTest struct {
	suite.Suite
}

func (me *datasetTest) Truncate(buf *SqlBuilder) *SqlBuilder {
	buf.Truncate(0)
	buf.args = make([]interface{}, 0)
	return buf
}

func (me *datasetTest) TestClone() {
	t := me.T()
	ds := From("test")
	assert.Equal(t, ds.Clone(), ds)
}

func (me *datasetTest) TestExpression() {
	t := me.T()
	ds := From("test")
	assert.Equal(t, ds.Expression(), ds)
}

func (me *datasetTest) TestAdapter() {
	t := me.T()
	ds := From("test")
	assert.Equal(t, ds.Adapter(), ds.adapter)
}

func (me *datasetTest) TestSetAdapter() {
	t := me.T()
	ds := From("test")
	adapter := NewAdapter("default", ds)
	ds.SetAdapter(adapter)
	assert.Equal(t, ds.Adapter(), adapter)
}

func (me *datasetTest) TestPrepared() {
	t := me.T()
	ds := From("test")
	preparedDs := ds.Prepared(true)
	assert.True(t, preparedDs.isPrepared)
	assert.False(t, ds.isPrepared)

	//should apply the prepared to any datasets created from the root
	assert.True(t, preparedDs.Where(Ex{"a": 1}).isPrepared)
}

func (me *datasetTest) TestLiteralUnsupportedType() {
	t := me.T()
	assert.EqualError(t, From("test").Literal(NewSqlBuilder(false), struct{}{}), "goqu: Unable to encode value {}")
}

type unknowExpression struct {
}

func (me unknowExpression) Expression() Expression {
	return me
}
func (me unknowExpression) Clone() Expression {
	return me
}
func (me *datasetTest) TestLiteralUnsupportedExpression() {
	t := me.T()
	assert.EqualError(t, From("test").Literal(NewSqlBuilder(false), unknowExpression{}), "goqu: Unsupported expression type goqu.unknowExpression")
}

func (me *datasetTest) TestLiteralFloatTypes() {
	t := me.T()
	ds := From("test")
	var float float64
	buf := NewSqlBuilder(false)
	assert.NoError(t, ds.Literal(buf, float32(10.01)))
	assert.Equal(t, buf.String(), "10.010000228881836")
	assert.NoError(t, ds.Literal(me.Truncate(buf), float64(10.01)))
	assert.Equal(t, buf.String(), "10.01")
	assert.NoError(t, ds.Literal(me.Truncate(buf), &float))
	assert.Equal(t, buf.String(), "0")

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(buf, float32(10.01)))
	assert.Equal(t, buf.args, []interface{}{float64(float32(10.01))})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), float64(10.01)))
	assert.Equal(t, buf.args, []interface{}{float64(10.01)})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), &float))
	assert.Equal(t, buf.args, []interface{}{float})
	assert.Equal(t, buf.String(), "?")
}

func (me *datasetTest) TestLiteralIntTypes() {
	t := me.T()
	ds := From("test")
	var i int64
	buf := NewSqlBuilder(false)
	assert.NoError(t, ds.Literal(buf, int(10)))
	assert.Equal(t, buf.String(), "10")
	assert.NoError(t, ds.Literal(me.Truncate(buf), int8(10)))
	assert.Equal(t, buf.String(), "10")
	assert.NoError(t, ds.Literal(me.Truncate(buf), int16(10)))
	assert.Equal(t, buf.String(), "10")
	assert.NoError(t, ds.Literal(me.Truncate(buf), int32(10)))
	assert.Equal(t, buf.String(), "10")
	assert.NoError(t, ds.Literal(me.Truncate(buf), int64(10)))
	assert.Equal(t, buf.String(), "10")
	assert.NoError(t, ds.Literal(me.Truncate(buf), uint(10)))
	assert.Equal(t, buf.String(), "10")
	assert.NoError(t, ds.Literal(me.Truncate(buf), uint8(10)))
	assert.Equal(t, buf.String(), "10")
	assert.NoError(t, ds.Literal(me.Truncate(buf), uint16(10)))
	assert.Equal(t, buf.String(), "10")
	assert.NoError(t, ds.Literal(me.Truncate(buf), uint32(10)))
	assert.Equal(t, buf.String(), "10")
	assert.NoError(t, ds.Literal(me.Truncate(buf), uint64(10)))
	assert.Equal(t, buf.String(), "10")
	assert.NoError(t, ds.Literal(me.Truncate(buf), &i))
	assert.Equal(t, buf.String(), "0")

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(buf, int(10)))
	assert.Equal(t, buf.args, []interface{}{int64(10)})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), int8(10)))
	assert.Equal(t, buf.args, []interface{}{int64(10)})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), int16(10)))
	assert.Equal(t, buf.args, []interface{}{int64(10)})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), int32(10)))
	assert.Equal(t, buf.args, []interface{}{int64(10)})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), int64(10)))
	assert.Equal(t, buf.args, []interface{}{int64(10)})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), uint(10)))
	assert.Equal(t, buf.args, []interface{}{int64(10)})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), uint8(10)))
	assert.Equal(t, buf.args, []interface{}{int64(10)})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), uint16(10)))
	assert.Equal(t, buf.args, []interface{}{int64(10)})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), uint32(10)))
	assert.Equal(t, buf.args, []interface{}{int64(10)})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), uint64(10)))
	assert.Equal(t, buf.args, []interface{}{int64(10)})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), &i))
	assert.Equal(t, buf.args, []interface{}{i})
	assert.Equal(t, buf.String(), "?")
}

func (me *datasetTest) TestLiteralStringTypes() {
	t := me.T()
	ds := From("test")
	var str string
	buf := NewSqlBuilder(false)
	assert.NoError(t, ds.Literal(me.Truncate(buf), "Hello"))
	assert.Equal(t, buf.String(), "'Hello'")
	//should esacpe single quotes
	assert.NoError(t, ds.Literal(me.Truncate(buf), "hello'"))
	assert.Equal(t, buf.String(), "'hello'''")
	assert.NoError(t, ds.Literal(me.Truncate(buf), &str))
	assert.Equal(t, buf.String(), "''")

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), "Hello"))
	assert.Equal(t, buf.args, []interface{}{"Hello"})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), "hello'"))
	assert.Equal(t, buf.args, []interface{}{"hello'"})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), &str))
	assert.Equal(t, buf.args, []interface{}{str})
	assert.Equal(t, buf.String(), "?")
}

func (me *datasetTest) TestLiteralBytesTypes() {
	t := me.T()
	ds := From("test")
	var b string
	buf := NewSqlBuilder(false)
	assert.NoError(t, ds.Literal(me.Truncate(buf), []byte("Hello")))
	assert.Equal(t, buf.Bytes(), []byte("'Hello'"))
	//should escape single quotes
	assert.NoError(t, ds.Literal(me.Truncate(buf), []byte("hello'")))
	assert.Equal(t, buf.Bytes(), []byte("'hello'''"))
	assert.NoError(t, ds.Literal(me.Truncate(buf), (&b)))
	assert.Equal(t, buf.Bytes(), []byte("''"))
	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), []byte("Hello")))
	assert.Equal(t, buf.args, []interface{}{[]byte("Hello")})
	assert.Equal(t, buf.Bytes(), []byte("?"))
	assert.NoError(t, ds.Literal(me.Truncate(buf), []byte("hello'")))
	assert.Equal(t, buf.args, []interface{}{[]byte("hello'")})
	assert.Equal(t, buf.Bytes(), []byte("?"))
	assert.NoError(t, ds.Literal(me.Truncate(buf), []byte(*(&b))))
	assert.Equal(t, buf.args, []interface{}{[]byte(b)})
	assert.Equal(t, buf.Bytes(), []byte("?"))
}

func (me *datasetTest) TestLiteralBoolTypes() {
	t := me.T()
	var b bool
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), true))
	assert.Equal(t, buf.String(), "TRUE")
	assert.NoError(t, ds.Literal(me.Truncate(buf), false))
	assert.Equal(t, buf.String(), "FALSE")
	assert.NoError(t, ds.Literal(me.Truncate(buf), &b))
	assert.Equal(t, buf.String(), "FALSE")

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), true))
	assert.Equal(t, buf.args, []interface{}{true})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), false))
	assert.Equal(t, buf.args, []interface{}{false})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), &b))
	assert.Equal(t, buf.args, []interface{}{b})
	assert.Equal(t, buf.String(), "?")
}

func (me *datasetTest) TestLiteralTimeTypes() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	now := time.Now().UTC()
	assert.NoError(t, ds.Literal(me.Truncate(buf), now))
	assert.Equal(t, buf.String(), "'"+now.Format(time.RFC3339Nano)+"'")
	assert.NoError(t, ds.Literal(me.Truncate(buf), &now))
	assert.Equal(t, buf.String(), "'"+now.Format(time.RFC3339Nano)+"'")

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), now))
	assert.Equal(t, buf.args, []interface{}{now})
	assert.Equal(t, buf.String(), "?")
	assert.NoError(t, ds.Literal(me.Truncate(buf), &now))
	assert.Equal(t, buf.args, []interface{}{now})
	assert.Equal(t, buf.String(), "?")
}

func (me *datasetTest) TestLiteralNilTypes() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), nil))
	assert.Equal(t, buf.String(), "NULL")

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), nil))
	assert.Equal(t, buf.args, []interface{}{nil})
	assert.Equal(t, buf.String(), "?")
}

type datasetValuerType int64

func (j datasetValuerType) Value() (driver.Value, error) {
	return []byte(fmt.Sprintf("Hello World %d", j)), nil
}

func (me *datasetTest) TestLiteralValuer() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")

	assert.NoError(t, ds.Literal(me.Truncate(buf), datasetValuerType(10)))
	assert.Equal(t, buf.String(), "'Hello World 10'")

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), datasetValuerType(10)))
	assert.Equal(t, buf.args, []interface{}{[]byte("Hello World 10")})
	assert.Equal(t, buf.String(), "?")

}

func (me *datasetTest) TestLiteraSlice() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), []string{"a", "b", "c"}))
	assert.Equal(t, buf.String(), `('a', 'b', 'c')`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), []string{"a", "b", "c"}))
	assert.Equal(t, buf.args, []interface{}{"a", "b", "c"})
	assert.Equal(t, buf.String(), `(?, ?, ?)`)
}

func (me *datasetTest) TestLiteralDataset() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), From("a")))
	assert.Equal(t, buf.String(), `(SELECT * FROM "a")`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), From("a").As("b")))
	assert.Equal(t, buf.String(), `(SELECT * FROM "a") AS "b"`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), From("a")))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `(SELECT * FROM "a")`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), From("a").As("b")))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `(SELECT * FROM "a") AS "b"`)
}

func (me *datasetTest) TestLiteralColumnList() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), cols("a", Literal("true"))))
	assert.Equal(t, buf.String(), `"a", true`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), cols("a", Literal("true"))))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `"a", true`)
}

func (me *datasetTest) TestLiteralExpressionList() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), And(I("a").Eq("b"), I("c").Neq(1))))
	assert.Equal(t, buf.String(), `(("a" = 'b') AND ("c" != 1))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Or(I("a").Eq("b"), I("c").Neq(1))))
	assert.Equal(t, buf.String(), `(("a" = 'b') OR ("c" != 1))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Or(I("a").Eq("b"), And(I("c").Neq(1), I("d").Eq(Literal("NOW()"))))))
	assert.Equal(t, buf.String(), `(("a" = 'b') OR (("c" != 1) AND ("d" = NOW())))`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), And(I("a").Eq("b"), I("c").Neq(1))))
	assert.Equal(t, buf.args, []interface{}{"b", 1})
	assert.Equal(t, buf.String(), `(("a" = ?) AND ("c" != ?))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Or(I("a").Eq("b"), I("c").Neq(1))))
	assert.Equal(t, buf.args, []interface{}{"b", 1})
	assert.Equal(t, buf.String(), `(("a" = ?) OR ("c" != ?))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Or(I("a").Eq("b"), And(I("c").Neq(1), I("d").Eq(Literal("NOW()"))))))
	assert.Equal(t, buf.args, []interface{}{"b", 1})
	assert.Equal(t, buf.String(), `(("a" = ?) OR (("c" != ?) AND ("d" = NOW())))`)
}

func (me *datasetTest) TestLiteralLiteralExpression() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")

	assert.NoError(t, ds.Literal(me.Truncate(buf), Literal(`"b"::DATE = '2010-09-02'`)))
	assert.Equal(t, buf.String(), `"b"::DATE = '2010-09-02'`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Literal(`"b" = ? or "c" = ? or d IN ?`, "a", 1, []int{1, 2, 3, 4})))
	assert.Equal(t, buf.String(), `"b" = 'a' or "c" = 1 or d IN (1, 2, 3, 4)`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Literal(`"b"::DATE = '2010-09-02'`)))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `"b"::DATE = '2010-09-02'`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Literal(`"b" = ? or "c" = ? or d IN ?`, "a", 1, []int{1, 2, 3, 4})))
	assert.Equal(t, buf.args, []interface{}{"a", 1, 1, 2, 3, 4})
	assert.Equal(t, buf.String(), `"b" = ? or "c" = ? or d IN (?, ?, ?, ?)`)
}

func (me *datasetTest) TestLiteralAliasedExpression() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").As("b")))
	assert.Equal(t, buf.String(), `"a" AS "b"`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Literal("count(*)").As("count")))
	assert.Equal(t, buf.String(), `count(*) AS "count"`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").As("b")))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `"a" AS "b"`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Literal("count(*)").As("count")))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `count(*) AS "count"`)

	buf = NewSqlBuilder(false)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").As(I("b"))))
	assert.Equal(t, buf.String(), `"a" AS "b"`)
}

func (me *datasetTest) TestBooleanExpression() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Eq(1)))
	assert.Equal(t, buf.String(), `("a" = 1)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Eq(true)))
	assert.Equal(t, buf.String(), `("a" IS TRUE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Eq(false)))
	assert.Equal(t, buf.String(), `("a" IS FALSE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Eq(nil)))
	assert.Equal(t, buf.String(), `("a" IS NULL)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Eq([]int64{1, 2, 3})))
	assert.Equal(t, buf.String(), `("a" IN (1, 2, 3))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Eq(From("test2").Select("id"))))
	assert.Equal(t, buf.String(), `("a" IN (SELECT "id" FROM "test2"))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Neq(1)))
	assert.Equal(t, buf.String(), `("a" != 1)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Neq(true)))
	assert.Equal(t, buf.String(), `("a" IS NOT TRUE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Neq(false)))
	assert.Equal(t, buf.String(), `("a" IS NOT FALSE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Neq(nil)))
	assert.Equal(t, buf.String(), `("a" IS NOT NULL)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Neq([]int64{1, 2, 3})))
	assert.Equal(t, buf.String(), `("a" NOT IN (1, 2, 3))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Neq(From("test2").Select("id"))))
	assert.Equal(t, buf.String(), `("a" NOT IN (SELECT "id" FROM "test2"))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Is(nil)))
	assert.Equal(t, buf.String(), `("a" IS NULL)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Is(false)))
	assert.Equal(t, buf.String(), `("a" IS FALSE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Is(true)))
	assert.Equal(t, buf.String(), `("a" IS TRUE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").IsNot(nil)))
	assert.Equal(t, buf.String(), `("a" IS NOT NULL)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").IsNot(false)))
	assert.Equal(t, buf.String(), `("a" IS NOT FALSE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").IsNot(true)))
	assert.Equal(t, buf.String(), `("a" IS NOT TRUE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Gt(1)))
	assert.Equal(t, buf.String(), `("a" > 1)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Gte(1)))
	assert.Equal(t, buf.String(), `("a" >= 1)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Lt(1)))
	assert.Equal(t, buf.String(), `("a" < 1)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Lte(1)))
	assert.Equal(t, buf.String(), `("a" <= 1)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").In([]int{1, 2, 3})))
	assert.Equal(t, buf.String(), `("a" IN (1, 2, 3))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").NotIn([]int{1, 2, 3})))
	assert.Equal(t, buf.String(), `("a" NOT IN (1, 2, 3))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Like("a%")))
	assert.Equal(t, buf.String(), `("a" LIKE 'a%')`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Like(regexp.MustCompile("(a|b)"))))
	assert.Equal(t, buf.String(), `("a" ~ '(a|b)')`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").NotLike("a%")))
	assert.Equal(t, buf.String(), `("a" NOT LIKE 'a%')`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").NotLike(regexp.MustCompile("(a|b)"))))
	assert.Equal(t, buf.String(), `("a" !~ '(a|b)')`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").ILike("a%")))
	assert.Equal(t, buf.String(), `("a" ILIKE 'a%')`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").ILike(regexp.MustCompile("(a|b)"))))
	assert.Equal(t, buf.String(), `("a" ~* '(a|b)')`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").NotILike("a%")))
	assert.Equal(t, buf.String(), `("a" NOT ILIKE 'a%')`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").NotILike(regexp.MustCompile("(a|b)"))))
	assert.Equal(t, buf.String(), `("a" !~* '(a|b)')`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Eq(1)))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `("a" = ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Eq(true)))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `("a" IS TRUE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Eq(false)))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `("a" IS FALSE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Eq(nil)))
	assert.Equal(t, buf.args, []interface{}{nil})
	assert.Equal(t, buf.String(), `("a" IS ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Eq([]int64{1, 2, 3})))
	assert.Equal(t, buf.args, []interface{}{1, 2, 3})
	assert.Equal(t, buf.String(), `("a" IN (?, ?, ?))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Neq(1)))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `("a" != ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Neq(true)))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `("a" IS NOT TRUE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Neq(false)))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `("a" IS NOT FALSE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Neq(nil)))
	assert.Equal(t, buf.args, []interface{}{nil})
	assert.Equal(t, buf.String(), `("a" IS NOT ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Neq([]int64{1, 2, 3})))
	assert.Equal(t, buf.args, []interface{}{1, 2, 3})
	assert.Equal(t, buf.String(), `("a" NOT IN (?, ?, ?))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Is(nil)))
	assert.Equal(t, buf.args, []interface{}{nil})
	assert.Equal(t, buf.String(), `("a" IS ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Is(false)))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `("a" IS FALSE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Is(true)))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `("a" IS TRUE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").IsNot(nil)))
	assert.Equal(t, buf.args, []interface{}{nil})
	assert.Equal(t, buf.String(), `("a" IS NOT ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").IsNot(false)))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `("a" IS NOT FALSE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").IsNot(true)))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `("a" IS NOT TRUE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Gt(1)))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `("a" > ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Gte(1)))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `("a" >= ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Lt(1)))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `("a" < ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Lte(1)))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `("a" <= ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").In([]int{1, 2, 3})))
	assert.Equal(t, buf.args, []interface{}{1, 2, 3})
	assert.Equal(t, buf.String(), `("a" IN (?, ?, ?))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").NotIn([]int{1, 2, 3})))
	assert.Equal(t, buf.args, []interface{}{1, 2, 3})
	assert.Equal(t, buf.String(), `("a" NOT IN (?, ?, ?))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Like("a%")))
	assert.Equal(t, buf.args, []interface{}{"a%"})
	assert.Equal(t, buf.String(), `("a" LIKE ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Like(regexp.MustCompile("(a|b)"))))
	assert.Equal(t, buf.args, []interface{}{"(a|b)"})
	assert.Equal(t, buf.String(), `("a" ~ ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").NotLike("a%")))
	assert.Equal(t, buf.args, []interface{}{"a%"})
	assert.Equal(t, buf.String(), `("a" NOT LIKE ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").NotLike(regexp.MustCompile("(a|b)"))))
	assert.Equal(t, buf.args, []interface{}{"(a|b)"})
	assert.Equal(t, buf.String(), `("a" !~ ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").ILike("a%")))
	assert.Equal(t, buf.args, []interface{}{"a%"})
	assert.Equal(t, buf.String(), `("a" ILIKE ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").ILike(regexp.MustCompile("(a|b)"))))
	assert.Equal(t, buf.args, []interface{}{"(a|b)"})
	assert.Equal(t, buf.String(), `("a" ~* ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").NotILike("a%")))
	assert.Equal(t, buf.args, []interface{}{"a%"})
	assert.Equal(t, buf.String(), `("a" NOT ILIKE ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").NotILike(regexp.MustCompile("(a|b)"))))
	assert.Equal(t, buf.args, []interface{}{"(a|b)"})
	assert.Equal(t, buf.String(), `("a" !~* ?)`)

}

func (me *datasetTest) TestLiteralOrderedExpression() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Asc()))
	assert.Equal(t, buf.String(), `"a" ASC`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Desc()))
	assert.Equal(t, buf.String(), `"a" DESC`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Asc().NullsLast()))
	assert.Equal(t, buf.String(), `"a" ASC NULLS LAST`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Desc().NullsLast()))
	assert.Equal(t, buf.String(), `"a" DESC NULLS LAST`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Asc().NullsFirst()))
	assert.Equal(t, buf.String(), `"a" ASC NULLS FIRST`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Desc().NullsFirst()))
	assert.Equal(t, buf.String(), `"a" DESC NULLS FIRST`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Asc()))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `"a" ASC`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Desc()))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `"a" DESC`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Asc().NullsLast()))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `"a" ASC NULLS LAST`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Desc().NullsLast()))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `"a" DESC NULLS LAST`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Asc().NullsFirst()))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `"a" ASC NULLS FIRST`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Desc().NullsFirst()))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `"a" DESC NULLS FIRST`)
}

func (me *datasetTest) TestLiteralUpdateExpression() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Set(1)))
	assert.Equal(t, buf.String(), `"a"=1`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Set(1)))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `"a"=?`)
}

func (me *datasetTest) TestLiteralSqlFunctionExpression() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), Func("MIN", I("a"))))
	assert.Equal(t, buf.String(), `MIN("a")`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), MIN("a")))
	assert.Equal(t, buf.String(), `MIN("a")`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), COALESCE(I("a"), "a")))
	assert.Equal(t, buf.String(), `COALESCE("a", 'a')`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Func("MIN", I("a"))))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `MIN("a")`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), MIN("a")))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `MIN("a")`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), COALESCE(I("a"), "a")))
	assert.Equal(t, buf.args, []interface{}{"a"})
	assert.Equal(t, buf.String(), `COALESCE("a", ?)`)
}

func (me *datasetTest) TestLiteralCastExpression() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Cast("DATE")))
	assert.Equal(t, buf.String(), `CAST("a" AS DATE)`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a").Cast("DATE")))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `CAST("a" AS DATE)`)
}

func (me *datasetTest) TestCompoundExpression() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), Union(From("b"))))
	assert.Equal(t, buf.String(), ` UNION (SELECT * FROM "b")`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), UnionAll(From("b"))))
	assert.Equal(t, buf.String(), ` UNION ALL (SELECT * FROM "b")`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Intersect(From("b"))))
	assert.Equal(t, buf.String(), ` INTERSECT (SELECT * FROM "b")`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), IntersectAll(From("b"))))
	assert.Equal(t, buf.String(), ` INTERSECT ALL (SELECT * FROM "b")`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Union(From("b"))))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), ` UNION (SELECT * FROM "b")`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), UnionAll(From("b"))))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), ` UNION ALL (SELECT * FROM "b")`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Intersect(From("b"))))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), ` INTERSECT (SELECT * FROM "b")`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), IntersectAll(From("b"))))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), ` INTERSECT ALL (SELECT * FROM "b")`)
}

func (me *datasetTest) TestLiteralIdentifierExpression() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a")))
	assert.Equal(t, buf.String(), `"a"`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a.b")))
	assert.Equal(t, buf.String(), `"a"."b"`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a.b.c")))
	assert.Equal(t, buf.String(), `"a"."b"."c"`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a.b.*")))
	assert.Equal(t, buf.String(), `"a"."b".*`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a.*")))
	assert.Equal(t, buf.String(), `"a".*`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a")))
	assert.Equal(t, buf.String(), `"a"`)
	assert.Equal(t, buf.args, []interface{}{})
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a.b")))
	assert.Equal(t, buf.String(), `"a"."b"`)
	assert.Equal(t, buf.args, []interface{}{})
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a.b.c")))
	assert.Equal(t, buf.String(), `"a"."b"."c"`)
	assert.Equal(t, buf.args, []interface{}{})
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a.b.*")))
	assert.Equal(t, buf.String(), `"a"."b".*`)
	assert.Equal(t, buf.args, []interface{}{})
	assert.NoError(t, ds.Literal(me.Truncate(buf), I("a.*")))
	assert.Equal(t, buf.String(), `"a".*`)
	assert.Equal(t, buf.args, []interface{}{})
}

func (me *datasetTest) TestLiteralExpressionMap() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": 1}))
	assert.Equal(t, buf.String(), `("a" = 1)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": true}))
	assert.Equal(t, buf.String(), `("a" IS TRUE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": false}))
	assert.Equal(t, buf.String(), `("a" IS FALSE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": nil}))
	assert.Equal(t, buf.String(), `("a" IS NULL)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": []string{"a", "b", "c"}}))
	assert.Equal(t, buf.String(), `("a" IN ('a', 'b', 'c'))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"neq": 1}}))
	assert.Equal(t, buf.String(), `("a" != 1)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"isnot": true}}))
	assert.Equal(t, buf.String(), `("a" IS NOT TRUE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"gt": 1}}))
	assert.Equal(t, buf.String(), `("a" > 1)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"gte": 1}}))
	assert.Equal(t, buf.String(), `("a" >= 1)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"lt": 1}}))
	assert.Equal(t, buf.String(), `("a" < 1)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"lte": 1}}))
	assert.Equal(t, buf.String(), `("a" <= 1)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"like": "a%"}}))
	assert.Equal(t, buf.String(), `("a" LIKE 'a%')`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"notLike": "a%"}}))
	assert.Equal(t, buf.String(), `("a" NOT LIKE 'a%')`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"notLike": "a%"}}))
	assert.Equal(t, buf.String(), `("a" NOT LIKE 'a%')`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"in": []string{"a", "b", "c"}}}))
	assert.Equal(t, buf.String(), `("a" IN ('a', 'b', 'c'))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"notIn": []string{"a", "b", "c"}}}))
	assert.Equal(t, buf.String(), `("a" NOT IN ('a', 'b', 'c'))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"is": nil, "eq": 10}}))
	assert.Equal(t, buf.String(), `(("a" = 10) OR ("a" IS NULL))`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": 1}))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `("a" = ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": true}))
	assert.Equal(t, buf.String(), `("a" IS TRUE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": false}))
	assert.Equal(t, buf.String(), `("a" IS FALSE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": nil}))
	assert.Equal(t, buf.args, []interface{}{nil})
	assert.Equal(t, buf.String(), `("a" IS ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": []string{"a", "b", "c"}}))
	assert.Equal(t, buf.args, []interface{}{"a", "b", "c"})
	assert.Equal(t, buf.String(), `("a" IN (?, ?, ?))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"neq": 1}}))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `("a" != ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"isnot": true}}))
	assert.Equal(t, buf.args, []interface{}{})
	assert.Equal(t, buf.String(), `("a" IS NOT TRUE)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"gt": 1}}))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `("a" > ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"gte": 1}}))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `("a" >= ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"lt": 1}}))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `("a" < ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"lte": 1}}))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `("a" <= ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"like": "a%"}}))
	assert.Equal(t, buf.args, []interface{}{"a%"})
	assert.Equal(t, buf.String(), `("a" LIKE ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"notLike": "a%"}}))
	assert.Equal(t, buf.args, []interface{}{"a%"})
	assert.Equal(t, buf.String(), `("a" NOT LIKE ?)`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"in": []string{"a", "b", "c"}}}))
	assert.Equal(t, buf.args, []interface{}{"a", "b", "c"})
	assert.Equal(t, buf.String(), `("a" IN (?, ?, ?))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"notIn": []string{"a", "b", "c"}}}))
	assert.Equal(t, buf.args, []interface{}{"a", "b", "c"})
	assert.Equal(t, buf.String(), `("a" NOT IN (?, ?, ?))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), Ex{"a": Op{"is": nil, "eq": 10}}))
	assert.Equal(t, buf.args, []interface{}{10, nil})
	assert.Equal(t, buf.String(), `(("a" = ?) OR ("a" IS ?))`)
}

func (me *datasetTest) TestLiteralExpressionOrMap() {
	t := me.T()
	buf := NewSqlBuilder(false)
	ds := From("test")
	assert.NoError(t, ds.Literal(me.Truncate(buf), ExOr{"a": 1, "b": true}))
	assert.Equal(t, buf.String(), `(("a" = 1) OR ("b" IS TRUE))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), ExOr{"a": 1, "b": []string{"a", "b", "c"}}))
	assert.Equal(t, buf.String(), `(("a" = 1) OR ("b" IN ('a', 'b', 'c')))`)

	buf = NewSqlBuilder(true)
	assert.NoError(t, ds.Literal(me.Truncate(buf), ExOr{"a": 1, "b": true}))
	assert.Equal(t, buf.args, []interface{}{1})
	assert.Equal(t, buf.String(), `(("a" = ?) OR ("b" IS TRUE))`)
	assert.NoError(t, ds.Literal(me.Truncate(buf), ExOr{"a": 1, "b": []string{"a", "b", "c"}}))
	assert.Equal(t, buf.args, []interface{}{1, "a", "b", "c"})
	assert.Equal(t, buf.String(), `(("a" = ?) OR ("b" IN (?, ?, ?)))`)

}

func TestDatasetSuite(t *testing.T) {
	suite.Run(t, new(datasetTest))
}
