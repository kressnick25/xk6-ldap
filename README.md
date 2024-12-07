# k6 LDAP Extension

This is a [k6](https://k6.io) extension that enables LDAP operations within k6 test scripts. It provides a wrapper around the `ldap.v3` package, allowing you to perform LDAP operations such as binding, searching, adding, and deleting entries during load testing.

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

export default function () {
  // Create LDAP connection
  // Note that if you do not want to create a connection per VU, 
  // move this to the init context
  const conn = ldap.dialURL('ldap://your-ldap-server:389');

  try {
    // Bind to LDAP server
    conn.bind('cn=admin,dc=example,dc=com', 'admin_password');

    // Perform a search
    const searchRequest = ldap.newSearchRequest(
      'dc=example,dc=com',      // Base DN
      'WholeSubtree',           // Scope
      0,                        // Size Limit (0 for no limit)
      0,                        // Time Limit (0 for no limit)
      '(objectClass=person)',   // Filter
      ['cn', 'mail']           // Attributes to retrieve
    );

    const result = conn.search(searchRequest);
    console.log(result.entries);
  } finally {
    // Always close the connection
    conn.close();
  }
}
```

See also `examples/example.js`

## API Reference

### Ldap

#### `dialURL(addr: string): Conn`
Establishes a new connection to an LDAP server using the provided URL.

#### `escapeFilter(filter: string): string`
Escapes special characters in LDAP filter strings.

#### `newAddRequest(dn: string): AddRequest`
Creates a new request to add an entry to the LDAP directory.

#### `newDelRequest(dn: string): DelRequest`
Creates a new request to delete an entry from the LDAP directory.

#### `newSearchRequest(baseDn: string, scope: "BaseObject" | "SingleLevel" | "WholeSubtree", sizeLimit: number, timeLimit: number, filter: string, attributes: string[]): SearchRequest`
Creates a new search request with the specified parameters.

### Conn

#### `add(addRequest: AddRequest): void`
Adds a new entry to the LDAP directory.

#### `del(delRequest: DelRequest): void`
Deletes an entry from the LDAP directory.

#### `bind(username: string, password: string): void`
Authenticates with the LDAP server.

#### `search(searchRequest: SearchRequest): SearchResult`
Performs a search operation. Returns a SearchResult containing matched entries.

#### `close(): void`
Closes the LDAP connection.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open-source. Please ensure you check the license terms before using it.

## See Also

- [k6 Documentation](https://k6.io/docs/)
- [LDAP v3 Package Documentation](https://pkg.go.dev/gopkg.in/ldap.v3)

