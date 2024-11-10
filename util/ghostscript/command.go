package ghostscript

import (
	"errors"
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
	"github.com/rs/zerolog/log"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DefaultFilePrefix          = "gser"
	TempInFilePattern          = "%s-%08d-in-file"
	TempMultipleOutFilePattern = "%s-%08d-out-file-%%d.%s"
	TempSingleOutFilePattern   = "%s-%08d-out-file.%s"
	PDFSettings                = "PDFSETTINGS"
	PDFSettingsSwitch          = "-dPDFSETTINGS="
	PDFSettingsScreen          = "/screen"
	PDFSettingsPrinter         = "/printer"
	PDFSettingsPrePress        = "/prepress"
	PDFSettingsDefault         = "/default"
	Device                     = "DEVICE"
	DeviceSwitch               = "-sDEVICE="
	DevicePDFWrite             = "pdfwrite"
	DeviceJPEG                 = "jpeg"
	Batch                      = "BATCH"
	BatchSwitch                = "-dBATCH"
	NoPause                    = "NOPAUSE"
	NoPauseSwitch              = "-dNOPAUSE"
	Quiet                      = "QUIET"
	QuietSwitch                = "-q"
	JpegQuality                = "JPEGQ"
	JpegQualitySwitch          = "-dJPEGQ="
	OutputResolution           = "OUTPUT-RES"
	OutputResolutionSwitch     = "-r"
	InFile                     = "INFILE"
	InFileSwitch               = "-f"
	OutFile                    = "OUTFILE"
	OutFileSwitch              = "-sOutputFile="
)

type CmdOutputData struct {
	Ct []byte
	Fn string
}

type CmdArgument struct {
	argId     string
	argSwitch string
	argValue  string
}

func (ca CmdArgument) String() string {
	return ca.argSwitch + ca.argValue
}

type GsCmd struct {
	cmdId                  int
	verbose                bool
	workFolder             string
	keepTempFiles          bool
	numExpectedOutputFiles int
	args                   map[string]CmdArgument
	files                  []string
	singleOutputFile       bool
}

func (cmd *GsCmd) InputTempFileName() string {
	return filepath.Join(cmd.workFolder, fmt.Sprintf(TempInFilePattern, DefaultFilePrefix, cmd.cmdId))
}

func (cmd *GsCmd) OutputTempFilePattern() (string, error) {

	outfilePattern := TempMultipleOutFilePattern
	if cmd.singleOutputFile {
		outfilePattern = TempSingleOutFilePattern
	}

	if len(cmd.args) == 0 {
		return "", errors.New("")
	}

	dev, ok := cmd.args[Device]
	if !ok {
		return "", errors.New("cannot compute out-file-pattern device is undefined")
	}

	ext := ""
	switch dev.argValue {
	case DevicePDFWrite:
		ext = "pdf"
	case DeviceJPEG:
		ext = "jpg"
	default:

	}

	outPattern := fmt.Sprintf(outfilePattern, DefaultFilePrefix, cmd.cmdId, ext)
	return outPattern, nil
}

type CmdOption func(cfg *GsCmd)

func WithSingleOutFile() CmdOption {
	return func(cfg *GsCmd) {
		cfg.singleOutputFile = true
	}
}

func WithVerbose(b bool) CmdOption {
	return func(cfg *GsCmd) {
		cfg.verbose = b
	}
}

func WithCommandId(i int) CmdOption {
	return func(cfg *GsCmd) {
		cfg.cmdId = i
	}
}

func WithNumExpectedOutputFiles(n int) CmdOption {
	return func(cfg *GsCmd) {
		cfg.numExpectedOutputFiles = n
	}
}

func WithKeepTempFiles(f bool) CmdOption {
	return func(cfg *GsCmd) {
		cfg.keepTempFiles = f
	}
}

func WithFile(f string) CmdOption {
	return func(cfg *GsCmd) {
		cfg.files = append(cfg.files, f)
	}
}

func WithPdfSettings(s string) CmdOption {
	return func(cfg *GsCmd) {
		if len(cfg.args) == 0 {
			cfg.args = make(map[string]CmdArgument)
		}

		cfg.args[PDFSettings] = CmdArgument{argId: PDFSettings, argSwitch: PDFSettingsSwitch, argValue: s}
	}
}

func WithDevice(d string) CmdOption {
	return func(cfg *GsCmd) {
		if len(cfg.args) == 0 {
			cfg.args = make(map[string]CmdArgument)
		}

		cfg.args[Device] = CmdArgument{argId: Device, argSwitch: DeviceSwitch, argValue: d}
	}
}

func WithNoInteractive() CmdOption {
	return func(cfg *GsCmd) {
		if len(cfg.args) == 0 {
			cfg.args = make(map[string]CmdArgument)
		}

		cfg.args[Batch] = CmdArgument{argId: Batch, argSwitch: BatchSwitch}
		cfg.args[NoPause] = CmdArgument{argId: NoPause, argSwitch: NoPauseSwitch}
	}
}

func WithQuiet() CmdOption {
	return func(cfg *GsCmd) {
		if len(cfg.args) == 0 {
			cfg.args = make(map[string]CmdArgument)
		}

		cfg.args[Quiet] = CmdArgument{argId: Quiet, argSwitch: QuietSwitch}
	}
}

func WithOutputResolution(r int) CmdOption {
	return func(cfg *GsCmd) {
		if len(cfg.args) == 0 {
			cfg.args = make(map[string]CmdArgument)
		}

		cfg.args[OutputResolution] = CmdArgument{argId: OutputResolution, argSwitch: OutputResolutionSwitch, argValue: strconv.Itoa(r)}
	}
}

