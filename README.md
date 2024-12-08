# xk6-ldap

A [k6](https://k6.io) extension that enables LDAP operations within k6 scripts. It is mostly a wrapper around the golang `ldap.v3` package.

## Features

LDAP Operations
- BIND
- Search
- Add entry
- Delete entry

Utils
- Escape search filter

### Not implemented (yet)
- LDAP Modify
- operation controls

## Installation

To use this extension, you'll need to build k6 with the extension enabled. Follow these steps:

1. Clone this repository
2. Build k6 with the extension using xk6:
```bash
xk6 build --with github.com/kressnick25/xk6-ldap
```

## Usage

```javascript
import ldap from 'k6/x/ldap';

// Create LDAP connection using an LDAPURL
// The connection will be established on the first command call.
const conn = ldap.dialURL('ldaps://your-ldap-server:636');

export default function () {

  try {
    // Bind to LDAP server
    conn.bind('cn=admin,dc=example,dc=com', 'admin_password');

    // Perform a search
    const searchRequest = {
        filter: '(objectClass=person)', // Search Filter
        baseDn: 'dc=example,dc=org', // Base DN
        attributes: ['cn', 'mail'], // [] for all attributes
        scope: 'WholeSubtree', // options: BaseObject, SingleLevel, WholeSubtree
        sizeLimit: 0, // 0 for no limit
        timeLimit: 0, // (seconds) 0 for not limit
        derefAliases: 0,
        typesOnly: false
    }

    const result = conn.search(searchRequest);
    console.log(result.entries);
  } finally {
    // Always close the connection
    conn.close();
  }
}
```

See also `examples/example.js`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open-source. Please ensure you check the license terms before using it.

## See Also

- [k6 Documentation](https://k6.io/docs/)
- [LDAP v3 Package Documentation](https://pkg.go.dev/gopkg.in/ldap.v3)


## API Reference

### Module Import
```javascript
import ldap from 'k6/x/ldap'
```

### LDAP Module Methods

#### dialURL(address: string)
Establishes a connection to an LDAP server.

**Parameters:**
- `address` (string): LDAP URL in the format `ldap://host:port`

**Returns:**
- Connection object

#### escapeFilter(filter: string)
Escapes special characters in LDAP filter strings to prevent injection.

**Parameters:**
- `filter` (string): LDAP filter string to escape

**Returns:**
- Escaped filter string

### Connection Methods

#### bind(username: string, password: string)
Authenticates the connection with the LDAP server.

**Parameters:**
- `username` (string): DN of the user to authenticate as
- `password` (string): Password for authentication

#### search(options: object)
Performs an LDAP search operation.

**Parameters:**
- `options` (object):
  - `filter` (string, optional): LDAP search filter. Default: "*"
  - `baseDn` (string, optional): Base DN for search. Default: ""
  - `attributes` (string[], optional): Attributes to return. Default: []
  - `scope` (string, optional): Search scope. Default: "WholeSubtree"
    - Valid values: "BaseObject", "SingleLevel", "WholeSubtree"
  - `sizeLimit` (number, optional): Maximum entries to return. Default: 0 (unlimited)
  - `timeLimit` (number, optional): Search time limit in seconds. Default: 0 (unlimited)
  - `derefAliases` (number, optional): Alias dereferencing option. Default: 0
  - `typesOnly` (boolean, optional): Return attribute names only. Default: false

**Returns:**
- Search result object containing matching entries

#### add(dn: string, attributes: string)
Adds a new entry to the LDAP directory.

**Parameters:**
- `dn` (string): Distinguished Name for the new entry
- `attributes` (object): Map of attribute names to arrays of values

#### del(dn: string)
Deletes an entry from the LDAP directory.

**Parameters:**
- `dn` (string): Distinguished Name of the entry to delete

#### close()
Closes the LDAP connection and releases resources.