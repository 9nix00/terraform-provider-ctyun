package mongodb_test

//func TestAccCtyunMongodbInstances(t *testing.T) {
//
//	dnd := utils.GenerateRandomString()
//
//	datasourceName := "data.ctyun_mongodb_instances." + dnd
//	datasourceFile := "datasource_ctyun_mongodb_instances.tf"
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
//		Steps: []resource.TestStep{
//			// 绑定IP验证
//			{
//				Config: utils.LoadTestCase(datasourceFile, dnd, ""),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttrSet(datasourceName, "mongodb_instances.#"),
//				),
//			},
//		},
//	})
//}
