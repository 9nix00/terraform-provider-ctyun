package business

const (
	LbResourceTypeInternal = "internal" //内网负载均衡
	LbResourceTypeExternal = "external" //公网负载均衡

	AdminStatusDown   = "down"   //管理状态: DOWN
	AdminStatusActive = "active" //管理状态: ACTIVE

	// elb rule 状态
	ElbRuleStatusACTIVE = "ACTIVE"
	ElbRuleStatusDOWN   = "DOWN"

	// condition 类型
	ElbRuleConditionTypeServerName = "server_name"
	ElbRuleConditionTypeUrlPath    = "url_path"
	// 匹配类型
	ElbRuleMatchTypeABSOLUTE = "ABSOLUTE"
	ElbRuleMatchTypePREFIX   = "PREFIX"
	ElbRuleMatchTypeREG      = "REG"

	ElbTargetIPStatusOffline = "offline"
	ElbTargetIPStatusOnline  = "online"
	ElbTargetIPStatusUnknown = "unknown"

	ElbTargetTypeVM = "VM"
	ElbTargetTypeBM = "BM"

	ElbRuleActionTypeForward  = "forward"
	ElbRuleActionTypeRedirect = "redirect"
	ElbRuleActionTypeDeny     = "deny"

	ElbTargetInstanceTypeVM  = "VM"
	ElbTargetInstanceTypeBM  = "BM"
	ElbTargetInstanceTypeECI = "ECI"
	lbTargetInstanceTypeIP   = "IP"
)

var LbResourceType = []string{LbResourceTypeInternal, LbResourceTypeExternal}
var AdminStatusName = []string{AdminStatusActive, AdminStatusDown}
var ElbRuleStatus = []string{ElbRuleStatusACTIVE, ElbRuleStatusDOWN}
var ElbRuleConditionTypes = []string{ElbRuleConditionTypeServerName, ElbRuleConditionTypeUrlPath}
var ElbRuleMatchTypes = []string{ElbRuleMatchTypeABSOLUTE, ElbRuleMatchTypePREFIX, ElbRuleMatchTypeREG}
var ElbTargetIpStatus = []string{ElbTargetIPStatusOffline, ElbTargetIPStatusOnline, ElbTargetIPStatusUnknown}
var ElbTargetType = []string{ElbTargetTypeVM, ElbTargetTypeBM}
var ElbRuleActionType = []string{ElbRuleActionTypeForward, ElbRuleActionTypeRedirect, ElbRuleActionTypeDeny}
var ElbTargetInstanceType = []string{ElbTargetInstanceTypeVM, ElbTargetInstanceTypeBM, ElbTargetInstanceTypeECI, lbTargetInstanceTypeIP}
