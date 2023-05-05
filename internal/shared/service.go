// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package shared

import (
	"os"
	"path"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	pwg "github.com/sethvargo/go-password/password"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.breu.io/ctrlplane/internal/shared/logger"
)

type (
	service struct {
		Name    string `env:"SERVICE_NAME" env-default:"service"`
		Debug   bool   `env:"DEBUG" env-default:"false"`
		Secret  string `env:"SECRET" env-default:""`
		CLI     *cli   `env-prefix:"CTRLPLANE_" env-allow-empty:"true"`
		version string `env:"VERSION" env-default:""`
	}

	cli struct {
		BaseURL string `env:"BASE_URL" env-default:"http://api.ctrlplane.ai"`
		APIKEY  string `env:"API_KEY" env-default:""`
	}

	EchoValidator struct {
		Validator *validator.Validate
	}
)

var (
	Logger    *logger.ZapAdapter
	Service   = &service{}
	Validator *validator.Validate
)

// Version creates the version string as per [calver].
//
// The scheme currently being followed is YYYY.0M.0D.<git commit hash>-<channel> where:
//   - YYYY.0M.0D is the date of the commit
//   - <git commit hash> is the first 8 characters of the git commit hash
//   - <channel> is the channel of the build (e.g. dev, alpha, beta, rc, stable).
//
// For out purposes, -<channel> is optional and will be set to "dev" if the git is dirty.
//
// [calver]: https://calver.org/
func (s *service) Version() string {
	if s.version == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			var (
				revision  string
				modified  string
				timestamp time.Time
				version   string
			)

			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					revision = setting.Value
				}

				if setting.Key == "vcs.modified" {
					modified = setting.Value
				}

				if setting.Key == "vcs.time" {
					timestamp, _ = time.Parse(time.RFC3339, setting.Value)
				}
			}

			if len(revision) > 0 && len(modified) > 0 && timestamp.Unix() > 0 {
				version = timestamp.Format("2006.01.02") + "." + revision[:8]
			} else {
				version = "debug"
			}

			if modified == "true" {
				version += "-dev"
			}

			s.version = version
		}
	}

	return s.version
}

// ReadEnv reads the environment variables and initializes the service.
func (s *service) ReadEnv() {
	if err := cleanenv.ReadEnv(s); err != nil {
		panic("Failed to read environment variables")
	}
}

func (s *service) GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".ctrlplane", "config.json"), nil
}

func (s *service) ReadConfig() error {
	conf, err := s.GetConfigPath()

	if err != nil {
		return err
	}

	return cleanenv.ReadConfig(conf, s)
}

// InitValidator sets up global validator.
func (s *service) InitValidator() { InitValidator() }

// InitLogger sets up global logger.
func (s *service) InitLogger() {
	var zl *zap.Logger

	if s.Debug {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.FunctionKey = "func"
		zl, _ = config.Build()
	} else {
		println("production logger")
		config := zap.NewProductionConfig()
		config.EncoderConfig.FunctionKey = "func"
		zl, _ = config.Build()
	}

	Logger = logger.NewZapAdapter(zl)
}

// InitCLI sets up ctrlplane-cli.
func (s *service) InitCLI() error {
	s.Name = "ctrlplane-cli"
	s.CLI = &cli{}

	if err := cleanenv.ReadEnv(s); err != nil {
		return err
	}

	if s.Version() == "debug" {
		println("Warning: You are using a development version of the CLI. Please use the stable version.")

		s.Debug = true
		s.CLI.BaseURL = "http://localhost:8080"
	}

	s.InitLogger()

	return nil
}

func (ev *EchoValidator) Validate(i interface{}) error {
	if err := ev.Validator.Struct(i); err != nil {
		return err
	}

	return nil
}

func InitValidator() {
	Validator = validator.New()
	// by default, the validator will try to get json tag.
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func InitForTest() {
	Service = &service{
		Name:   "test",
		Secret: pwg.MustGenerate(32, 8, 0, false, false),
		Debug:  true,
	}

	Service.InitLogger()
	Service.InitValidator()
}
