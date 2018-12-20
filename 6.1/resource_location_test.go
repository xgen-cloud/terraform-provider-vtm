// Copyright (C) 2018, Pulse Secure, LLC. 
// Licensed under the terms of the MPL 2.0. See LICENSE file for details.

package main

/*
 * This test covers the following cases:
 *   - Creation and deletion of a vtm_location object with minimal configuration
 */

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	vtm "github.com/pulse-vadc/go-vtm/6.1"
)

func TestResourceLocation(t *testing.T) {
	objName := acctest.RandomWithPrefix("TestLocation")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: getBasicLocationConfig(objName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationExists,
				),
			},
		},
	})
}

func testAccCheckLocationExists(s *terraform.State) error {
	for _, tfResource := range s.RootModule().Resources {
		if tfResource.Type != "vtm_location" {
			continue
		}
		objectName := tfResource.Primary.Attributes["name"]
		tm := testAccProvider.Meta().(*vtm.VirtualTrafficManager)
		if _, err := tm.GetLocation(objectName); err != nil {
			return fmt.Errorf("Location %s does not exist: %#v", objectName, err)
		}
	}

	return nil
}

func testAccCheckLocationDestroy(s *terraform.State) error {
	for _, tfResource := range s.RootModule().Resources {
		if tfResource.Type != "vtm_location" {
			continue
		}
		objectName := tfResource.Primary.Attributes["name"]
		tm := testAccProvider.Meta().(*vtm.VirtualTrafficManager)
		if _, err := tm.GetLocation(objectName); err == nil {
			return fmt.Errorf("Location %s still exists", objectName)
		}
	}

	return nil
}

func getBasicLocationConfig(name string) string {
	return fmt.Sprintf(`
        resource "vtm_location" "test_vtm_location" {
			name = "%s"
			identifier = 10

        }`,
		name,
	)
}
