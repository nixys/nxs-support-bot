package redmine

import (
	rdmn "github.com/nixys/nxs-go-redmine/v4"
)

// Redmine it is a module context structure
type Redmine struct {
	c    rdmn.Context
	host string
}

// Init initiates context to interact with Redmine
func Init(host, key string) Redmine {

	var r Redmine

	r.c.SetEndpoint(host)
	r.c.SetAPIKey(key)

	r.host = host

	return r
}
