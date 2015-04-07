package smtp

import (
	"fmt"
	"log"
	"net"
	"regexp"
)

type MailAddress struct {
	Local  string
	Domain string
}

func (m *MailAddress) String() string {
	return fmt.Sprintf("%s@%s", m.Local, m.Domain)
}

// Validate the email adress
/*
   RFC 5321

   address-literal  = "[" ( IPv4-address-literal /
                    IPv6-address-literal /
                    General-address-literal ) "]"
                    ; See Section 4.1.3

   Mailbox        = Local-part "@" ( Domain / address-literal )

   Local-part     = Dot-string / Quoted-string
                  ; MAY be case-sensitive


   Dot-string     = Atom *("."  Atom)

   Atom           = 1*atext

   Quoted-string  = DQUOTE *QcontentSMTP DQUOTE

   QcontentSMTP   = qtextSMTP / quoted-pairSMTP

   quoted-pairSMTP  = %d92 %d32-126
                    ; i.e., backslash followed by any ASCII
                    ; graphic (including itself) or SPace

   qtextSMTP      = %d32-33 / %d35-91 / %d93-126
                  ; i.e., within a quoted string, any
                  ; ASCII graphic or space is permitted
                  ; without blackslash-quoting except
                  ; double-quote and the backslash itself.

   String         = Atom / Quoted-string
*/

// some regexes we don't want to compile for each request
var (
	localRegex = regexp.MustCompile("^[a-zA-Z0-9,!#\\$%&'\\*\\+/=\\?\\^_`\\{\\|}~-]+(\\.[a-zA-Z0-9,!#\\$%&'\\*\\+/=\\?\\^_`\\{\\|}~-]+)*$")
	// TODO: quoted-string and more special chars
)

func (m *MailAddress) Validate() (bool, string) {
	// Check lengths
	if len(m.Local) > 64 {
		return false, "Local too long"
	}
	if len(m.Domain) > 253 {
		return false, "Domain too long"
	}
	if len(m.Domain)+len(m.Local) > 254 {
		return false, "MailAddress too long"
	}
	if !localRegex.MatchString(m.Local) {
		return false, "Invalid local part"
	}
	return true, ""
}

/*
   RFC 5321

   The maximum total length of a domain name or number is 255 octets.
*/

// Check if m.Domain reverses to conn.
func (m *MailAddress) HasReverseDns(conn *conn) bool {
	// TODO
	// check for IP address
	ip := net.ParseIP(m.Domain)
	connAddr, ok := (conn.c.RemoteAddr()).(*net.TCPAddr)
	if !ok {
		log.Printf("    > Connection %s isn't a tcp connection", conn.c.RemoteAddr())
		return false
	}

	if ip != nil {
		// it's an IP
		if !ip.Equal(connAddr.IP) {
			log.Printf("    > IP in from(%s) doesn't match real IP(%s)", ip, connAddr.IP)
			return false
		}

	} else {
		// try to interpret is as a domain
		// check for rDNS of client IP
		domains, err := net.LookupAddr(connAddr.IP.String())
		if err != nil {
			log.Printf("    > rDNS lookup failed: %s", err)
			return false
		}

		if !stringInSlice(m.Domain, domains) {
			log.Printf("    > rDNS(%s) didn't match Domain(%s)", domains, m.Domain)
			return false
		}

		// if no rDNS match found, check for the SPF record
		// TODO
	}

	return true
}

// Check if we are m.Domain.
func (m *MailAddress) IsLocal(conn *conn) bool {
	// TODO: Check the domain for real :p
	return m.Domain == "gopistolet.be"
}

func stringInSlice(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
