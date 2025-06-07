package mysql_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"strconv"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccCtyunMysqlSpecs(t *testing.T) {

	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}

	dnd := utils.GenerateRandomString()

	datasourceName := "data.ctyun_mysql_specs." + dnd
	datasourceFile := "datasource_ctyun_mysql_specs.tf"

	prodType := "1"
	prodCode := "MYSQL"
	instanceType := "1"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(datasourceFile, dnd, prodType, prodCode, instanceType),
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						ds := s.RootModule().Resources[datasourceName].Primary
						count, err := strconv.Atoi(ds.Attributes["specs.#"])
						if err != nil || count == 0 {
							return fmt.Errorf("specs 无效: %v", ds.Attributes)
						}
						return nil
					},
				),
			},
		},
	})
}
