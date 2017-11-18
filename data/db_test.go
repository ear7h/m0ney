package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

var _testQuotes []Quote

func must(e error) {
	if e != nil {
		panic(e)
	}
}

// initialize database
func init() {
	db, err := sql.Open("mysql", "root@(:3306)/?multiStatements=true")
	must(err)

	byt, err := ioutil.ReadFile("database_structure.sql")
	must(err)

	_, err = db.Exec(string(byt))
	must(err)
	db.Close()

	byt, err = ioutil.ReadFile("test_data.json")
	must(err)

	_testQuotes = []Quote{}
	err = json.Unmarshal(byt, &_testQuotes)
	must(err)
}

func TestNewMoneyDB(t *testing.T) {
	_, err := NewMoneyDB("root", "", "127.0.0.1", "3306", "money_test")
	if err != nil {
		panic(err)
	}
}

func TestMoneyDB_InsertQuotes(t *testing.T) {

	mdb, err := NewMoneyDB("root", "", "127.0.0.1", "3306", "money_test")
	if err != nil {
		panic(err)
	}

	err = mdb.InsertQuotes(_testQuotes)
	if err != nil {
		panic(err)
	}

}

func TestMoneyDB_Nightly(t *testing.T) {
	mdb, err := NewMoneyDB("root", "", "127.0.0.1", "3306", "money_test")
	if err != nil {
		panic(err)
	}

	err = mdb.InsertQuotes(_testQuotes)
	must(err)

	// move to partition
	mdb.Nightly()

	fmt.Println("")

	q := fmt.Sprintf(`SELECT updated_at FROM %s ORDER BY updated_at ASC`, mdb.lastPart)
	rows, err := mdb.Query(q)
	must(err)

	for _, v := range _testQuotes {
		ts := time.Time{}

		if !rows.Next() {
			panic(fmt.Errorf("no row for test quote %v", v))
		}
		err = rows.Scan(&ts)
		must(err)

		if ts.Format(SQL_TIME) != v.UpdatedAt.Format(SQL_TIME) {
			fmt.Println(ts.Format(SQL_TIME), v.UpdatedAt.Format(SQL_TIME))
			t.Fail()
			return
		}
	}
}

func TestMoneyDB_Partitions(t *testing.T) {
	mdb, err := NewMoneyDB("root", "", "127.0.0.1", "3306", "money_test")
	if err != nil {
		panic(err)
	}

	parts, err := mdb.Partitions()
	if err != nil {
		panic(err)
	}

	fmt.Println(parts)
}

func TestMoneyDB_Runs(t *testing.T) {
	mdb, err := NewMoneyDB("root", "", "127.0.0.1", "3306", "money_test")
	if err != nil {
		panic(err)
	}

	err = mdb.InsertQuotes(_testQuotes)
	must(err)

	mdb.Nightly()

	runs, err := mdb.Runs()
	must(err)

	// only one run from test sequence
	if len(runs) != 1 ||
		// start time should be same timestamp as first in test sequence
		runs[0].Start.Format(SQL_TIME) !=
			_testQuotes[0].UpdatedAt.Format(SQL_TIME) ||
		// end time should be same timestamp as last in test sequence
		runs[0].End.Format(SQL_TIME) !=
			_testQuotes[len(_testQuotes)-1].UpdatedAt.Format(SQL_TIME) {
		fmt.Println("runs not correct: ", runs)
		fmt.Println(runs[0].Start.Format(SQL_TIME) !=
			_testQuotes[0].UpdatedAt.Format(SQL_TIME))
		fmt.Println(runs[0].End.Format(SQL_TIME) !=
			_testQuotes[len(_testQuotes)-1].UpdatedAt.Format(SQL_TIME))
		t.Fail()
		return
	}
}

func TestMoneyDB_Run(t *testing.T) {
	mdb, err := NewMoneyDB("root", "", "127.0.0.1", "3306", "money_test")
	if err != nil {
		panic(err)
	}

	err = mdb.InsertQuotes(_testQuotes)
	must(err)

	mdb.Nightly()

	run, err := mdb.Run(1)
	must(err)

	for k, v := range run {
		if v.UpdatedAt.Format(SQL_TIME) != _testQuotes[k].UpdatedAt.Format(SQL_TIME) {
			fmt.Println("fail at ", k, ": ", v, _testQuotes[k])
			t.Fail()
			return
		}
	}

}
