/*
Copyright (c) 2020 AT&T. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package airshipopenstack

import (
	"log"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
)

// gets OpenStack flavors that are available for the tenant
// https://docs.openstack.org/nova/latest/user/flavors.html
func GetFlavors() []map[string]string {
	osp := NewOpenstackPlugin()
	m := make([]map[string]string, 0)

	err := flavors.ListDetail(computeClientHelper(osp), flavors.ListOpts{AccessType: flavors.AllAccess}).EachPage(
		func(page pagination.Page) (bool, error) {
			flavorList, err := flavors.ExtractFlavors(page)

			if err != nil {
				log.Printf("compute flavor Error: %s\n", err)
			}

			for _, flavor := range flavorList {
				name := flavor.Name
				flavorCache[flavor.ID] = name

				m = append(m, map[string]string{
					"Name":           name,
					"VCPUs":          strconv.Itoa(flavor.VCPUs),
					"RAM":            strconv.Itoa(flavor.RAM),
					"Root Disk":      strconv.Itoa(flavor.Disk),
					"Ephemeral Disk": strconv.Itoa(flavor.Ephemeral),
					"Swap Disk":      strconv.Itoa(flavor.Swap),
					"RX/TX factor":   strconv.FormatFloat(flavor.RxTxFactor, 'f', 1, 64),
					"Public":         strconv.FormatBool(flavor.IsPublic),
				})
			}

			return true, nil
		})

	if err != nil {
		log.Printf("compute flavor list error: %s\n", err)
	}

	return m
}

// gets OpenStack images that are available for the tenant
// https://docs.openstack.org/image-guide/
func GetImages() []map[string]string {
	osp := NewOpenstackPlugin()
	m := make([]map[string]string, 0)

	err := images.ListDetail(computeClientHelper(osp), images.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			list, err := images.ExtractImages(page)

			if err != nil {
				log.Printf("Image list error: %s\n", err)
				return false, err
			}

			for _, image := range list {
				name := image.Name
				imageCache[image.ID] = name

				m = append(m, map[string]string{
					"Name":   name,
					"Type":   strconv.Itoa(image.Progress),
					"Status": image.Status,
				})
			}

			return true, nil
		})

	if err != nil {
		log.Printf("compute image list error: %s\n", err)
	}

	return m
}

// gets OpenStack vms that are configured for the tenant
// https://docs.openstack.org/python-openstackclient/latest/cli/command-objects/server.html
func GetVMs() []map[string]string {
	osp := NewOpenstackPlugin()
	m := make([]map[string]string, 0)

	err := servers.List(computeClientHelper(osp), servers.ListOpts{AllTenants: true}).EachPage(
		func(page pagination.Page) (bool, error) {
			list, err := servers.ExtractServers(page)

			if err != nil {
				log.Fatalf("Server list error: %s\n", err)
				return false, err
			}

			for _, server := range list {
				name := server.Name
				vmCache[server.ID] = name

				// pull together the addresses from the structure
				var addresses string
				var tmp []string
				for subnetName := range server.Addresses {
					for _, v := range server.Addresses[subnetName].([]interface{}) {
						tmp = append(tmp, v.(map[string]interface{})["addr"].(string))
					}
				}

				if len(tmp) > 1 {
					addresses = strings.Join(tmp, ", ")
				} else {
					addresses = tmp[0]
				}

				if len(imageCache) == 0 {
					GetImages()
				}
				if len(flavorCache) == 0 {
					GetFlavors()
				}
				m = append(m, map[string]string{
					"Instance Name": name,
					"Image Name":    imageCache[server.Image["id"].(string)],
					"IP Address":    addresses,
					"Flavor":        flavorCache[server.Flavor["id"].(string)],
					"Key Pair":      server.KeyName,
					"Status":        server.Status,
					"Created":       server.Created.Local().String(),
				})
			}

			return true, nil
		})

	if err != nil {
		log.Printf("compute server list error: %s\n", err)
	}

	return m
}

// helper function to create a compute specific gophercloud client
func computeClientHelper(osp *OpenstackPlugin) *gophercloud.ServiceClient {
	client, err := openstack.NewComputeV2(osp.provider, gophercloud.EndpointOpts{
		Type: "compute",
	})

	if err != nil {
		log.Fatalf("Compute Client Error: %s\n", err)
		return nil
	}

	return client
}
