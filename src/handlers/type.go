package handlers

import (
	"log"

	"gopkg.in/mgo.v2"
)

// Provider holds application wide variables
type Provider struct {
	log *log.Logger
	db  *mgo.Session
}

func NewProvider(log *log.Logger, db *mgo.Session) *Provider {
	return &Provider{
		log: log,
		db:  db,
	}
}

func (p *Provider) Logger() *log.Logger { return p.log }

func (p *Provider) DB() *mgo.Session { return p.db }
