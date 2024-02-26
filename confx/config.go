package confx

import (
	"os"
	"strings"

	"code.sikey.com.cn/luoyangwei/serverx/dbx"
	"code.sikey.com.cn/luoyangwei/serverx/gid"
	"code.sikey.com.cn/luoyangwei/serverx/rdbx"
	"code.sikey.com.cn/luoyangwei/serverx/zlog"
	"github.com/rotisserie/eris"
	"github.com/spf13/viper"
)

// SetEnvConfig 将环境变量设置到默认值
func SetEnvConfig(inp ...string) {
	inp = append(inp, nodeId)
	for _, in := range inp {
		viper.SetDefault(in, os.Getenv(in))
	}
}

type LoadOption func() error

// defaultConfigFileType 默认的配置文件格式
const defaultConfigFileType = "toml"

const (
	nodeId = "NODE_ID"
)

// LoadConfig 加载配置
func LoadConfig(file string, opts ...LoadOption) error {
	viper.SetConfigFile(file)
	viper.SetConfigType(defaultConfigFileType)
	if err := viper.ReadInConfig(); err != nil {
		panic(eris.Wrap(err, "无法加载配置"))
	}

	var err error
	for _, opt := range opts {
		if err = opt(); err != nil {
			return err
		}
	}

	gid.SetNodeId(viper.GetInt64(nodeId))
	if containsKey(zlog.KeyPrefix) {
		zlog.NewZlog()
	}
	if containsKey(dbx.KeyPrefix) {
		dbx.Connect()
	}
	if containsKey(dbx.KeyPrefix) {
		rdbx.Connect()
	}
	return nil
}

func GetServerPort() int {
	return viper.GetInt("port")
}

func GetServerName() string {
	return viper.GetString("name")
}

func containsKey(k string) bool {
	var keys = viper.AllKeys()
	for _, key := range keys {
		var prefix = key
		if strings.Contains(key, ".") {
			prefix = strings.Split(key, ".")[0]
		}
		if prefix == k {
			return true
		}
	}
	return false
}
