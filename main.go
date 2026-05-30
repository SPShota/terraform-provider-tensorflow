package main

import (
	"context"
	"flag"
	"log"

	"github.com/SPShota/terraform-provider-tensorflow/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with debugger support")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/SPShota/tensorflow",
		Debug:   debug,
	}

	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.Fatal(err.Error())
	}
}
