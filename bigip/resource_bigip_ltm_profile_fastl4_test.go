package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_FASTL4_NAME = fmt.Sprintf("/%s/test-fastl4", TEST_PARTITION)

var TEST_FASTL4_RESOURCE = `
resource "bigip_fastl4_profile" "test-fastl4" {
            name = "` + TEST_FASTL4_NAME + `"
            partition = "Common"
            defaults_from = "fastL4"
						client_timeout = 40
						idle_timeout = "200"
            explicitflow_migration = "enabled"
            hardware_syncookie = "enabled"
            iptos_toclient = "pass-through"
            iptos_toserver = "pass-through"
            keepalive_interval = "disabled"
 }
`

func TestBigipLtmFastl4_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckfastl4sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_FASTL4_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckfastl4Exists(TEST_FASTL4_NAME, true),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "name", TEST_FASTL4_NAME),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "partition", "Common"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "defaults_from", "fastL4"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "client_timeout", "40"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "explicitflow_migration", "enabled"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "hardware_syncookie", "enabled"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "idle_timeout", "200"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "hardware_syncookie", "enabled"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "iptos_toclient", "pass-through"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "iptos_toserver", "pass-through"),
					resource.TestCheckResourceAttr("bigip_fastl4_profile.test-fastl4", "keepalive_interval", "disabled"),
				),
			},
		},
	})
}

func TestBigipLtmfastl4_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckfastl4sDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_FASTL4_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckfastl4Exists(TEST_FASTL4_NAME, true),
				),
				ResourceName:      TEST_FASTL4_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckfastl4Exists(name string, exists bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)
		p, err := client.Fastl4(name)
		if err != nil {
			return err
		}
		if exists && p == nil {
			return fmt.Errorf("fastl4 ", name, " was not created.")
		}
		if !exists && p != nil {
			return fmt.Errorf("fastl4 ", name, " still exists.")
		}
		return nil
	}
}

func testCheckfastl4sDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_fastl4_profile" {
			continue
		}

		name := rs.Primary.ID
		fastl4, err := client.Fastl4(name)
		if err != nil {
			return err
		}
		if fastl4 == nil {
			return fmt.Errorf("fastl4 tata ", name, " not destroyed.")
		}
	}
	return nil
}
