module github.com/seanbit/gogateway

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/seanbit/ginserver v1.0.3
	github.com/seanbit/gokit v1.0.3
	github.com/seanbit/goserving v1.0.6
	github.com/sirupsen/logrus v1.4.2
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.0
