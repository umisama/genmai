package genmai

import (
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type testModel struct {
	Id   int64
	Name string
	Addr string
}

func newTestDB(t *testing.T) *DB {
	db, err := New(&SQLite3Dialect{}, ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{
		`CREATE TABLE test_model (
			id INTEGER NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			addr TEXT NOT NULL
		);`,
		`INSERT INTO test_model (id, name, addr) VALUES (1, 'test1', 'addr1');`,
		`INSERT INTO test_model (id, name, addr) VALUES (2, 'test2', 'addr2');`,
		`INSERT INTO test_model (id, name, addr) VALUES (3, 'test3', 'addr3');`,
		`INSERT INTO test_model (id, name, addr) VALUES (4, 'other', 'addr4');`,
		`INSERT INTO test_model (id, name, addr) VALUES (5, 'other', 'addr5');`,
		`INSERT INTO test_model (id, name, addr) VALUES (6, 'other1', 'addr6');`,
		`INSERT INTO test_model (id, name, addr) VALUES (7, 'other2', 'addr7');`,
	} {
		if _, err := db.db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	return db
}

func Test_Select(t *testing.T) {
	// SELECT * FROM test_model;
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{1, "test1", "addr1"},
			{2, "test2", "addr2"},
			{3, "test3", "addr3"},
			{4, "other", "addr4"},
			{5, "other", "addr5"},
			{6, "other1", "addr6"},
			{7, "other2", "addr7"},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()

	// SELECT * FROM test_model WHERE "id" = 1;
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual, db.Where("id", "=", 1)); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{1, "test1", "addr1"},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()

	// SELECT * FROM test_model ORDER BY "id" DESC LIMIT 2;
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual, db.Limit(2).OrderBy("id", DESC)); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{7, "other2", "addr7"}, {6, "other1", "addr6"},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()

	// SELECT * FROM test_model LIMIT 2 OFFSET 3;
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual, db.Limit(2).Offset(3)); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{4, "other", "addr4"}, {5, "other", "addr5"},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()

	// SELECT * FROM test_model WHERE "id" = 1 OR ("id" = 5 AND "name" = "other");
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual, db.Where("id", "=", 1).Or(db.Where("id", "=", 5).And("name", "=", "other"))); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{1, "test1", "addr1"}, {5, "other", "addr5"},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()

	// SELECT * FROM test_model WHERE "id" = 1 AND "name" = "test1";
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual, db.Where("id", "=", 1).And("name", "=", "test1")); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{1, "test1", "addr1"},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()

	// SELECT * FROM test_model WHERE "id" IN (2, 3);
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual, db.Where("id").In(2, 3)); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{2, "test2", "addr2"}, {3, "test3", "addr3"},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()

	// SELECT * FROM test_model WHERE "name" LIKE "%3";
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual, db.Where("name").Like("%3")); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{3, "test3", "addr3"},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()

	// SELECT * FROM test_model WHERE "name" = "other" LIMIT 1;
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual, db.Where("name", "=", "other").Limit(1).OrderBy("id", ASC)); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{4, "other", "addr4"},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()

	// SELECT * FROM test_model WHERE "name" = "other" LIMIT 1 OFFSET 1;
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual, db.Where("name", "=", "other").Limit(1).OrderBy("id", ASC).Offset(1)); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{5, "other", "addr5"},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()

	// SELECT * FROM test_model WHERE "id" BETWEEN 3 AND 5;
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual, db.Where("id").Between(3, 5)); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{3, "test3", "addr3"}, {4, "other", "addr4"}, {5, "other", "addr5"},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()

	// SELECT "id" FROM test_model;
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual, "id"); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{1, "", ""},
			{2, "", ""},
			{3, "", ""},
			{4, "", ""},
			{5, "", ""},
			{6, "", ""},
			{7, "", ""},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()

	// SELECT "name", "addr" FROM test_model;
	func() {
		db := newTestDB(t)
		defer db.Close()
		var actual []testModel
		if err := db.Select(&actual, []string{"name", "addr"}); err != nil {
			t.Fatal(err)
		}
		expected := []testModel{
			{0, "test1", "addr1"},
			{0, "test2", "addr2"},
			{0, "test3", "addr3"},
			{0, "other", "addr4"},
			{0, "other", "addr5"},
			{0, "other1", "addr6"},
			{0, "other2", "addr7"},
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Expect %v, but %v", expected, actual)
		}
	}()
}
