package setup

// cockroachDB connection with pgx

import (
	"context"
	"fmt"

	// "log"
	"os"

	// "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	// "github.com/google/uuid"
	// "github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool
var err error

// database connection
func ConnectDB() *pgxpool.Pool {
	// Read in connection string
	db_url := os.Getenv("DB_URL")

	// fmt.Println(db_url)
	DB, err = pgxpool.New(context.Background(), db_url)
	// DB, err = pgx.Connect(context.Background(), db_url)
	//context background ()  returns a non-nil, empty Context.
	//It is typically used by the main function, initialization,and tests, and as the top-level Context for incoming requests.
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("database connected successfully")

	return DB
}
