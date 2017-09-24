package orm

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"go-support/log"
)

type User struct {
	Name        string `type:"VARCHAR(256)" constraint:"NOT NULL"`
	Age         int
	TestFloat32 float32
	TestFloat64 float64
	Male        bool
}

func TestNewDB(t *testing.T) {
	db, err := NewDB("mysql", "root:123456@tcp(localhost:3306)/uc")
	if err != nil {
		t.Fatal(err)
	}
	db.Logger(log.NewLogger(nil)).ShowSql(true)
	table := db.GetTable("user", User{})
	t.Log(table)
	err = table.CreateIfNotExists()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Create table success...")
	id, err := table.NewInsert(&User{
		Name:"zhangsan",
		Age:18,
		TestFloat32:13.5,
		TestFloat64:15.6,
		Male:true,
	}).Exec()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("user.id = ", id)
	where := NewWhere("name=?", "zhangsan")
	user, err := table.NewQuery().Where(where).GetOne()
	t.Log(user)
	affects, err := table.NewDelete().Where(where).Exec()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("affect rows = ", affects)
}
