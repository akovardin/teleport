package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go.uber.org/zap"
)

type Config struct {
	Driver  string
	Connect string
	Debug   bool
}

func NewDatabase(c Config, log *zap.SugaredLogger) *gorm.DB {
	db, err := gorm.Open(c.Driver, c.Connect)
	if err != nil {
		log.Fatalw("failed to connect database", "err", err)
	}

	db.LogMode(c.Debug)
	db.Debug()

	return db
}

type Sorting struct {
	Sort  string
	Order string
}

type Paginating struct {
	Start int
	End   int
}

type Param struct {
	Field string
	Value interface{}
}

type Condition struct {
	Params     []Param
	Sorting    Sorting
	Pagination Paginating
	Joins      []string
	Preload    []string
}

func NewCondition(sort, order string) Condition {
	condition := Condition{
		Params: []Param{},
		Sorting: Sorting{
			Sort:  sort,
			Order: order,
		},
	}

	return condition
}

func Find(q *gorm.DB, condition Condition) *gorm.DB {
	for _, join := range condition.Joins {
		q = q.Joins(join)
	}

	for _, preload := range condition.Preload {
		q = q.Preload(preload)
	}

	for _, param := range condition.Params {
		q = q.Where(param.Field+" = ?", param.Value)
	}

	sort := condition.Sorting.Sort
	order := condition.Sorting.Order

	if sort != "" && order != "" {
		q = q.Order(gorm.Expr(sort + " " + order))
	}

	start := condition.Pagination.Start
	end := condition.Pagination.End

	if end > start {
		q = q.Offset(start).Limit(end - start)
	}

	return q
}
