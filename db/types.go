package db

type NodeType int

const (
	NodeTypeHostname NodeType = iota
	NodeTypeIP
	NodeTypePort
	// any file/dir found on one of the ports ( web stuff, ftp stuff, etc )
	NodeTypeContent
	// any info such as whois, emails, names, etc ... we'll make this more
	// specific once we'll have more osint modules
	NodeTypeInfo
)

func (t NodeType) String() string {
	switch t {
	case NodeTypeHostname:
		return "host"

	case NodeTypeIP:
		return "ip"

	case NodeTypePort:
		return "port"

	case NodeTypeContent:
		return "content"

	case NodeTypeInfo:
		return "info"
	}

	return "???"
}
