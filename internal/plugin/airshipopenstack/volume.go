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
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
	"github.com/gophercloud/gophercloud/pagination"
)

// gets OpenStack volumes viewable by the tenant
// https://docs.openstack.org/cinder/latest/
func GetVolumes() []map[string]string {
	osp := NewOpenstackPlugin()
	m := make([]map[string]string, 0)

	err := volumes.List(volumeClientHelper(osp), volumes.ListOpts{}).EachPage(
		func(page pagination.Page) (bool, error) {
			volumeList, err := volumes.ExtractVolumes(page)

			if err != nil {
				log.Printf("Broken at volume list %v\n", err)
				return false, err
			}

			for _, volume := range volumeList {
				// extract the potentially multiple devices
				var attachedDevices []string
				if len(vmCache) == 0 {
					GetVMs()
				}
				for index := range volume.Attachments {
					attachment := volume.Attachments[index]
					attachedDevices = append(attachedDevices, vmCache[attachment.ServerID]+" at "+attachment.Device)
				}

				m = append(m, map[string]string{
					"Name":              volume.Name,
					"Description":       volume.Description,
					"Size":              strconv.Itoa(volume.Size) + "GiB",
					"Status":            volume.Status,
					"Type":              volume.VolumeType,
					"Attached To":       strings.Join(attachedDevices, ", "),
					"Availability Zone": volume.AvailabilityZone,
					"Bootable":          volume.Bootable,
					"Encrypted":         strconv.FormatBool(volume.Encrypted),
				})
			}

			return true, nil
		})

	if err != nil {
		log.Printf("volume list error: %s\n", err)
	}

	return m
}

// helper function to create an volume specific gophercloud client
func volumeClientHelper(osp *OpenstackPlugin) *gophercloud.ServiceClient {
	client, err := openstack.NewBlockStorageV2(osp.provider, gophercloud.EndpointOpts{})

	if err != nil {
		log.Fatalf("NewIdentityV3 error: %v", err)
	}

	return client
}
