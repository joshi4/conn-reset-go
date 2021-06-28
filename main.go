package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	db.SetConnMaxIdleTime(1 * time.Hour)
	defer db.Close()

	ctx := context.Background()
	GetCourses(ctx, db)
	fmt.Println("initial calls ended")
	// initial queries end.
	for {
		<-time.After(50 * time.Minute)
		GetCourses(ctx, db)
	}
}

func GetCourses(ctx context.Context, db *sql.DB) {
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			getCourses(ctx, num, db)
		}(i)
	}
	wg.Wait()
}

func getCourses(ctx context.Context, num int, db *sql.DB) {
	var numRows int
	err := db.QueryRowContext(ctx, "select count(*) from courses").Scan(&numRows)
	if err != nil {
		fmt.Printf("number of rows: %d", numRows)
		return
	}
	fmt.Printf("err: %v", err)
}
