package replay

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

type ResponseAnalyzer interface {
	ResponseAnalyze(*HttpResponse)
}

// ForwardHost where to forward requests
type ForwardHost struct {
	Url   string
	Limit int

	Stat *RequestStat
}

// ReplaySettings ListenerSettings contain all the needed configuration for setting up the replay
type ReplaySettings struct {
	Port int
	Host string

	Address string

	ForwardAddress string

	FileToReplyPath string

	Verbose bool

	ElastiSearchURI string

	ResponseAnalyzePlugins []ResponseAnalyzer
}

var Settings ReplaySettings = ReplaySettings{}

func (r *ReplaySettings) RegisterResponseAnalyzePlugin(p ResponseAnalyzer) {
	r.ResponseAnalyzePlugins = append(r.ResponseAnalyzePlugins, p)
}

// ForwardedHosts implements forwardAddress syntax support for multiple hosts (coma separated), and rate limiting by specifing "|maxRps" after host name.
//
//    -f "host1,http://host2|10,host3"
//
func (r *ReplaySettings) ForwardedHosts() (hosts []*ForwardHost) {
	hosts = make([]*ForwardHost, 0, 10)

	for _, address := range strings.Split(r.ForwardAddress, ",") {
		host_info := strings.Split(address, "|")

		if strings.Index(host_info[0], "http") == -1 {
			host_info[0] = "http://" + host_info[0]
		}

		host := &ForwardHost{Url: host_info[0]}
		host.Stat = NewRequestStats(host)

		if len(host_info) > 1 {
			host.Limit, _ = strconv.Atoi(host_info[1])
		}

		hosts = append(hosts, host)
	}

	return
}

func (r *ReplaySettings) Parse() {
	r.Address = r.Host + ":" + strconv.Itoa(r.Port)

	// Register Plugins
	// Elasticsearch Plugin
	if Settings.ElastiSearchURI != "" {
		esp := &ESPlugin{}
		esp.Init(Settings.ElastiSearchURI)

		r.RegisterResponseAnalyzePlugin(esp)
	}
}

func init() {
	if len(os.Args) < 2 || os.Args[1] != "replay" {
		return
	}

	const (
		defaultPort = 28020
		defaultHost = "0.0.0.0"

		defaultForwardAddress = "http://localhost:8080"
	)

	flag.IntVar(&Settings.Port, "p", defaultPort, "specify port number")

	flag.StringVar(&Settings.Host, "ip", defaultHost, "ip addresses to listen on")

	flag.StringVar(&Settings.ForwardAddress, "f", defaultForwardAddress, "http address to forward traffic.\n\tYou can limit requests per second by adding `|num` after address.\n\tIf you have multiple addresses with different limits. For example: http://staging.example.com|100,http://dev.example.com|10")

	flag.StringVar(&Settings.FileToReplyPath, "file", "", "File to replay captured requests from")

	flag.BoolVar(&Settings.Verbose, "verbose", false, "Log requests")

	flag.StringVar(&Settings.ElastiSearchURI, "es", "", "enable elasticsearch\n\tformat: hostname:9200/index_name")
}
