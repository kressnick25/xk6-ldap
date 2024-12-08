import { check } from 'k6'
import ldap from 'k6/x/ldap'


// Open the connection in the init section, called once per run
// It's a good idea not to open a connection per VU unless you specifically want to test that
const ldapUrl = 'ldap://localhost:1389'
console.log(`Dialing LDAP URL: ${ldapUrl}`)
let ldapConn = ldap.dialURL(ldapUrl)

let bindDn = 'cn=admin,dc=example,dc=org'
let bindPassword = 'adminpassword'
console.log(`Binding to LDAP with DN: ${bindDn}`)
ldapConn.bind(bindDn, bindPassword)

export default function () {
    let searchReq = {
        filter: '(cn=*)',
        baseDn: 'dc=example,dc=org',
        attributes: ['cn', 'sn', 'objectClass'], // [] for all attributes
        scope: 'WholeSubtree', // options: BaseObject, SingleLevel, WholeSubtree
        sizeLimit: 0, // 0 for unlimited
        timeLimit: 0, // (seconds). 0 for unlimited
        derefAliases: 0,
        typesOnly: false
    }

    let result = ldapConn.search(searchReq)
    console.log(`Search found ${result.entries.length} results`)
    check(result.entries, {
        'expected results': (r) => r.length === 3
    })

    let addAttributes = {
        'sn': ['Smith'],
        'objectClass': ['inetOrgPerson']
    }
    console.log('Running Add request')
    ldapConn.add('cn=test,dc=example,dc=org', addAttributes)


    // use default search attributes
    result = ldapConn.search({
        filter: '(cn=*)',
        baseDn: 'dc=example,dc=org'
    })
    console.log(`Search found ${result.entries.length} results`)
    check(result.entries, {
        'expected results': (r) => r.length === 4
    })

    console.log('Running Delete request')
    ldapConn.del('cn=test,dc=example,dc=org')

}

export function teardown() {
    ldapConn.close()
    console.log('LDAP connection closed')
}
