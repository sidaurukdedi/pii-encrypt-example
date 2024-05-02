package hook

import (
	"os"

	"pii-encrypt-example/entity"

	"github.com/sirupsen/logrus"
)

type stdoutLoggerHook struct {
	logger *logrus.Logger
}

func NewStdoutLoggerHook(logger *logrus.Logger, formatter logrus.Formatter) logrus.Hook {
	logger.Out = os.Stdout
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(formatter)

	return &stdoutLoggerHook{logger: logger}
}

// Levels implements logrus.Hook interface, this hook applies to all defined levels
func (d *stdoutLoggerHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel}
}

// Fire implements logrus.Hook interface, attaches trace and span details found in entry context
func (d *stdoutLoggerHook) Fire(e *logrus.Entry) error {
	var clientDevice entity.ClientDevice
	e.Logger = d.logger
	ctx := e.Context
	if ctx == nil {
		return nil
	}

	clientDevice, ok := ctx.Value(entity.ClientContextKey{}).(entity.ClientDevice)
	if !ok {
		return nil
	}

	e.Data["mpv.client.user_agent"] = clientDevice.UserAgent
	e.Data["mpv.client.remote_address"] = clientDevice.RemoteAddress
	e.Data["mpv.client.x_forwarded_for"] = clientDevice.XForwardedFor
	e.Data["mpv.client.x_real_ip"] = clientDevice.XRealIP

	return nil
}
