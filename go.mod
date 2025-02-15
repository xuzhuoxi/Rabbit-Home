module github.com/xuzhuoxi/Rabbit-Home

go 1.16

require (
	github.com/json-iterator/go v1.1.12
	github.com/xuzhuoxi/infra-go v1.0.4
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/xuzhuoxi/infra-go v1.0.4 => ../infra-go
)
