package pq_array_test

import (
	"database/sql"
	_ "github.com/lib/pq"
	. "github.com/smagch/pq_array"
	"testing"
)

var (
	db *sql.DB
)

func drop(t *testing.T) {
	_, err := db.Exec(`DROP TABLE pq_array_test_table`)
	if err != nil {
		t.Fatal(err)
	}
}

func openDB(t *testing.T) {
	var err error
	if db == nil {
		db, err = sql.Open("postgres", "postgres://:@localhost:5432/test?sslmode=disable")
		if err != nil {
			t.Fatal(err)
		}
	}
	_, err = db.Exec(`
		DROP TABLE IF EXISTS pq_array_test_table;
		CREATE TABLE pq_array_test_table (nums integer[])`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNULLArray(t *testing.T) {
	openDB(t)
	defer drop(t)
	var arr IntArray
	err := db.QueryRow(`
		INSERT INTO pq_array_test_table (nums) VALUES ($1) RETURNING nums
	`, arr).Scan(&arr)
	if err != nil {
		t.Fatal(err)
	}
	if arr != nil {
		t.Fatal("array should be nil")
	}
}

func TestParseArray(t *testing.T) {
	openDB(t)
	defer drop(t)
	testCases := []IntArray{
		IntArray{1, 2, 3, 4, 5},
		IntArray{100, 200000, 300, 4, 1, -300},
		IntArray{-100, -100, -100},
		IntArray{0, -345, 3, 4, 5},
		IntArray{},
	}

	for testNum, tc := range testCases {
		var arr IntArray
		err := db.QueryRow(`INSERT INTO pq_array_test_table (nums) VALUES ($1) RETURNING nums`,
			tc).Scan(&arr)
		if err != nil {
			t.Fatal(err)
		}
		if len(arr) != len(tc) {
			t.Fatal("Unexpected: ", arr)
		}
		for i := 0; i < len(arr); i++ {
			if tc[i] != arr[i] {
				t.Errorf("Unexpected value %d in index %d.", arr[i], i)
				t.Fatalf("%d:Expected: %v", testNum, tc)
			}
		}
	}
}
