package eyeshade

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/brave-intl/bat-go/datastore/grantserver"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

var (
	// ErrLimitRequired signifies that the query requires a limit
	ErrLimitRequired = errors.New("query requires a limit")
	// ErrLimitReached signifies that the limit has been exceeded
	ErrLimitReached = errors.New("query limit reached")
	// ErrNeedTimeConstraints needs both a start and an end time
	ErrNeedTimeConstraints = errors.New("need both start and end date")
)

// AccountEarnings holds results from querying account earnings
type AccountEarnings struct {
	Channel   string          `json:"channel" db:"channel"`
	Earnings  decimal.Decimal `json:"earnings" db:"earnings"`
	AccountID string          `json:"account_id" db:"account_id"`
}

// AccountSettlementEarnings holds results from querying account earnings
type AccountSettlementEarnings struct {
	Channel   string          `json:"channel" db:"channel"`
	Paid      decimal.Decimal `json:"paid" db:"paid"`
	AccountID string          `json:"account_id" db:"account_id"`
}

// AccountEarningsOptions receives all options pertaining to account earnings calculations
type AccountEarningsOptions struct {
	Type      string
	Ascending bool
	Limit     int
}

// AccountSettlementEarningsOptions receives all options pertaining to account settlement earnings calculations
type AccountSettlementEarningsOptions struct {
	Type      string
	Ascending bool
	Limit     int
	StartDate *time.Time
	UntilDate *time.Time
}

// Datastore holds methods for interacting with database
type Datastore interface {
	grantserver.Datastore
	GetAccountEarnings(
		ctx context.Context,
		options AccountEarningsOptions,
	) (*[]AccountEarnings, error)
	GetAccountSettlementEarnings(
		ctx context.Context,
		options AccountSettlementEarningsOptions,
	) (*[]AccountSettlementEarnings, error)
}

// Postgres is a Datastore wrapper around a postgres database
type Postgres struct {
	grantserver.Postgres
}

// NewFromConnection creates new datastores from connection objects
func NewFromConnection(pg *grantserver.Postgres, name string) Datastore {
	return &DatastoreWithPrometheus{
		base:         &Postgres{*pg},
		instanceName: name,
	}
}

// NewDB creates a new Postgres Datastore
func NewDB(
	databaseURL string,
	performMigration bool,
	name string,
	dbStatsPrefix ...string,
) (Datastore, error) {
	pg, err := grantserver.NewPostgres(
		databaseURL,
		performMigration,
		"eyeshade", // to follow the migration track for eyeshade
		dbStatsPrefix...,
	)
	if err != nil {
		return nil, err
	}
	return NewFromConnection(pg, name), nil
}

// NewConnections creates postgres connections
func NewConnections() (Datastore, Datastore, error) {
	var eyeshadeRoPg Datastore
	eyeshadePg, err := NewDB(
		os.Getenv("EYESHADE_DB_URL"),
		true,
		"eyeshade_datastore",
		"eyeshade_db",
	)
	if err != nil {
		sentry.CaptureException(err)
		log.Panic().Err(err).Msg("Must be able to init postgres connection to start")
	}

	roDB := os.Getenv("EYESHADE_DB_RO_URL")
	if len(roDB) > 0 {
		eyeshadeRoPg, err = NewDB(
			roDB,
			false,
			"eyeshade_ro_datastore",
			"eyeshade_read_only_db",
		)
		if err != nil {
			sentry.CaptureException(err)
			log.Error().Err(err).Msg("Could not start reader postgres connection")
		}
	}
	return eyeshadePg, eyeshadeRoPg, err
}

// GetAccountEarnings gets the account earnings for a subset of ids
func (pg Postgres) GetAccountEarnings(
	ctx context.Context,
	options AccountEarningsOptions,
) (*[]AccountEarnings, error) {
	order := "desc"
	if options.Ascending {
		order = "asc"
	}
	limit := options.Limit
	if limit < 1 {
		return nil, ErrLimitRequired
	} else if limit > 1000 {
		return nil, ErrLimitReached
	}
	// remove the `s`
	txType := options.Type[:len(options.Type)-1]
	statement := fmt.Sprintf(`
select
	channel,
	coalesce(sum(amount), 0.0) as earnings,
	account_id
from account_transactions
where account_type = 'owner'
	and transaction_type = $1
group by (account_id, channel)
order by earnings %s
limit $2`, order)
	earnings := []AccountEarnings{}
	err := pg.RawDB().SelectContext(
		ctx,
		&earnings,
		statement,
		txType,
		limit,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &earnings, nil
}

// GetAccountSettlementEarnings gets the account settlement earnings for a subset of ids
func (pg Postgres) GetAccountSettlementEarnings(
	ctx context.Context,
	options AccountSettlementEarningsOptions,
) (*[]AccountSettlementEarnings, error) {
	order := "desc"
	if options.Ascending {
		order = "asc"
	}
	limit := options.Limit
	if limit < 1 {
		return nil, ErrLimitRequired
	} else if limit > 1000 {
		return nil, ErrLimitReached
	}
	txType := options.Type[:len(options.Type)-1]
	txType = fmt.Sprintf(`%s_settlement`, txType)
	timeConstraints := ""
	constraints := []interface{}{
		txType,
		limit,
	}
	if options.UntilDate != nil {
		// if end date exists, start date must also
		if options.StartDate == nil {
			return nil, ErrNeedTimeConstraints
		}
		constraints = append(
			constraints,
			*options.StartDate,
			*options.UntilDate,
		)
		timeConstraints = `
and created_at >= $3
and created_at < $4`
	}
	statement := fmt.Sprintf(`
select
	channel,
	coalesce(sum(-amount), 0.0) as paid,
	account_id
from account_transactions
where
		account_type = 'owner'
and transaction_type = $1%s
group by (account_id, channel)
order by paid %s
limit $2`, timeConstraints, order)
	earnings := []AccountSettlementEarnings{}

	err := pg.RawDB().SelectContext(
		ctx,
		&earnings,
		statement,
		constraints...,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &earnings, nil
}
