package mongodb_test

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"os"
	"testing"
)

func TestAccCtyunMongodbInstances(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	dnd := utils.GenerateRandomString()

	datasourceName := "data.ctyun_mongodb_instances." + dnd
	datasourceFile := "datasource_ctyun_mongodb_instances.tf"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 绑定IP验证
			{
				Config: utils.LoadTestCase(datasourceFile, dnd, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "mongodb_instances.#"),
				),
			},
		},
	})
}
