package main

import (
	"embed"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog"
	"github.com/zaffka/iotas/pkg/gen"
	"github.com/zaffka/iotas/pkg/parse"
)

//go:embed iotas.tpl
var iotasTpl embed.FS

const (
	typeSeparator = ","

	templateNameMsg = "template_name"
	templateName    = "iotas.tpl"

	appVersionMsg = "app_ver"
	appVersion    = "dev"

	dirNameMsg = "dir"
	dirName    = "."
)

func main() {
	log := initLogger()

	types := flag.String("type", "", "-type=TypeName1,TypeName2 (at least one TypeName required)")
	flag.Parse()

	if len(*types) == 0 {
		log.Error().Msg("flag -type must be set and have at least one TypeName")

		return
	}

	typeNames := strings.Split(*types, typeSeparator)
	workDir := dirName

	if dirArg := flag.Arg(0); dirArg != "" {
		workDir = dirArg
	}

	dir, err := filepath.Abs(workDir)
	if err != nil {
		log.Error().Err(err).Str(dirNameMsg, workDir).Msg("failed to get an absolute filepath for the directory")

		return
	}

	log.Info().
		Str(dirNameMsg, workDir).
		Str(appVersionMsg, appVersion).
		Msg("Start parsing a directory")

	parser, err := parse.NewParser(dir, typeNames)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set-up a parser")

		return
	}

	if err := parser.Parse(); err != nil {
		log.Error().Err(err).Msg("Parsing finished with errors")

		return
	}

	tpl, err := template.ParseFS(iotasTpl, templateName)
	if err != nil {
		log.Error().Err(err).Str(templateNameMsg, templateName).Msg("Failed to parse a template")

		return
	}

	gen.Generator{
		AppVersion: appVersion,
		DirName:    dir,
		PkgName:    parser.GetPackageName(),
		Data:       parser.GetConstantsByType(),
		Tpl:        tpl,
		Logger:     log,
	}.Exec()
}

func initLogger() zerolog.Logger {
	outFmt := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
	}

	return zerolog.New(outFmt).With().Timestamp().Logger()
}
