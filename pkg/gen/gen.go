package gen

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/rs/zerolog"
)

const (
	filePermission = 0600
	suffix         = "_iotas.go"
	outFileMsg     = "out_file"
)

type Generator struct {
	AppVersion string
	DirName    string
	PkgName    string

	Data map[string][]string
	Tpl  *template.Template

	Logger zerolog.Logger
}

type tplData struct {
	AppVersion string
	PkgName    string
	TypeName   string
	ConstNames []string
}

func (g Generator) Exec() {
	wg := &sync.WaitGroup{}

	for typeName, constantNames := range g.Data {
		if len(constantNames) == 0 {
			continue
		}

		wg.Add(1)

		go g.worker(typeName, constantNames, wg)
	}

	wg.Wait()
}

func (g *Generator) worker(typeName string, constantNames []string, wg *sync.WaitGroup) {
	defer wg.Done()

	log := g.Logger.With().
		Str("type_name", typeName).
		Int("constants_count", len(constantNames)).
		Logger()

	var buf bytes.Buffer
	err := g.Tpl.Execute(&buf, tplData{
		AppVersion: g.AppVersion,
		PkgName:    g.PkgName,
		TypeName:   typeName,
		ConstNames: constantNames,
	})

	if err != nil {
		log.Err(err).Msg("Failed to execute a template")

		return
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Err(err).Msg("Failed to format the generated source code")

		return
	}

	fileName := strings.ToLower(typeName + suffix)
	filePath := filepath.Join(g.DirName, fileName)

	if err := os.WriteFile(filePath, src, filePermission); err != nil {
		log.Error().Err(err).Str(outFileMsg, filePath).Msg("Failed to write the out file")

		return
	}

	log.Info().Str(outFileMsg, filePath).Msg("The out file has been written")
}