func WithJPEGQuality(r int) CmdOption {
	return func(cfg *GsCmd) {
		if len(cfg.args) == 0 {
			cfg.args = make(map[string]CmdArgument)
		}

		cfg.args[JpegQuality] = CmdArgument{argId: JpegQuality, argSwitch: JpegQualitySwitch, argValue: strconv.Itoa(r)}
	}
}

// WithInFilename This is dangerous since it must be put at the end of params. Better is to use the WithFile equivalent that uses a
// different syntax on the command line and accepts multiple values.
func WithInFilename(fn string) CmdOption {
	return func(cfg *GsCmd) {
		if len(cfg.args) == 0 {
			cfg.args = make(map[string]CmdArgument)
		}

		cfg.args[InFile] = CmdArgument{argId: InFile, argSwitch: InFileSwitch, argValue: fn}
	}
}

func WithOutFilename(fn string) CmdOption {
	return func(cfg *GsCmd) {
		if len(cfg.args) == 0 {
			cfg.args = make(map[string]CmdArgument)
		}

		cfg.args[OutFile] = CmdArgument{argId: OutFile, argSwitch: OutFileSwitch, argValue: fn}
	}
}

func (cmd *GsCmd) BuildArgList() ([]string, error) {
	args := make([]string, 0, 10)
	for _, a := range cmd.args {
		cmdArg := a.String()
		switch a.argId {
		case OutFile:
			fallthrough
		case InFile:
			if a.argValue == "" {
				return nil, fmt.Errorf("missing value for %s switch", a.argId)
			}
			if filepath.Dir(a.argValue) == "." {
				cmdArg = a.argSwitch + filepath.Join(cmd.workFolder, a.argValue)
			}
		default:
		}

		args = append(args, cmdArg)
	}

	for _, f := range cmd.files {
		args = append(args, f)
	}

	if cmd.verbose {
		for i := range args {
			log.Trace().Str("arg", args[i]).Msg("command arguments")
		}
	}

	return args, nil
}

func (gs *Ghostscripter) NewCommand(opts ...CmdOption) (*GsCmd, error) {

	cfg := &GsCmd{workFolder: gs.cfg.workFolderMountPoint, numExpectedOutputFiles: -1}
	for _, o := range opts {
		o(cfg)
	}

	return cfg, nil
}

func (gs *Ghostscripter) ExecuteCommand(cfg *GsCmd) error {

	args, err := cfg.BuildArgList()
	if err != nil {
		return err
	}

	cmd := exec.Command(gs.cfg.whichGs, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return err
}

func (gs *Ghostscripter) ExecuteCommandOnBlockData(cfg *GsCmd, b []byte) ([]CmdOutputData, error) {

	outFormat, err := cfg.OutputTempFilePattern()
	if err != nil {
		return nil, err
	}

	inputFileName := cfg.InputTempFileName()
	err = ioutil.WriteFile(inputFileName, b, fs.ModePerm)
	if err != nil {
		return nil, err
	}

	additionalOpts := []CmdOption{
		WithOutFilename(outFormat),
		WithFile(cfg.InputTempFileName()),
	}

	for _, o := range additionalOpts {
		o(cfg)
	}

	err = gs.ExecuteCommand(cfg)
	if err != nil {
		_ = clearFiles(!cfg.keepTempFiles, inputFileName, nil)
		return nil, err
	}

	fileNames, err := collectOutputFiles(gs.cfg.workFolderMountPoint, outFormat, cfg.numExpectedOutputFiles)
	if err != nil {
		_ = clearFiles(!cfg.keepTempFiles, inputFileName, fileNames)
		return nil, err
	}

	numFiles2Read := len(fileNames)
	if cfg.numExpectedOutputFiles > 0 {
		numFiles2Read = cfg.numExpectedOutputFiles
	}

	resp := make([]CmdOutputData, 0)
	for i := 0; i < numFiles2Read; i++ {
		fn := fileNames[i]
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			_ = clearFiles(!cfg.keepTempFiles, inputFileName, fileNames)
			return nil, err
		}

		resp = append(resp, CmdOutputData{Fn: fn, Ct: b})
	}

	_ = clearFiles(!cfg.keepTempFiles, inputFileName, fileNames)

	return resp, nil
}

func clearFiles(doClear bool, inFn string, outFn []string) error {

	if !doClear {
		return nil
	}

	var firstErr error
	err := os.Remove(inFn)
	if err != nil {
		log.Warn().Err(err).Str("fn", inFn).Msg("cannot remove file of ghostscript input")
		firstErr = err
	}

	for _, fn := range outFn {
		err := os.Remove(fn)
		if err != nil {
			log.Warn().Err(err).Str("fn", fn).Msg("cannot remove file of ghostscript output")
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}

func collectOutputFiles(workDir string, filePattern string, minExpected int) ([]string, error) {

	files := make([]string, 0)
	if strings.Contains(filePattern, "%d") || strings.Contains(filePattern, "%02d") {
		keepLooping := true
		fileNum := 0
		for keepLooping {
			fileNum++
			fn := filepath.Join(workDir, fmt.Sprintf(filePattern, fileNum))
			if fileutil.FileExists(fn) {
				files = append(files, fn)
			} else {
				keepLooping = false
			}
		}
	} else {
		fn := filepath.Join(workDir, filePattern)
		if fileutil.FileExists(fn) {
			files = append(files, fn)
		}
	}

	if minExpected > 0 && len(files) < minExpected {
		return files, fmt.Errorf("expected at least %d output files but found %d", minExpected, len(files))
	}

	return files, nil
}
