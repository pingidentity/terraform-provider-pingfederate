package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/pingidentity/pingfederate",
		Debug:   debug,
	})

	if err != nil {
		fmt.Println(err)
	}
}
