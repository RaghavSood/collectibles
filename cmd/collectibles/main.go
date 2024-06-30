package main

import (
	"github.com/RaghavSood/collectibles/web"
)

func main() {
	webServer := web.NewServer()
	webServer.Serve()
}
