import { check } from 'k6'
import ldap from 'k6/x/ldap'

const ldapUrl = 'ldaps://localhost:1636'

export default function () {
    let ldapConn
    console.log(`Dialing LDAP URL: ${ldapUrl}`)
    try {
        // Should fail, container started with self-signed cert
        ldapConn = ldap.dialURL(ldapUrl, {
            insecureSkipTlsVerify: false
        })

        throw new Error(tlsErrMsg)
    } catch (err) {
        if (err.message.includes("tls: failed to verify certificate")) {
            // expected
            console.debug("This error is expected: " + err)
        } else {
            console.error("TLS verification should have thrown, but didn't")
        }

        // retry, ignoring certificate errors
        ldapConn = ldap.dialURL(ldapUrl, {
            // DO NOT EVER DO THIS IN PRODUCTION
            insecureSkipTlsVerify: true
        })
    }
    try {
        let bindDn = 'cn=admin,dc=example,dc=org'
        let bindPassword = 'adminpassword'
        console.log(`Binding to LDAP with DN: ${bindDn}`)
        ldapConn.bind(bindDn, bindPassword)
        
        test(ldapConn)
    } finally {
        ldapConn.close()
        console.log('LDAP connection closed')
    }
}

function test(ldapConn) {
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
        'expected results': (r) => r.length > 1
    })

    let addAttributes = {
        'sn': ['Smith'],
        'objectClass': ['inetOrgPerson']
    }
    console.log('Running Add request')
    ldapConn.add('cn=test-tls,dc=example,dc=org', addAttributes)


    // use default search attributes
    result = ldapConn.search({
        filter: '(cn=*)',
        baseDn: 'dc=example,dc=org'
    })
    console.log(`Search found ${result.entries.length} results`)
    check(result.entries, {
        'expected results': (r) => r.length > 1
    })

    console.log('Running Delete request')
    ldapConn.del('cn=test-tls,dc=example,dc=org')

}

