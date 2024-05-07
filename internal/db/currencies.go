package db

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github/Igo87/crypt/models"
	"github/Igo87/crypt/pkg/api"
)

var ErrGetCurrencies = errors.New("failed to get currencies from the repository: connection not established")
var ErrScanRow = errors.New("unable to scan row")
var Err = errors.New("unable to execute query")
var ErrIterateRows = errors.New("unable to iterate rows")
var ErrGetConnection = errors.New("unable to get connection")
var ErrFetchFromAPI = errors.New("failed to fetch data from the remote API")
var ErrFailedToFetchData = errors.New("failed to fetch data from the repository")
var ErrFailedToInsertData = errors.New("failed to insert data")
var ErrFailedConnectToDB = errors.New("failed to connect to the database")
var ErrFailedBeginTRX = errors.New("failed to begin transaction")
var ErrFailedCommitTRX = errors.New("failed to commit transaction")
var ErrFailedRollbackTRX = errors.New("failed to rollback transaction")
var ErrFailedToAddCurrencies = errors.New("failed to add currencies")
var ErrNilPointerRecevied = errors.New("nil pointer received")
var ErrCreateDB = errors.New("failed to create DB")

//go:generate mockgen -source=currencies.go -destination=mocks/currencies.go -package=mocks
type Currencies interface {
	InitDB(ctx context.Context, l slog.Logger) error
	AddCurrencies(ctx context.Context, data []models.Currencies) error
	Run(ctx context.Context, l slog.Logger) error
	SelectCurrenciesByLastDay(ctx context.Context) (models.Currencies, error)
	SelectCurrenciesByToday(ctx context.Context) (models.Currencies, error)
	Fetch(ctx context.Context, l slog.Logger) ([]models.Currencies, error)
	SelectCurrenciesByName(ctx context.Context, date string) (models.Currencies, error)
}

//go:generate mockgen -source=currencies.go -destination=mocks/currencies.go -package=mocks
type CurrencyRepository struct {
	conn *pgxpool.Pool
}

// Fetch fetches data from the API and returns it as a models.Crypto struct.
//
// ctx: The context for the fetch operation.
// l: The logger for logging any errors during the fetch operation.
// Returns a models.Crypto struct and an error.

func (c *CurrencyRepository) Fetch(ctx context.Context, log slog.Logger) (models.Crypto, error) {
	if c == nil || c.conn == nil {
		log.Error(ErrNilPointerRecevied.Error())
		return models.Crypto{}, ErrNilPointerRecevied
	}

	data, err := api.FetchData(ctx)

	if err != nil {
		return models.Crypto{}, ErrFetchFromAPI
	}
	return data, nil
}

// Run runs the CurrencyRepository in an infinite loop, fetching data and adding it to the repository.
//
// ctx: The context for the run operation.
// l: The logger for logging any errors during the run operation.
// Returns an error if there was a problem fetching or adding data.
func (c *CurrencyRepository) Run(ctx context.Context, log slog.Logger) error {
	for {
		data, err := c.Fetch(ctx, log)
		if err != nil {
			log.Error(ErrFailedToFetchData.Error())
			return ErrFailedToFetchData
		}
		err = c.AddCurrencies(ctx, data)
		if err != nil {
			log.Error(ErrFailedToInsertData.Error())
			return ErrFailedToInsertData
		}
	}
}

// NewСurrencyRepository creates a new instance of the СurrencyRepository struct.
//
// It takes a connString parameter of type string, which represents the connection string to the database.
// The function returns a pointer to a СurrencyRepository object and an error.
func NewCurrencyRepository(connString string) (*CurrencyRepository, error) {
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, ErrFailedConnectToDB
	}

	return &CurrencyRepository{conn: conn}, nil
}

// Close closes the connection to the database, if it exists.
func (c *CurrencyRepository) Close() error {
	if c.conn == nil {
		return nil
	}

	c.conn.Close()
	c.conn = nil
	return nil
}

