package log

import (
	"reflect"
	"github.com/go-ozzo/ozzo-config"
	"github.com/go-ozzo/ozzo-log"
)
var (
	Logger *log.Logger
)

func NewLogger(c *config.Config) {
	targets := c.Get("logger.Targets")
	if targets != nil {
		Logger = log.NewLogger()
		if _targets, ok := targets.([]interface{}); ok {
			for _, target := range _targets {
				//获取type
				value := reflect.ValueOf(target)
				switch value.Kind() {
				case reflect.Map:
					//_type := c.GetString("components.logger.Targets."+strconv.Itoa(key)+".type")
					_type := value.MapIndex(reflect.ValueOf("type")).Interface().(string)
					switch _type {
					case "ConsoleTarget":
						if err := c.Register(_type, log.NewConsoleTarget); err != nil {
							panic(err)
						}
					case "FileTarget":
						if err := c.Register(_type, log.NewFileTarget); err != nil {
							panic(err)
						}
					case "mailTarget":
						if err := c.Register(_type, log.NewMailTarget); err != nil {
							panic(err)
						}
					case "networkTarget":
						if err := c.Register(_type, log.NewNetworkTarget); err != nil {
							panic(err)
						}
					default:
						panic("the target of logger is not validate:" + _type)
					}
				default:
					panic("target type is undeclare")
				}
			}

			if err := c.Configure( Logger,"logger"); err != nil {
				panic(err)
			}

			 Logger.Open()
		}
	}
}

func Info(category string, format string, a ...interface{}) {
	logger :=  Logger.GetLogger(category)
	logger.Info(format, a ...)
}

func Error(category string, format string, a ...interface{}) {
	logger :=  Logger.GetLogger(category)
	logger.Error(format, a ...)
}

func Debug(category string, format string, a ...interface{}) {
	logger :=  Logger.GetLogger(category)
	logger.Debug(format, a ...)
}

func Warning(category string, format string, a ...interface{}) {
	logger :=  Logger.GetLogger(category)
	logger.Warning(format, a ...)
}

func Emergency(category string, format string, a ...interface{}) {
	logger :=  Logger.GetLogger(category)
	logger.Emergency(format, a ...)
}

func Alert(category string, format string, a ...interface{}) {
	logger :=  Logger.GetLogger(category)
	logger.Alert(format, a ...)
}
func Critical(category string, format string, a ...interface{}) {
	logger :=  Logger.GetLogger(category)
	logger.Critical(format, a ...)
}
func Notice(category string, format string, a ...interface{}) {
	logger :=  Logger.GetLogger(category)
	logger.Notice(format, a ...)
}
