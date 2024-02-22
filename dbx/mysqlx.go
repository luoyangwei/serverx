package dbx

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	sql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const KeyPrefix = "mysql"

type mysqlConfig struct {
	SSH                       bool          `toml:"ssh"` // SSH 是否开启SSH
	Dsn                       string        // Dsn 数据源地址
	SkipDefaultTransaction    bool          // SkipDefaultTransaction 跳过默认事务
	SlowThreshold             time.Duration // SlowThreshold 慢 SQL 阈值
	IgnoreRecordNotFoundError bool          // IgnoreRecordNotFoundError 忽略记录未找到的错误
	MaxLifetime               time.Duration // MaxLifetime 连接的有效时长
	MaxOpenConns              int           // MaxOpenConns 打开数据库连接的最大数量。
	MaxIdleConns              int           // MaxIdleConns 空闲连接池中连接的最大数量
}

type driverConfig struct {
	username string
	password string
	protocol string
	address  string
	port     int
	db       string
	params   string
}

func (vc *driverConfig) formatDSN() string {
	return vc.username + ":" + vc.password + "@" +
		vc.protocol + "(" + vc.address + ":" + strconv.Itoa(vc.port) + ")/" + vc.db + "?" + vc.params
}

type Dialer struct {
	client *ssh.Client
}

func (v *Dialer) Dial(ctx context.Context, address string) (net.Conn, error) {
	return v.client.Dial("tcp", address)
}

// Connect 初始化 mysql 连接
func Connect() *gorm.DB {
	cfg := readMysqlConfig()
	conn, err := gorm.Open(sql.New(sql.Config{
		DSN:                       cfg.Dsn,
		DefaultStringSize:         255,
		SkipInitializeWithVersion: false,
	}), getGormConfig(cfg))
	if err != nil {
		log.Panicln(err)
	}

	sqlDB, _ := conn.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	log.Printf("Mysql connected to %s \n", cfg.Dsn)
	return conn
}

func withDsn(dsn string) *driverConfig {
	// root:qq123123@tcp(127.0.0.1:3306)/sikey?charset=utf8mb4&parseTime=true&loc=Local
	var user = strings.Split(dsn, "@")
	var ua = strings.Split(user[0], ":")
	var protocol = strings.Split(user[1], "(")
	var address = strings.Split(protocol[1], ")")
	var addr = strings.Split(address[0], ":")
	var port, err = strconv.ParseInt(addr[1], 10, 64)
	if err != nil {
		port = 3306
	}
	return &driverConfig{
		username: ua[0],
		password: ua[1],
		protocol: protocol[0],
		address:  addr[0],
		port:     int(port),
		db:       strings.Split(address[1], "?")[0][1:],
		params:   strings.Split(address[1], "?")[1],
	}
}

// readMysqlConfig 加载配置
func readMysqlConfig() mysqlConfig {
	var cfg mysqlConfig
	if err := viper.UnmarshalKey("mysql", &cfg); err != nil {
		log.Fatalln(err)
	}

	if cfg.SSH {
		config := &ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.Password("RHTUH2z49aEXnsgz"),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		var err error
		var clt *ssh.Client
		if clt, err = ssh.Dial("tcp", "106.75.230.4:22", config); err != nil {
			log.Fatalln(err)
		}

		var protocol = "ssh"
		vc := withDsn(cfg.Dsn)
		vc.protocol = protocol
		cfg.Dsn = vc.formatDSN()
		mysql.RegisterDialContext(protocol, (&Dialer{client: clt}).Dial)
	}
	return cfg
}

// getGormConfig 获取 gorm 配置
func getGormConfig(cfg mysqlConfig) *gorm.Config {
	return &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   cfg.SkipDefaultTransaction,
		Logger:                                   defaultLogger(cfg),
	}
}

// defaultLogger 默认的日志打印
func defaultLogger(cfg mysqlConfig) logger.Interface {
	return logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             cfg.SlowThreshold * time.Millisecond, // Slow SQL threshold
		LogLevel:                  logger.Silent,                        // Log level
		IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,        // Ignore ErrRecordNotFound error for logger
	})
}
