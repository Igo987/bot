package db_test

import (
	"context"
	"github/Igo87/crypt/config"
	"github/Igo87/crypt/internal/db"
	"github/Igo87/crypt/models"
	"github/Igo87/crypt/pkg/api"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// add tests with pgxmock

func TestNewCurrencyRepository(t *testing.T) {
	t.Parallel()

	// Create a connection pool with a smaller max connections
	pool, err := pgxpool.New(context.Background(), config.Cfg.GetConnString())
	require.NoError(t, err)
	defer pool.Close()

	repo, err := db.NewCurrencyRepository(pool.Config().ConnString())
	require.NoError(t, err)
	assert.NotNil(t, repo)
}

func TestAddCurrencies(t *testing.T) {
	t.Parallel()
	cxt, cancel := context.WithTimeout(context.Background(), 120*time.Second)

	defer cancel()

	pgContainer, err := postgres.RunContainer(context.Background(),
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("../../", "migration", "1_init_db_up.sql")),
		postgres.WithDatabase("currency"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(cxt, connStr)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("../../migration/1_init_db_up.sql")
	if err != nil {
		t.Error(err)
	}

	stringArray := make([]string, len(data))

	for i, b := range data {
		stringArray[i] = string(b)
	}
	_, _, err = pgContainer.Exec(cxt, stringArray)
	if err != nil {
		t.Error(err)
	}

	defer pool.Close()

	info, err := api.FetchData(cxt)
	if err != nil {
		t.Fatal(err)
	}
	repo, err := db.NewCurrencyRepository(pool.Config().ConnString())
	if err != nil {
		t.Fatal(err)
	}
	err = repo.AddCurrencies(cxt, info)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSelectCurrenciesByToday(t *testing.T) {
	t.Parallel()
	pool, err := pgxpool.New(context.Background(), config.Cfg.GetConnString())
	if err != nil {
		t.Fatal(err)
	}

	defer pool.Close()

	query := `
		SELECT name, min(price) AS min_price, max(price) AS max_price,
		(MAX(price) - min(price)) / min(price) * 100 AS price_change_percent
		FROM currencies
		WHERE date(last_updated) = CURRENT_DATE
		GROUP BY 1
		LIMIT 2
	`

	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var extr models.Currencies

	for rows.Next() {
		var currency models.Extremes
		err := rows.Scan(&currency.Name, &currency.Min, &currency.Max, &currency.Percent)
		if err != nil {
			t.Fatalf("unable to scan row: %v", err)
		}

		extr = append(extr, currency)
	}
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}

	// check rows
	for _, c := range extr {
		assert.NotEmpty(t, c)
		assert.NotEmpty(t, c.Name)
		assert.NotEmpty(t, c.Min)
		assert.NotEmpty(t, c.Max)
		assert.NotEmpty(t, c.Percent)
	}

	assert.NotEmpty(t, extr)
	assert.Len(t, extr, 2)
	assert.Equal(t, "Bitcoin", extr[0].Name)
	assert.Equal(t, "Ethereum", extr[1].Name)
	assert.NoError(t, rows.Err())
}

func TestSelectCurrenciesByLastDay(t *testing.T) {
	t.Parallel()
	pool, err := pgxpool.New(context.Background(), config.Cfg.GetConnString())
	if err != nil {
		t.Fatal(err)
	}

	defer pool.Close()
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	// check date = yesterday
	rows, err := pool.Query(context.Background(),
		`select name,min(price),max(price),(MAX(price) - MIN(price)) / MIN(price) * 100 AS price_change_percent
	from currencies WHERE date(last_updated) = $1 group by 1 LIMIT 2`, yesterday)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var extr models.Currencies

	for rows.Next() {
		var currency models.Extremes
		err := rows.Scan(&currency.Name, &currency.Min, &currency.Max, &currency.Percent)
		if err != nil {
			t.Fatalf("unable to scan row: %v", err)
		}

		extr = append(extr, currency)
	}

	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	// check rows
	for _, c := range extr {
		assert.NotEmpty(t, c)
		assert.NotEmpty(t, c.Name)
		assert.NotEmpty(t, c.Min)
		assert.NotEmpty(t, c.Max)
	}
	// test with assert
	assert.NotEmpty(t, extr)
	assert.Equal(t, "Bitcoin", extr[0].Name)
	assert.Equal(t, "Ethereum", extr[1].Name)
	require.NoError(t, rows.Err())

	t.Log("success")
}
