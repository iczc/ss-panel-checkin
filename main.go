package main

import (
	"flag"
	"log"

	"github.com/ss-panel-checkin/checkin"
)

var (
	host   string
	email  string
	passwd string
)

func init() {
	flag.StringVar(&host, "host", "", "website host")
	flag.StringVar(&email, "email", "", "user email")
	flag.StringVar(&passwd, "passwd", "", "user password")
	flag.Parse()
}

func main() {
	c := checkin.NewCheckin(host, email, passwd)
	if err := c.Handle(); err != nil {
		log.Panic(err)
	}
}
