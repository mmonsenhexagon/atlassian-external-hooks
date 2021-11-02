package bitbucket

import (
	"strings"
	"time"

	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

type Volume string

func (volume Volume) Start(
	version string,
	opts StartOpts,
) (*Bitbucket, error) {
	if opts.PortHTTP == 0 {
		opts.PortHTTP = 7990
	}

	if opts.PortSSH == 0 {
		opts.PortSSH = 7999
	}

	if opts.AdminUser == "" {
		opts.AdminUser = "admin"
	}

	if opts.AdminPassword == "" {
		opts.AdminPassword = "admin"
	}

	instance := &Instance{
		version: version,
		volume:  string(volume),
	}

	instance.opts.StartOpts = opts

	if opts.ContainerID != "" {
		log.Infof(
			karma.
				Describe("container", opts.ContainerID).
				Describe("volume", string(volume)).
				Describe("opts", opts),
			"{bitbucket %s} re-using existing container",
			version,
		)

		instance.container = opts.ContainerID
	} else {
		log.Infof(
			karma.
				Describe("volume", string(volume)).
				Describe("opts", opts),
			"{bitbucket %s} starting container",
			version,
		)

		err := instance.start()
		if err != nil {
			return nil, karma.Format(
				err,
				"start bitbicket container",
			)
		}
	}

	err := instance.connect()
	if err != nil {
		return nil, karma.Format(
			err,
			"connect to container",
		)
	}

	instance.stacktraceLogs, err = instance.startLogReader(false)
	if err != nil {
		return nil, karma.Format(err, "start log reader")
	}

	instance.testcaseLogs, err = instance.startLogReader(true)
	if err != nil {
		return nil, karma.Format(err, "start log reader")
	}

	var message string

	for {
		status, err := instance.getStartupStatus()
		if err != nil {
			return nil, karma.Format(
				err,
				"get container startup status",
			)
		}

		if status == nil {
			continue
		}

		if message != status.Progress.Message {
			log.Debugf(
				nil,
				"{bitbucket %s} setup: %3d%% %s | %s",
				version,
				status.Progress.Percentage,
				strings.ToLower(status.State),
				status.Progress.Message,
			)

			message = status.Progress.Message
		}

		if status.State == "STARTED" {
			break
		}

		time.Sleep(time.Millisecond * 20)
	}

	return New(instance)
}