func (c *CurrencyRepository) AddCurrencies(ctx context.Context, data models.Crypto) error {
	if c.conn == nil {
		return ErrGetCurrencies
	}

	transaction, err := c.conn.Begin(ctx)
	if err != nil {
		return ErrFailedBeginTRX
	}

	btc, eth := data.Data.Bitcoin, data.Data.Ethereum

	query := `INSERT INTO currencies (name, last_updated, price, percent_change_1h, percent_change_24h,
		percent_change_7d, percent_change_30d)
		VALUES ($1, $2, $3, $4, $5, $6, $7), ($8, $9, $10, $11, $12, $13, $14) returning id`

	btcValues := []interface{}{
		btc.Name, time.Now(), btc.Quote.Rub.Price,
		btc.Quote.Rub.PercentChange1H, btc.Quote.Rub.PercentChange24H,
		btc.Quote.Rub.PercentChange7D, btc.Quote.Rub.PercentChange30D,
	}

	ethValues := []interface{}{
		eth.Name, time.Now(), eth.Quote.Rub.Price,
		eth.Quote.Rub.PercentChange1H, eth.Quote.Rub.PercentChange24H,
		eth.Quote.Rub.PercentChange7D, eth.Quote.Rub.PercentChange30D,
	}

	_, err = transaction.Exec(ctx, query, append(btcValues, ethValues...)...)
	if err != nil {
		err = transaction.Rollback(ctx)
		if err != nil {
			return ErrFailedRollbackTRX
		}
		return ErrFailedToAddCurrencies
	}

	err = transaction.Commit(ctx)
	if err != nil {
		err = transaction.Rollback(ctx)
		if err != nil {
			return ErrFailedRollbackTRX
		}
		return ErrFailedCommitTRX
	}

	return nil
	// if c.conn == nil {
	// 	return fmt.Errorf("failed to add currencies to the repository: connection not established")
	// }

	// go func() {
	// 	<-ctx.Done()
	// 	c.Close()
	// }()

	// tx, err := c.conn.Begin(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to begin transaction: %w", err)
	// }
	// /* trunk-ignore(golangci-lint/errcheck) */
	// defer tx.Rollback(ctx)

	// btc, eth := data.Data.Bitcoin, data.Data.Ethereum
	// fmt.Println(btc, eth)

	// query := `INSERT INTO currencies (name, last_updated, price, percent_change_1h, percent_change_24h, percent_change_7d, percent_change_30d)
	// 		  VALUES ($1, $2, $3, $4, $5, $6, $7), ($8, $9, $10, $11, $12, $13, $14)`

	// btcValues := []interface{}{
	// 	btc.Name, time.Now(), btc.Quote.Rub.Price,
	// 	btc.Quote.Rub.PercentChange1H, btc.Quote.Rub.PercentChange24H,
	// 	btc.Quote.Rub.PercentChange7D, btc.Quote.Rub.PercentChange30D,
	// }

	// ethValues := []interface{}{
	// 	eth.Name, time.Now(), eth.Quote.Rub.Price,
	// 	eth.Quote.Rub.PercentChange1H, eth.Quote.Rub.PercentChange24H,
	// 	eth.Quote.Rub.PercentChange7D, eth.Quote.Rub.PercentChange30D,
	// }

	// _, err = tx.Exec(ctx, query, append(btcValues, ethValues...)...)
	// if err != nil {
	// 	return fmt.Errorf("failed to add currencies to the repository: %w", err)
	// }

	// err = tx.Commit(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to commit transaction: %w", err)
	// }

	// return nil
}

// Максимальное и минимальное значение цены валюты за сегодня.
func (c *CurrencyRepository) SelectCurrenciesByToday(ctx context.Context) (models.Currencies, error) {

	if c.conn == nil {
		return nil, ErrGetConnection
	}

	query := `select name,min(price),max(price),(MAX(price) - MIN(price)) / MIN(price)
	* 100 AS price_change_percent from currencies WHERE date(last_updated) = CURRENT_DATE group by 1 LIMIT 2`
	rows, err := c.conn.Query(ctx, query)

	if err != nil {
		return nil, Err
	}

	defer rows.Close()

	extr := make(models.Currencies, 0)
	// read all rows in the result
	for rows.Next() {
		var ext models.Extremes
		err := rows.Scan(&ext.Name, &ext.Min, &ext.Max, &ext.Percent)
		if err != nil {
			rows.Close()
			return nil, ErrScanRow
		}
		extr = append(extr, ext)
	}
	return extr, nil
}

// Максимальное и минимальное значение цены валюты за вчера.
func (c *CurrencyRepository) SelectCurrenciesByLastDay(ctx context.Context) (models.Currencies, error) {
	if c.conn == nil {

		return nil, ErrGetCurrencies
	}
	query := `
		SELECT name, min(price), max(price),
		(MAX(price) - MIN(price)) / MIN(price) * 100 AS price_change_percent
		FROM currencies
		WHERE date(last_updated) = CURRENT_DATE - INTERVAL '1 day'
		GROUP BY 1
		LIMIT 2
	`

	rows, err := c.conn.Query(ctx, query)
	if err != nil {
		return nil, ErrGetConnection
	}
	defer rows.Close()

	extr := make(models.Currencies, 0)
	for rows.Next() {
		var ext models.Extremes
		if err := rows.Scan(&ext.Name, &ext.Min, &ext.Max, &ext.Percent); err != nil {
			rows.Close()

			return nil, ErrScanRow
		}
		extr = append(extr, ext)
	}

	if err := rows.Err(); err != nil {
		return nil, ErrIterateRows
	}

	return extr, nil
}

// request from the database by currency name for today.
func (c *CurrencyRepository) SelectCurrenciesByName(ctx context.Context, name string) (models.Currencies, error) {
	if c.conn == nil {
		return nil, ErrGetCurrencies
	}

	query := `select name,min(price),max(price),(MAX(price) - MIN(price)) / MIN(price) * 100 AS
	price_change_percent from currencies WHERE date(last_updated) = CURRENT_DATE AND name = $1 GROUP BY name LIMIT 1;`
	row := c.conn.QueryRow(ctx, query, name)

	var ext models.Extremes
	err := row.Scan(&ext.Name, &ext.Min, &ext.Max, &ext.Percent)
	if err != nil {
		return nil, ErrScanRow
	}

	extr := make(models.Currencies, 1)
	extr[0] = ext

	return extr, nil
}
