// Package ldap wraps ldap.v3 for xk6 extension
package ldap

import (
	"go.k6.io/k6/js/modules"
	"gopkg.in/ldap.v3"
)

func init() {
	modules.Register("k6/x/ldap", new(Ldap))
}

type Ldap struct {
}

func (l *Ldap) DialURL(addr string) (*Conn, error) {
	c, err := ldap.DialURL(addr)
	if err != nil {
		return nil, err
	}
	newConn := Conn{conn: c}
	return &newConn, nil
}

func (l *Ldap) EscapeFilter(filter string) string {
	return ldap.EscapeFilter(filter)
}

func (l *Ldap) NewAddRequest(dn string) *ldap.AddRequest {
	return ldap.NewAddRequest(dn, []ldap.Control{})
}

func (l *Ldap) NewDelRequest(dn string) *ldap.DelRequest {
	return ldap.NewDelRequest(dn, []ldap.Control{})
}

func (l *Ldap) NewSearchRequest(
	baseDn string,
	scope string,
	sizeLimit int,
	timeLimit int,
	filter string,
	attributes []string,
) *ldap.SearchRequest {
	var _scope int
	switch scope {
	case "BaseObject":
		_scope = ldap.ScopeBaseObject
	case "SingleLevel":
		_scope = ldap.ScopeSingleLevel
	default:
		_scope = ldap.ScopeWholeSubtree
	}

	// defaults
	derefAliases := 0
	typesOnly := false
	control := []ldap.Control{}

	lsr := ldap.NewSearchRequest(baseDn, _scope, derefAliases,
		sizeLimit, timeLimit, typesOnly, filter, attributes, control)

	return lsr
}

type Conn struct {
	conn *ldap.Conn
}

func (c *Conn) Add(addRequest *ldap.AddRequest) error {
	return c.conn.Add(addRequest)
}

func (c *Conn) Del(delRequest *ldap.DelRequest) error {
	return c.conn.Del(delRequest)
}

func (c *Conn) Bind(username string, password string) error {
	return c.conn.Bind(username, password)
}

func (c *Conn) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
	result, err := c.conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (c *Conn) Close() {
	c.conn.Close()
}
