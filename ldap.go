// Package ldap wraps ldap.v3 for xk6 extension
package ldap

import (
	"crypto/tls"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/ldap", new(Ldap))
}

type Ldap struct{}

func (l *Ldap) DialURL(addr string, opts map[string]bool) (*Conn, error) {
	var dialOpts []ldap.DialOpt
	insecureOptVal, ok := opts["insecureSkipTlsVerify"]
	if ok && insecureOptVal {
		conf := &tls.Config{InsecureSkipVerify: true} //nolint:all
		dialOpts = append(dialOpts, ldap.DialWithTLSConfig(conf))
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

func (c *Conn) Search(args map[string]any) (*ldap.SearchResult, error) {
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

	argsAttributes, ok := getOrDefault(args, "attributes", make([]any, 0)).([]any)
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

func (c *Conn) Modify(dn string, args map[string]map[string]string) error {
	modify := ldap.NewModifyRequest(dn, nil)

	for attribute, v := range args {
		switch v["operation"] {
		case "add":
			modify.Add(attribute, []string{v["value"]})
		case "replace":
			modify.Replace(attribute, []string{v["value"]})
		case "increment":
			modify.Increment(attribute, v["value"])
		case "delete":
			modify.Delete(attribute, make([]string, 0))
		default:
			return fmt.Errorf("Unsupported LDAP Modify operation for attribute %s", attribute)
		}
	}

	err := c.conn.Modify(modify)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func getOrDefault(m map[string]any, key string, defaultVal any) any {
	val, ok := m[key]
	if ok {
		return val
	}
	return defaultVal
}

func anyToStrSlice(a []any) []string {
	res := make([]string, len(a))
	for i, v := range a {
		res[i] = fmt.Sprint(v)
	}
	return res
}
