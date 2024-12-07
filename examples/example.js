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
    let filter = '(cn=*)'
    let baseDn = 'dc=example,dc=org'
    let attributes = ['cn', 'sn', 'objectClass'] // [] for all attributes
    let scope = 'WholeSubtree' // options: BaseObject, SingleLevel, WholeSubtree
    let sizeLimit = 0 // 0 for unlimited
    let timeLimit = 0 // (seconds). 0 for unlimited

    let searchReq = ldap.newSearchRequest(baseDn, scope, sizeLimit, timeLimit, filter, attributes) 

    let result = ldapConn.search(searchReq)
    console.log(`Search found ${result.entries.length} results`)

    let addRequest = ldap.newAddRequest('cn=test,dc=example,dc=org')
    addRequest.attribute('sn', ['Smith'])
    addRequest.attribute('objectClass', ['inetOrgPerson'])
    console.log('Running Add request')
    ldapConn.add(addRequest)


    result = ldapConn.search(searchReq)
    console.log(`Search found ${result.entries.length} results`)
    
    let delRequest = ldap.newDelRequest('cn=test,dc=example,dc=org')
    console.log('Running Delete request')
    ldapConn.del(addRequest)

}

export function teardown() {
    ldapConn.close()
    console.log('LDAP connection closed')
}
