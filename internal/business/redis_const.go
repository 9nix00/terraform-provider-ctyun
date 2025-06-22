package business

const (
	RedisVersionBasic = "BASIC"
	RedisVersionPlus  = "PLUS"

	RedisDiskTypeSas = "SAS"
	RedisDiskTypeSsd = "SSD"

	RedisHostTypeS  = "S"
	RedisHostTypeC  = "C"
	RedisHostTypeM  = "M"
	RedisHostTypeHS = "HS"
	RedisHostTypeHC = "HC"
	RedisHostTypeKS = "KS"
	RedisHostTypeKC = "KC"

	RedisEditionStandardSingle          = "StandardSingle"
	RedisEditionStandardDual            = "StandardDual"
	RedisEditionDirectClusterSingle     = "DirectClusterSingle"
	RedisEditionDirectCluster           = "DirectCluster"
	RedisEditionClusterOriginalProxy    = "ClusterOriginalProxy"
	RedisEditionOriginalMultipleReadLvs = "OriginalMultipleReadLvs"
)

var RedisEngineVersion = []string{"5.0", "6.0", "7.0"}

var RedisHostType = []string{
	RedisHostTypeS,
	RedisHostTypeC,
	RedisHostTypeM,
	RedisHostTypeHS,
	RedisHostTypeHC,
	RedisHostTypeKS,
	RedisHostTypeKC,
}
