// Package ldap wraps ldap.v3 for xk6 extension
package ldap

import (
	"fmt"
    "crypto/tls"

	"go.k6.io/k6/js/modules"
	"github.com/go-ldap/ldap/v3"
)

func init() {
	modules.Register("k6/x/ldap", new(Ldap))
}

type Ldap struct{}

func (l *Ldap) DialURL(addr string, opts map[string]bool) (*Conn, error) {
    var dialOpts []ldap.DialOpt
    insecureOptVal, ok := opts["insecureSkipTlsVerify"]
    if ok && insecureOptVal {
        dialOpts = append(dialOpts, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
    } 

    c, err := ldap.DialURL(addr, dialOpts...)
	if err != nil {
		return nil, err
	}
	newConn := Conn{conn: c}
	return &newConn, nil
}

func (l *Ldap) EscapeFilter(filter string) string {
	return ldap.EscapeFilter(filter)
}

type Conn struct {
	conn *ldap.Conn
}

func (c *Conn) Add(dn string, attributes map[string][]string) error {
	addReq := ldap.NewAddRequest(dn, []ldap.Control{})
	for k, v := range attributes {
		addReq.Attribute(k, v)
	}
	return c.conn.Add(addReq)
}

func (c *Conn) Del(dn string) error {
	delReq := ldap.NewDelRequest(dn, []ldap.Control{})
	return c.conn.Del(delReq)
}

func (c *Conn) Bind(username string, password string) error {
	return c.conn.Bind(username, password)
}

func (c *Conn) Search(args map[string]interface{}) (*ldap.SearchResult, error) {
	var _scope int
	switch getOrDefault(args, "scope", "WholeSubtree") {
	case "BaseObject":
		_scope = ldap.ScopeBaseObject
	case "SingleLevel":
		_scope = ldap.ScopeSingleLevel
	default:
		_scope = ldap.ScopeWholeSubtree
	}

	// defaults
	control := []ldap.Control{}

	errorMsg := "Invalid search argument type:"

	filter, ok := getOrDefault(args, "filter", "*").(string)
	if !ok {
		return nil, fmt.Errorf("%s %s", errorMsg, "filter")
	}
	baseDn, ok := getOrDefault(args, "baseDn", "").(string)
	if !ok {
		return nil, fmt.Errorf("%s %s", errorMsg, "baseDn")
	}
	derefAliases, ok := getOrDefault(args, "derefAliases", int64(0)).(int64)
	if !ok {
		return nil, fmt.Errorf("%s %s", errorMsg, "derefAliases")
	}
	sizeLimit, ok := getOrDefault(args, "sizeLimit", int64(0)).(int64)
	if !ok {
		return nil, fmt.Errorf("%s %s", errorMsg, "sizeLimit")
	}
	timeLimit, ok := getOrDefault(args, "timeLimit", int64(0)).(int64)
	if !ok {
		return nil, fmt.Errorf("%s %s", errorMsg, "timeLimit")
	}
	typesOnly, ok := getOrDefault(args, "typesOnly", false).(bool)
	if !ok {
		return nil, fmt.Errorf("%s %s", errorMsg, "typesOnly")
	}

	argsAttributes, ok := getOrDefault(args, "attributes", make([]interface{}, 0)).([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s %s", errorMsg, "attributes")
	}
	attributes := make([]string, len(argsAttributes))
	for i, v := range argsAttributes {
		attributes[i] = fmt.Sprint(v)
	}

	searchRequest := ldap.NewSearchRequest(
		baseDn,
		_scope,
		int(derefAliases),
		int(sizeLimit),
		int(timeLimit),
		typesOnly,
		filter,
		attributes,
		control,
	)

	result, err := c.conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Conn) Close() {
	c.conn.Close()
}

func getOrDefault(m map[string]interface{}, key string, defaultVal interface{}) interface{} {
	val, ok := m[key]
	if ok {
		return val
	}
	return defaultVal
}
