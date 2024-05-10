package logs

import (
	"fmt"
	"main/pkg/constants"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ApplicationLog *zap.Logger

func SetUpApplicationLogs() {
	config := zap.NewProductionConfig()

	// Set log file path
	logFilePath := constants.ApplicationConfig.Application.LogPath + "messagechannelBussinessLogic.log"
	config.OutputPaths = []string{logFilePath}

	// Create lumberjack logger for log rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100, // MB
		MaxAge:     28,  // Days
		MaxBackups: 30000,
		LocalTime:  true,
		Compress:   true,
	}

	// Add lumberjack logger as output
	// config.OutputPaths = append(config.OutputPaths, "stdout")           // Also log to stdout
	config.ErrorOutputPaths = append(config.ErrorOutputPaths, "stderr") // Log errors to stderr
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // Use ISO8601 time format
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	// config.Callers = zapcore.ShortCaller

	// Create Zap core with log rotation
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		zapcore.AddSync(lumberjackLogger),
		zap.NewAtomicLevelAt(zap.InfoLevel), // Change log level here if needed
		// zap.RegisterEncoder("",zapcore.ShortCallerEncoder)
	)

	// Create logger with the configured core
	logger := zap.New(core, zap.AddCaller())

	ApplicationLog = logger

	// Close the lumberjack logger when application exits
	defer lumberjackLogger.Close()
}

func InfoLog(format string, a ...interface{}) {
	stringMessage := fmt.Sprintf(format, a...)
	ApplicationLog.Info(stringMessage)
}

func ErrorLog(format string, a ...interface{}) {
	stringMessage := fmt.Sprintf(format, a...)
	ApplicationLog.Error(stringMessage)
}
