package common

import (
	"sync"
	ccse2 "terraform-provider-ctyun/internal/core/ccse"
	"terraform-provider-ctyun/internal/core/core"
	"terraform-provider-ctyun/internal/core/ctebm"
	ctebs2 "terraform-provider-ctyun/internal/core/ctebs"
	ctecs2 "terraform-provider-ctyun/internal/core/ctecs"
	sdkCtelb "terraform-provider-ctyun/internal/core/ctelb"
	ctvpc2 "terraform-provider-ctyun/internal/core/ctvpc"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctebs"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctimage"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"terraform-provider-ctyun/internal/core/ctzos"
)

var once sync.Once
var ctyunMetadata *CtyunMetadata

type CtyunMetadata struct {
	Apis          *Apis
	Credential    ctyunsdk.Credential
	extra         map[string]string
	SdkCredential core.Credential
}

// InitCtyunMetadata 初始化
func InitCtyunMetadata(apis *Apis, credential ctyunsdk.Credential, sdkCred core.Credential, extra map[string]string) {
	ctyunMetadata = &CtyunMetadata{Apis: apis, Credential: credential, SdkCredential: sdkCred, extra: extra}
}

// AcquireCtyunMetadata 获取实例对象
func AcquireCtyunMetadata() *CtyunMetadata {
	if ctyunMetadata == nil {
		panic("ctyun metadata not init!")
	}
	return ctyunMetadata
}

// GetExtra 获取默认设置的值
func (c CtyunMetadata) GetExtra(extraKey string) string {
	return c.extra[extraKey]
}

// GetExtraIfEmpty 如果目标值为空，获取默认设置的值，若目标值非空则返回目标值
func (c CtyunMetadata) GetExtraIfEmpty(target, extraKey string) string {
	if target == "" {
		return c.extra[extraKey]
	}
	return target
}

type Apis struct {
	CtEbsApis      *ctebs.Apis
	CtEcsApis      *ctecs.Apis
	CtIamApis      *ctiam.Apis
	CtImageApis    *ctimage.Apis
	CtVpcApis      *ctvpc.Apis
	CtEbmApis      *ctebm.Apis
	SdkCtEbsApis   *ctebs2.Apis
	SdkCtEcsApis   *ctecs2.Apis
	SdkCtVpcApis   *ctvpc2.Apis
	SdkCtZosApis   *ctzos.Apis
	SdkCcseApis    *ccse2.Apis
	SdkCtElbApis   *sdkCtelb.Apis
	SdkCtMysqlApis *mysql.Apis
}
