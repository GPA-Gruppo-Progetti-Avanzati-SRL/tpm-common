package ghostscript_test

import (
	"fmt"
	"github.com/mario-imperato/tpm-common/util/ghostscript"
	"github.com/stretchr/testify/require"

	"io/ioutil"
	"os"
	"testing"
)

const (
	testDocumentSigned  = "/Users/marioa.imperato/projects/poste/gect/doc-poste-ntt-rev.1.1/doc-firmati/NFEA_1701140920069300A7052325157.pdf.pdf"
	testDocumentCH7     = "/Users/marioa.imperato/projects/poste/gect/doc-poste-ntt-rev.1.1/doc-firmati/NFEC_1501021100047000A3015356006_ch7pg.pdf"
	testOutFileTemplate = "flattened_%02d_%%02d.pdf"
)

func TestGhostscripter(t *testing.T) {

	_ = os.Setenv(ghostscript.GhostScriptWorkFolderEnvVarName, "/tmp/temp-ghs")
	gsOpts := []ghostscript.GsOption{
		ghostscript.WithWorkFolder("/tmp/temp-ghs"),
	}

	fileContent, err := ioutil.ReadFile(testDocumentCH7)
	require.NoError(t, err)

	gs, err := ghostscript.NewGhostscripter(gsOpts...)
	require.NoError(t, err)

	for i := 0; i < 1; i++ {
		cmdOpts := []ghostscript.CmdOption{
			ghostscript.WithDevice(ghostscript.DevicePDFWrite),
			ghostscript.WithNoInteractive(),
			ghostscript.WithQuiet(),
			ghostscript.WithOutFilename(fmt.Sprintf(testOutFileTemplate, i)),
			ghostscript.WithPdfSettings(ghostscript.PDFSettingsPrinter),
			ghostscript.WithFile(testDocumentSigned),
		}

		cmd, err := gs.NewCommand(cmdOpts...)
		require.NoError(t, err)

		err = gs.ExecuteCommand(cmd)
		require.NoError(t, err)
	}

	cmdOpts := []ghostscript.CmdOption{
		ghostscript.WithDevice(ghostscript.DeviceJPEG),
		ghostscript.WithNoInteractive(),
		ghostscript.WithQuiet(),
		ghostscript.WithJPEGQuality(20),
		ghostscript.WithOutputResolution(150),
		ghostscript.WithKeepTempFiles(false),
	}

	cmd, err := gs.NewCommand(cmdOpts...)
	require.NoError(t, err)

	_, err = gs.ExecuteCommandOnBlockData(cmd, fileContent)
	require.NoError(t, err)

	cmdOpts = []ghostscript.CmdOption{
		ghostscript.WithDevice(ghostscript.DevicePDFWrite),
		ghostscript.WithNoInteractive(),
		ghostscript.WithQuiet(),
		ghostscript.WithJPEGQuality(20),
		ghostscript.WithOutputResolution(150),
		ghostscript.WithKeepTempFiles(false),
		ghostscript.WithSingleOutFile(),
	}

	cmd, err = gs.NewCommand(cmdOpts...)
	require.NoError(t, err)

	_, err = gs.ExecuteCommandOnBlockData(cmd, fileContent)
	require.NoError(t, err)
}
