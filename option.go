package authority

import (
	"database/sql"
	"gorm.io/gorm"
)

type Options struct {
	prefix  string
	db      *gorm.DB
	migrate sql.NullBool
}

type Option func(*Options)

func newOptions(opt ...Option) *Options {
	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	if !opts.migrate.Valid {
		opts.migrate = sql.NullBool{Valid: true, Bool: false}
	}

	return &opts
}

func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.prefix = prefix
	}
}

func WithDB(db *gorm.DB) Option {
	return func(o *Options) {
		o.db = db
	}
}

func WithMigrate(migrate bool) Option {
	return func(o *Options) {
		o.migrate = sql.NullBool{Valid: true, Bool: migrate}
	}
}
