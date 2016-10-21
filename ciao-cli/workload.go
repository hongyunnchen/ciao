//
// Copyright (c) 2016 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/01org/ciao/openstack/compute"
)

var workloadCommand = &command{
	SubCommands: map[string]subCommand{
		"list": new(workloadListCommand),
	},
}

type workloadListCommand struct {
	Flag     flag.FlagSet
	template string
}

func (cmd *workloadListCommand) usage(...string) {
	fmt.Fprintf(os.Stderr, `usage: ciao-cli [options] workload list

List all workloads

`)
	cmd.Flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
The template passed to the -f option operates on a 

[]struct {
	OSFLVDISABLEDDisabled  bool    // Not used
	Disk                   string  // Backing images associated with workload
	OSFLVEXTDATAEphemeral  int     // Not currently used
	OsFlavorAccessIsPublic bool    // Indicates whether the workload is available to all tenants
	ID                     string  // ID of the workload
	Links                  []Link  // Not currently used
	Name                   string  // Name of the workload
	RAM                    int     // Amount of RAM allocated to instances of this workload 
	Swap                   string  // Not currently used
	Vcpus                  int     // Number of Vcpus allocated to instances of this workload 
}
`)

	os.Exit(2)
}

func (cmd *workloadListCommand) parseArgs(args []string) []string {
	cmd.Flag.StringVar(&cmd.template, "f", "", "Template used to format output")
	cmd.Flag.Usage = func() { cmd.usage() }
	cmd.Flag.Parse(args)
	return cmd.Flag.Args()
}

func (cmd *workloadListCommand) run(args []string) error {
	if *tenantID == "" {
		fatalf("Missing required -tenant-id parameter")
	}

	var flavors compute.FlavorsDetails
	if *tenantID == "" {
		*tenantID = "faketenant"
	}

	url := buildComputeURL("%s/flavors/detail", *tenantID)

	resp, err := sendHTTPRequest("GET", url, nil, nil)
	if err != nil {
		fatalf(err.Error())
	}

	err = unmarshalHTTPResponse(resp, &flavors)
	if err != nil {
		fatalf(err.Error())
	}

	if cmd.template != "" {
		return outputToTemplate("workload-list", cmd.template,
			&flavors.Flavors)
	}

	for i, flavor := range flavors.Flavors {
		fmt.Printf("Workload %d\n", i+1)
		fmt.Printf("\tName: %s\n\tUUID:%s\n\tImage UUID: %s\n\tCPUs: %d\n\tMemory: %d MB\n",
			flavor.Name, flavor.ID, flavor.Disk, flavor.Vcpus, flavor.RAM)
	}
	return nil
}
