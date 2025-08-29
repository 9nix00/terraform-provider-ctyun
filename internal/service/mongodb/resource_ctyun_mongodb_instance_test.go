package mongodb_test

import (
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/service"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"strconv"
	"testing"
)

func TestAccCtyunMongodbInstanceSingleOnDemand(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_single_on_demand.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kqjwyk123="
	prodId := "Single34"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 120
	backupStorageType := "SATA"
	backupStorageSpace := 150
	azInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"master"}, {"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"backup"}]`

	//更新参数
	updatedName := "tf-mongodb-single-new-" + utils.GenerateRandomString()
	updatedFlavorName := "s7.large.4"
	updatedReadPort := 12348
	//updatedStorageType := ""
	updatedStorageSpace := 130
	//backupStorageType := "SATA"
	updatedBackupStorageSpace := 160
	updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"s"}]`

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, backupStorageSpace, azInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_space", strconv.Itoa(backupStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 更新mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedBackupStorageSpace, updatedAzInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_space", strconv.Itoa(updatedBackupStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedBackupStorageSpace, updatedAzInfo),
				Destroy: true,
			},
		},
	})
}

// 创建包周期，且无传AZ信息, 备份空间为os
func TestAccCtyunMongodbInstanceSingleOnDemandNoAz(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_single_cycle_no_az_os.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "month"
	cycleCount := 1
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kqjwyk123="
	prodId := "Single34"
	readPort := 12345
	storageType := "SATA"
	storageSpace := 100
	backupStorageType := "OS"
	//backupStorageSpace := 100
	//azInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"master"}, {"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"backup"}]`

	//更新参数
	updatedName := "tf-mongodb-single-new-" + utils.GenerateRandomString()
	updatedFlavorName := "s7.large.4"
	updatedReadPort := 12348
	//updatedStorageType := ""
	updatedStorageSpace := 110
	//backupStorageType := "SATA"
	//updatedBackupStorageSpace := 160

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, cycleCount, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType,
					storageSpace, backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 更新mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, cycleCount, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, cycleCount, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType),
				Destroy: true,
			},
		},
	})
}

// 验证副本集，传azList，OS存储
func TestAccCtyunMongodbInstanceReplicaOs(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_replica_on_demand_os.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kqjwyk123="
	prodId := "Replica3R34"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "OS"
	azInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"master"},
				{"availability_zone_name":"cn-huadong1-jsnj2A-public-ctcloud","availability_zone_count":1,"node_type":"master"},
				{"availability_zone_name":"cn-huadong1-jsnj3A-public-ctcloud","availability_zone_count":1,"node_type":"master"}]`

	//更新参数
	updatedName := "tf-mongodb-single-new-" + utils.GenerateRandomString()
	updatedFlavorName := "s7.large.4"
	updatedReadPort := 12348
	updatedProdId := "Replica5R34"

	//updatedStorageType := ""
	updatedStorageSpace := 110
	updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`

	updatedSpecAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"ms"},
				{"availability_zone_name":"cn-huadong1-jsnj2A-public-ctcloud","availability_zone_count":1,"node_type":"ms"},
				{"availability_zone_name":"cn-huadong1-jsnj3A-public-ctcloud","availability_zone_count":1,"node_type":"ms"}]`
	//backupStorageType := "SATA"
	//updatedBackupStorageSpace := 160
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"s"}]`

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, azInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 更新mongodb实例，升级存储空间和flavor_name
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedSpecAzInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 更新mongodb实例，升级存储空间
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, updatedProdId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedAzInfo),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, updatedProdId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedAzInfo),
				Destroy: true,
			},
		},
	})
}

