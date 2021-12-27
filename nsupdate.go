package nsupdate // import "golang.voidptr.de/nsupdate"

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
)

func SendUpdate() {
	var m dns.Msg
	m.SetUpdate("example.com.")

	var newRRs []dns.RR

	rr, _ := dns.NewRR("test.example.com. 60 IN A 127.2.3.4")
	newRRs = append(newRRs, rr)

	m.Insert(newRRs)

	var client dns.Client

	tsigOptions, _ := ReadKeyFile("example-key.conf")

	fmt.Println(tsigOptions)

	m.SetTsig(tsigOptions.Name, tsigOptions.Algorithm, 300, time.Now().Unix())
	client.TsigSecret = map[string]string{
		tsigOptions.Name: tsigOptions.Secret,
	}

	res1, res2, err := client.Exchange(&m, "192.0.2.2:53")
	fmt.Println(res1, res2, err)
}
