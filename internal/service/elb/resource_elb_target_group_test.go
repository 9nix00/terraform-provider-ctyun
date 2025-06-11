package elb_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-ctyun/internal/service"
	"terraform-provider-ctyun/internal/utils"
	"testing"
)

func TestAccCtyunElbTargetGroup(t *testing.T) {

	rnd := utils.GenerateRandomString()
	dnd := utils.GenerateRandomString()

	resourceName := "ctyun_elb_target_group." + rnd
	resourceFile := "resource_ctyun_elb_target_group.tf"

	datasourceName := "data.ctyun_elb_target_groups." + dnd
	datasourceFile := "datasource_ctyun_elb_target_groups.tf"

	name := "target_groups_" + utils.GenerateRandomString()
	algorithm := "wrr"

	updatedName := "target_groups_new_" + utils.GenerateRandomString()
	updatedAlgorithm := "lc"

	healthCheckID := dependence.healthCheckID
	sessionStickyMode := "SOURCE_IP"
	//cookieExpire := 30
	//rewriteCookieName := "cookie_name" + utils.GenerateRandomString()
	sourceIpTimeout := 30
	proxyProtocol := 1
	protocol := "TCP"

	tfHealthCheckID := fmt.Sprintf(`health_check_id="%s"`, healthCheckID)
	tfSessionStickyMode := fmt.Sprintf(`session_sticky_mode="%s"`, sessionStickyMode)
	//tfCookieExpire := fmt.Sprintf(`cookie_expire=%d`, cookieExpire)
	//tfRewriteCookieName := fmt.Sprintf(`rewrite_cookie_name="%s"`, rewriteCookieName)
	tfSourceIpTimeout := fmt.Sprintf(`source_ip_timeout=%d`, sourceIpTimeout)
	tfProxyProtocol := fmt.Sprintf(`proxy_protocol=%d`, proxyProtocol)
	tfProtocol := fmt.Sprintf(`protocol="%s"`, protocol)

	insertSessionStickyMode := "INSERT"
	rewriteSessionStickyMode := "REWRITE"
	updatedCookieExpire := 1
	updatedRewriteCookieName := "cookie_name_new" + utils.GenerateRandomString()
	updatedSourceIpTimeout := 1
	updatedProxyProtocol := 0
	updatedProtocol := "HTTP"

	insertTfSessionStickyMode := fmt.Sprintf(`session_sticky_mode="%s"`, insertSessionStickyMode)
	rewriteTfSessionStickyMode := fmt.Sprintf(`session_sticky_mode="%s"`, rewriteSessionStickyMode)
	updatedTfCookieExpire := fmt.Sprintf(fmt.Sprintf(`cookie_expire=%d`, updatedCookieExpire))
	updatedTfRewriteCookieName := fmt.Sprintf(`rewrite_cookie_name="%s"`, updatedRewriteCookieName)
	updatedTfSourceIpTimeout := fmt.Sprintf(`source_ip_timeout=%d`, updatedSourceIpTimeout)
	updatedTfProxyProtocol := fmt.Sprintf(`proxy_protocol=%d`, updatedProxyProtocol)
	updatedTfProtocol := fmt.Sprintf(`protocol="%s"`, updatedProtocol)

	closedTfSessionStickyMode := fmt.Sprintf(`session_sticky_mode="%s"`, "CLOSE")
	// 代码合并需要整改
	vpcId := dependence.vpcID

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
			// 1. 基础功能测试
			// 1.1 create验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, vpcId, algorithm, "", "", "", "", "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "algorithm", algorithm),
				),
			},
			// 1.2 update 验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, vpcId, updatedAlgorithm, "", "", "", "", "", "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "algorithm", updatedAlgorithm),
				),
			},
			// 1.3 datasource验证
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, vpcId, updatedAlgorithm, "", "", "", "", "", "", "") +
					utils.LoadTestCase(datasourceFile, dnd, fmt.Sprintf(`ids=%s.id`, resourceName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "target_groups.0.name", updatedName),
					resource.TestCheckResourceAttr(datasourceName, "target_groups.0.vpc_id", vpcId),
					resource.TestCheckResourceAttr(datasourceName, "target_groups.0.algorithm", updatedAlgorithm),
				),
			},
			// 1.4 delete 验证
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, updatedName, vpcId, updatedAlgorithm, "", "", "", "", "", "", ""),
				Destroy: true,
			},

			// 2. 详细参数创建，包括健康检查id， algorithm=wrr, sessionStickyMode=SOURCE_IP
			// 2.1 create验证，包括创建一个健康检查
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, vpcId, algorithm, tfHealthCheckID, tfSessionStickyMode, "", "", tfSourceIpTimeout, tfProxyProtocol, tfProtocol),
				//Config: utils.LoadTestCase(resourceFile, rnd, updatedName, vpcId, algorithm, regionID, tfHealthCheckID, updatedTfSessionStickyMode, updatedTfCookieExpire, "", "", updatedTfProxyProtocol, updatedTfProtocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "algorithm", algorithm),
					//resource.TestCheckResourceAttr(resourceName, "health_check_id", healthCheckID),
					resource.TestCheckResourceAttr(resourceName, "session_sticky_mode", sessionStickyMode),
					//resource.TestCheckResourceAttr(resourceName, "cookie_expire", strconv.Itoa(cookieExpire)),
					//resource.TestCheckResourceAttr(resourceName, "rewrite_cookie_name", rewriteCookieName),
					resource.TestCheckResourceAttr(resourceName, "source_ip_timeout", strconv.Itoa(sourceIpTimeout)),
					resource.TestCheckResourceAttr(resourceName, "proxy_protocol", strconv.Itoa(proxyProtocol)),
				),
			},
			// 2.2 update验证, algorithm=wrr, sessionStickyMode=SOURCE_IP,
			{
				Config: utils.LoadTestCase(resourceFile, rnd, updatedName, vpcId, algorithm, tfHealthCheckID, tfSessionStickyMode, "", "", updatedTfSourceIpTimeout, tfProxyProtocol, tfProtocol),
				//Config: utils.LoadTestCase(resourceFile, rnd, updatedName, vpcId, algorithm, regionID, "", tfSessionStickyMode, "", "", tfSourceIpTimeout, tfProxyProtocol, tfProtocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "algorithm", algorithm),
					resource.TestCheckResourceAttr(resourceName, "health_check_id", healthCheckID),
					resource.TestCheckResourceAttr(resourceName, "session_sticky_mode", sessionStickyMode),
					//resource.TestCheckResourceAttr(resourceName, "cookie_expire", strconv.Itoa(updatedCookieExpire)),
					//resource.TestCheckResourceAttr(resourceName, "rewrite_cookie_name", updatedRewriteCookieName),
					resource.TestCheckResourceAttr(resourceName, "source_ip_timeout", strconv.Itoa(updatedSourceIpTimeout)),
					resource.TestCheckResourceAttr(resourceName, "proxy_protocol", strconv.Itoa(proxyProtocol)),
				),
			},
			// 2.3 销毁
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, updatedName, vpcId, algorithm, "", tfSessionStickyMode, "", "", updatedTfSourceIpTimeout, updatedTfProxyProtocol, tfProtocol),
				Destroy: true,
			},
			// 3. 验证SessionSticky修改,算法修改
			// 3.1 create algorithm=wrr, sessionStickyMode=INSERT, cookieExpire=1, proxyProtocol=0, protocol=http
			// 当proxy=http, proxyProtocol=1已验证，不可创建
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, vpcId, algorithm, "", insertTfSessionStickyMode, updatedTfCookieExpire, "", "", updatedTfProxyProtocol, updatedTfProtocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "algorithm", algorithm),
					//resource.TestCheckResourceAttr(resourceName, "health_check_id", healthCheckID),
					resource.TestCheckResourceAttr(resourceName, "session_sticky_mode", insertSessionStickyMode),
					resource.TestCheckResourceAttr(resourceName, "cookie_expire", strconv.Itoa(updatedCookieExpire)),
					resource.TestCheckResourceAttr(resourceName, "proxy_protocol", strconv.Itoa(updatedProxyProtocol)),
					resource.TestCheckResourceAttr(resourceName, "protocol", updatedProtocol),
				),
			},
			// 3.2 updated, algorithm=wrr, sessionStickyMode=REWRITE, cookieExpire=1, proxyProtocol=0, protocol=http
			// 除验证更改sessionSticky外，也需要验证算法=lc/sh时，sessionSticky非CLOSE选项是否可行，理论上不可行 ---- 结论：无法创建
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, vpcId, algorithm, "", rewriteTfSessionStickyMode, "", updatedTfRewriteCookieName, "", updatedTfProxyProtocol, updatedTfProtocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "algorithm", algorithm),
					//resource.TestCheckResourceAttr(resourceName, "health_check_id", healthCheckID),
					resource.TestCheckResourceAttr(resourceName, "session_sticky_mode", rewriteSessionStickyMode),
					resource.TestCheckResourceAttr(resourceName, "rewrite_cookie_name", updatedRewriteCookieName),
					resource.TestCheckResourceAttr(resourceName, "proxy_protocol", strconv.Itoa(updatedProxyProtocol)),
					resource.TestCheckResourceAttr(resourceName, "protocol", updatedProtocol),
				),
			},
			// 3.3 updated， algorithm=lc, sessionStickyMode=CLOSE， proxyProtocol=0, protocol=http
			{
				Config: utils.LoadTestCase(resourceFile, rnd, name, vpcId, updatedAlgorithm, "", closedTfSessionStickyMode, "", "", "", updatedTfProxyProtocol, updatedTfProtocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", vpcId),
					resource.TestCheckResourceAttr(resourceName, "algorithm", updatedAlgorithm),
					//resource.TestCheckResourceAttr(resourceName, "health_check_id", healthCheckID),
					resource.TestCheckResourceAttr(resourceName, "session_sticky_mode", "CLOSE"),
					resource.TestCheckResourceAttr(resourceName, "proxy_protocol", strconv.Itoa(updatedProxyProtocol)),
					resource.TestCheckResourceAttr(resourceName, "protocol", updatedProtocol),
				),
			},
			// 3.4 Destroy
			{
				Config:  utils.LoadTestCase(resourceFile, rnd, name, vpcId, updatedAlgorithm, "", closedTfSessionStickyMode, "", "", "", updatedTfProxyProtocol, updatedTfProtocol),
				Destroy: true,
			},
		},
	})
}