// 副本集，不传azList, 存储为SATA
func TestAccCtyunMongodbInstanceReplicaSATANoAzList(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_replica_on_demand_no_az.tf"
	// datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kqjwyk123="
	prodId := "Replica3R34"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "SATA"
	backupStorageSpace := 120
	//azInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"master"},
	//			{"availability_zone_name":"cn-huadong1-jsnj2A-public-ctcloud","availability_zone_count":1,"node_type":"master"},
	//			{"availability_zone_name":"cn-huadong1-jsnj3A-public-ctcloud","availability_zone_count":1,"node_type":"master"}]`
	replicaNum := 3

	//更新参数
	updatedName := "tf-mongodb-single-new-" + utils.GenerateRandomString()
	updatedFlavorName := "s7.large.4"
	updatedReadPort := 12348
	//updatedStorageType := ""
	updatedStorageSpace := 110
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"master"}]`
	updatedReplicaNum := 5
	//backupStorageType := "SATA"
	updatedBackupStorageSpace := 160
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"s"}]`

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, backupStorageSpace, replicaNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "replica_num", strconv.Itoa(replicaNum)),
				),
			},
			// 更新mongodb实例，升级存储空间
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedBackupStorageSpace, updatedReplicaNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "replica_num", strconv.Itoa(updatedReplicaNum)),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, prodId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedBackupStorageSpace, updatedReplicaNum),
				Destroy: true,
			},
		},
	})
}

// 集群版，传azList， OS存储
func TestAccCtyunMongodbInstanceClusterOs(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_cluster_on_demand_os.tf"
	resourceFile1 := "resource_ctyun_mongodb_instance_cluster_on_demand_os_update.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kqjwyk123="
	prodId := "Replica3R34"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "OS"
	azInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"master"},
				{"availability_zone_name":"cn-huadong1-jsnj2A-public-ctcloud","availability_zone_count":1,"node_type":"master"},
				{"availability_zone_name":"cn-huadong1-jsnj3A-public-ctcloud","availability_zone_count":1,"node_type":"master"}]`
	shardNum := 2
	mongosNum := 2

	//更新参数
	updatedName := "tf-mongodb-single-new-" + utils.GenerateRandomString()
	updatedFlavorName := "s7.large.4"
	updatedReadPort := 12348
	updatedProdId := "Replica5R34"

	//updatedStorageType := ""
	updatedStorageSpace := 110
	updatedSpecAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`
	updatedProdIDAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`
	updatedShardNum := 3
	updatedMongosNum := 3
	//
	upgradeNodeTypeMongos := "mongos"
	upgradeNodeTypeShard := "shard"

	//backupStorageType := "SATA"
	//updatedBackupStorageSpace := 160
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"s"}]`

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, azInfo, shardNum, mongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(shardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(mongosNum)),
				),
			},
			// 更新mongodb实例，升级存储空间、mongos的spec
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, updatedProdId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedSpecAzInfo, shardNum, mongosNum, upgradeNodeTypeMongos),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 扩容
			// 更新mongodb实例，shard的spec
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, updatedProdId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedSpecAzInfo, shardNum, mongosNum, upgradeNodeTypeShard),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 扩容
			// 更新shard数量和mongos数量
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, updatedProdId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedProdIDAzInfo, updatedShardNum, updatedMongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(updatedShardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(updatedMongosNum)),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, updatedProdId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedSpecAzInfo, shardNum, mongosNum),
				Destroy: true,
			},
		},
	})
}

