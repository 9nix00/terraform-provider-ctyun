package mysql_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"os"
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

	prodType := "RDS"
	prodCode := "MYSQL"
	instanceType := "通用型"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: utils.LoadTestCase(datasourceFile, dnd, prodType, prodCode, instanceType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "specs"),
				),
			},
		},
	})
}
