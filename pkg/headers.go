package pkg

type HeaderMatchName string

const (
	MatchNameAddress                   HeaderMatchName = "address"
	MatchNameAuthorized                HeaderMatchName = "authorized"
	MatchNameBlocksIncomingConnections HeaderMatchName = "blocksincomingconnections"
	MatchNameClientVersion             HeaderMatchName = "clientversion"
	MatchNameExitStatus                HeaderMatchName = "exitstatus"
	MatchNameFullname                  HeaderMatchName = "fullname"
	MatchNameIpv4                      HeaderMatchName = "ipv4"
	MatchNameIpv6                      HeaderMatchName = "ipv6"
	MatchNameHostname                  HeaderMatchName = "hostname"
	MatchNameLastSeen                  HeaderMatchName = "lastseen"
	MatchNameLastSeenAgo               HeaderMatchName = "lastseen.ago"
	MatchNameMachine                   HeaderMatchName = "machine"
	MatchNameName                      HeaderMatchName = "name"
	MatchNameNo                        HeaderMatchName = "no"
	MatchNameOS                        HeaderMatchName = "os"
	MatchNameTags                      HeaderMatchName = "tags"
	MatchNameUser                      HeaderMatchName = "user"
	MatchNameVersion                   HeaderMatchName = "version"
)

type (
	Header struct {
		ReqEnriched bool
		MatchName   HeaderMatchName
		Title       string

		// TODO: factor in alias names?
		// Aliases []string

		// TODO: should we plan for this?
		// Width int
	}
)

var (
	HdrAddress     = Header{Title: "Address", MatchName: MatchNameAddress}
	HdrAuthorized  = Header{Title: "Authorized", MatchName: MatchNameAuthorized}
	HdrExitStatus  = Header{Title: "Exit Status", MatchName: MatchNameExitStatus}
	HdrIpv4        = Header{Title: "Ipv4", MatchName: MatchNameIpv4}
	HdrIpv6        = Header{Title: "Ipv6", MatchName: MatchNameIpv6}
	HdrLastSeenAgo = Header{Title: "Last Seen", MatchName: MatchNameLastSeenAgo, ReqEnriched: true}
	HdrMachine     = Header{Title: "Machine", MatchName: MatchNameMachine}
	HdrNo          = Header{Title: "No", MatchName: MatchNameNo}
	HdrTags        = Header{Title: "Tags", MatchName: MatchNameTags}
	HdrUser        = Header{Title: "User", MatchName: MatchNameUser}
	HdrVersion     = Header{Title: "Version", MatchName: MatchNameVersion}

	// AllHeadersList must contain the complete list of headers.
	AllHeadersList = []Header{
		HdrAddress,
		HdrAuthorized,
		HdrExitStatus,
		HdrIpv4,
		HdrIpv6,
		HdrLastSeenAgo,
		HdrMachine,
		HdrNo,
		HdrTags,
		HdrUser,
		HdrVersion,
	}

	// AllHeadersMap initializes a map of HeaderMatchName to Header.
	AllHeadersMap = func() map[HeaderMatchName]Header {
		a := make(map[HeaderMatchName]Header)
		for _, hdr := range AllHeadersList {
			a[hdr.MatchName] = hdr
		}
		return a
	}()

	// DefaultColumnSet is the column set that ships out of the box.
	// Order matters which is why it's created as a slice.
	DefaultColumnSet = []Header{
		HdrNo,
		HdrMachine,
		HdrIpv4,
		HdrTags,
		HdrUser,
		HdrVersion,
		HdrExitStatus,
		HdrLastSeenAgo,
	}
)