// 集群版，不传azList
func TestAccCtyunMongodbInstanceClusterNoAz(t *testing.T) {
	err := os.Setenv("TF_ACC", "1")
	if err != nil {
		return
	}
	rnd := utils.GenerateRandomString()
	//dnd := utils.GenerateRandomString()
	resourceName := "ctyun_mongodb_instance." + rnd
	//datasourceName := "data.ctyun_mongodb_instances." + dnd

	resourceFile := "resource_ctyun_mongodb_instance_cluster_on_demand_no_az.tf"
	resourceFile1 := "resource_ctyun_mongodb_instance_cluster_on_demand_no_az_update.tf"
	//datasourceFile := "datasource_ctyun_mongodb_instances.tf"
	// 创建参数
	cycleType := "on_demand"
	vpcID := dependence.vpcID
	flavorName := "s7.large.2"
	subnetID := dependence.subnetID
	securityGroupID := dependence.securityGroupID
	name := "tf-mongodb-single-" + utils.GenerateRandomString()
	password := "Kqjwyk123="
	prodId := "Replica3R34"
	readPort := 12345
	storageType := "SAS"
	storageSpace := 100
	backupStorageType := "SATA"
	backupStorageSpace := 120
	shardNum := 2
	mongosNum := 2

	//更新参数
	updatedName := "tf-mongodb-single-new-" + utils.GenerateRandomString()
	updatedFlavorName := "s7.large.4"
	updatedReadPort := 12348
	updatedProdId := "Replica5R34"

	//updatedStorageType := ""
	updatedStorageSpace := 110
	//updatedSpecAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`
	//updatedProdIDAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":2,"node_type":"ms"}]`
	updatedShardNum := 3
	updatedMongosNum := 3
	//
	upgradeNodeTypeMongos := "mongos"
	upgradeNodeTypeShard := "shard"

	//backupStorageType := "SATA"
	updatedBackupStorageSpace := 150
	//updatedAzInfo := `[{"availability_zone_name":"cn-huadong1-jsnj1A-public-ctcloud","availability_zone_count":1,"node_type":"s"}]`

	resource.Test(t, resource.TestCase{
		CheckDestroy: func(s *terraform.State) error {
			_, exists := s.RootModule().Resources[resourceName]
			if exists {
				return fmt.Errorf("resource destroy failed")
			}
			return nil
		},
		ProtoV6ProviderFactories: service.GetTestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// 创建一个单节点的mongodb实例
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, flavorName, subnetID, securityGroupID, name, password, prodId, readPort, storageType, storageSpace,
					backupStorageType, backupStorageSpace, shardNum, mongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(storageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(shardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(mongosNum)),
				),
			},
			// 更新mongodb实例，升级存储空间、mongos的spec
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, updatedProdId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedBackupStorageSpace, shardNum, mongosNum, upgradeNodeTypeMongos),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 扩容
			// 更新mongodb实例，shard的spec
			{
				Config: utils.LoadTestCase(resourceFile1, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, updatedProdId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedBackupStorageSpace, shardNum, mongosNum, upgradeNodeTypeShard),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(updatedReadPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", updatedFlavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
				),
			},
			// 扩容
			// 更新shard数量和mongos数量
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, updatedProdId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedBackupStorageSpace, updatedShardNum, updatedMongosNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					//resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "read_port", strconv.Itoa(readPort)),
					resource.TestCheckResourceAttr(resourceName, "storage_type", storageType),
					resource.TestCheckResourceAttr(resourceName, "storage_space", strconv.Itoa(updatedStorageSpace)),
					resource.TestCheckResourceAttr(resourceName, "backup_storage_type", backupStorageType),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcID),
					resource.TestCheckResourceAttr(resourceName, "flavor_name", flavorName),
					resource.TestCheckResourceAttr(resourceName, "subnet_id", subnetID),
					resource.TestCheckResourceAttr(resourceName, "security_group_id", securityGroupID),
					resource.TestCheckResourceAttr(resourceName, "shard_num", strconv.Itoa(updatedShardNum)),
					resource.TestCheckResourceAttr(resourceName, "mongos_num", strconv.Itoa(updatedMongosNum)),
				),
			},
			{
				Config: utils.LoadTestCase(resourceFile, rnd, cycleType, vpcID, updatedFlavorName, subnetID, securityGroupID, updatedName, password, updatedProdId, updatedReadPort,
					storageType, updatedStorageSpace, backupStorageType, updatedBackupStorageSpace, updatedShardNum, updatedMongosNum),
				Destroy: true,
			},
		},
	})
}
