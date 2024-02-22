package zlog

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"code.sikey.com.cn/serverbackend/Serverx/format"
)

const KeyPrefix = "logger"

const (
	consoleLoggerEncoder = "console"
	jsonLoggerEncoder    = "json"
	defaultLoggerEncoder = consoleLoggerEncoder

	LoggerEncoderKey = "logger.encoder"
)

const (
	slsAccessKeyId     = "logger.sls.accessKeyId"
	slsAccessKeySecret = "logger.sls.accessKeySecret"
	slsEndpoint        = "logger.sls.endpoint"
	slsLogStoreName    = "logger.sls.logStoreName"
	slsProjectName     = "logger.sls.projectName"
	slsSourceIp        = "logger.sls.sourceIp"
	slsTopic           = "logger.sls.topic"
)

var lowPriority = zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
	return lv >= zapcore.DebugLevel
})

var logger *zap.SugaredLogger

func NewZlog() {
	// 控制台展示方便调试，使用 TEXT 的方式
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(format.DateParseAllUnixMilliFormat)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	stdCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), lowPriority)

	// 日志格式化
	productionEncoderConfig := zap.NewProductionEncoderConfig()
	productionEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(format.DateParseAllUnixMilliFormat)
	jsonEnc := zapcore.NewJSONEncoder(productionEncoderConfig)

	var core zapcore.Core
	encoder := viper.GetString(LoggerEncoderKey)
	if encoder == consoleLoggerEncoder {
		core = stdCore
	} else {
		syncer := zapcore.AddSync(NewSLSWriter(SLSWriterConfig{
			AccessKeyId:     viper.GetString(slsAccessKeyId),
			AccessKeySecret: viper.GetString(slsAccessKeySecret),
			Endpoint:        viper.GetString(slsEndpoint),
			ProjectName:     viper.GetString(slsProjectName),
			LogStoreName:    viper.GetString(slsLogStoreName),
			SourceIP:        viper.GetString(slsSourceIp),
			Topic:           viper.GetString(slsTopic),
		}))
		core = zapcore.NewTee(stdCore, zapcore.NewCore(jsonEnc, syncer, lowPriority))
	}

	logger = zap.New(core).WithOptions(zap.AddCallerSkip(1), zap.AddCaller()).Sugar()
}

type SLSWriter struct {
	cfg      SLSWriterConfig
	client   sls.ClientInterface
	project  *sls.LogProject
	logStore *sls.LogStore
}

type SLSWriterConfig struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string

	ProjectName  string
	LogStoreName string
	Topic        string
	SourceIP     string
	Environment  string
}

func NewSLSWriter(cfg SLSWriterConfig) *SLSWriter {
	// Create a logging service client.
	provider := sls.NewStaticCredentialsProvider(cfg.AccessKeyId, cfg.AccessKeySecret, "")
	client := sls.CreateNormalInterfaceV2(cfg.Endpoint, provider)

	project, err := client.GetProject(cfg.ProjectName)
	if err != nil {
		log.Fatalln(err)
	}
	logStore, err := client.GetLogStore(project.Name, cfg.LogStoreName)
	if err != nil {
		log.Fatalln(err)
	}
	return &SLSWriter{
		cfg:      cfg,
		client:   client,
		project:  project,
		logStore: logStore,
	}
}
func (w *SLSWriter) GetEnv() string {
	return w.cfg.Environment
}

func (w *SLSWriter) Write(p []byte) (int, error) {
	var err error
	var attributes map[string]string
	if err = json.Unmarshal(p, &attributes); err != nil {
		return 0, err
	}

	contents := make([]*sls.LogContent, 0)
	for k, v := range attributes {
		contents = append(contents, &sls.LogContent{Key: proto.String(k), Value: proto.String(v)})
	}

	unix := uint32(time.Now().UTC().Unix())
	logs := []*sls.Log{
		{
			Time:     proto.Uint32(unix),
			Contents: contents,
		},
	}
	lg := &sls.LogGroup{
		Topic:  proto.String(w.cfg.Topic),
		Source: proto.String(w.cfg.SourceIP),
		Logs:   logs,
	}
	return 1, w.client.PutLogs(w.project.Name, w.logStore.Name, lg)
}

const (
	requestId = "Request-Id"
)

func WithRequestId(ctx context.Context) *zap.SugaredLogger {
	rid := ctx.Value(requestId)
	if rid != nil && rid != "" {
		return logger.With(requestId, rid)
	}
	return logger
}

func DPanic(args ...interface{}) {
	logger.DPanic(args...)
}
func DPanicf(template string, args ...interface{}) {
	logger.DPanicf(template, args...)
}
func DPanicw(msg string, keysAndValues ...interface{}) {
	logger.DPanicw(msg, keysAndValues...)
}
func Debug(args ...interface{}) {
	logger.Debug(args...)
}
func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}
func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}
func Desugar() *zap.Logger {
	return logger.Desugar()
}
func Error(args ...interface{}) {
	logger.Error(args...)
}
func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}
func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}
func Fatalw(msg string, keysAndValues ...interface{}) {
	logger.Fatalw(msg, keysAndValues...)
}
func Info(args ...interface{}) {
	logger.Info(args...)
}
func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}
func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}
func Named(name string) *zap.SugaredLogger {
	return logger.Named(name)
}
func Panic(args ...interface{}) {
	logger.Panic(args...)
}
func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}
func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Panicw(msg, keysAndValues...)
}
func Sync() error {
	return logger.Sync()
}
func Warn(args ...interface{}) {
	logger.Warn(args...)
}
func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}
func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}
func With(args ...interface{}) *zap.SugaredLogger {
	return logger.With(args...)
}
