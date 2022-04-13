package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var ConfigValueRegexp = regexp.MustCompile("(\\$\\{([a-zA-Z][a-zA-Z0-9_]*)\\})")

func ResolveConfigValue(v string) string {

	matches := ConfigValueRegexp.FindAllSubmatch([]byte(v), -1)

	for _, m := range matches {
		env, ok := os.LookupEnv(string(m[2]))
		if ok {
			v = strings.ReplaceAll(v, string(m[1]), env)
		}
	}

	return v
}

func ReadFileAndResolveEnvVars(cfgFile string) ([]byte, error) {

	fsz := FileSize(cfgFile)
	if fsz < 0 {
		return nil, fmt.Errorf("error reading file %s", cfgFile)
	}

	file, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sb bytes.Buffer
	sb.Grow(int(fsz))

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sb.WriteString(ResolveConfigValue(scanner.Text()))
		sb.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sb.Bytes(), nil
}

/*
func ReadFileAndResolveEnvVars(cfgFile string) (string, error) {

	fsz := FileSize(cfgFile)
	if fsz < 0 {
		return "", fmt.Errorf("error reading file %s", cfgFile)
	}

	file, err := os.Open(cfgFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var sb strings.Builder
	sb.Grow(int(fsz))

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sb.WriteString(ResolveConfigValue(scanner.Text()))
		sb.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}


func ReadConfig(fileConfigPathEnvVar string, defaultConfigContent string, resolveEnv bool) (string, string, error) {

	configPath := os.Getenv(fileConfigPathEnvVar)
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			cfgContent, rerr := ReadFileAndResolveEnvVars(configPath)
			if rerr != nil {
				return configPath, "", err
			}
			return configPath, cfgContent, nil
		}
		return configPath, "", fmt.Errorf("the %s env variable has been set but no file cannot be found at %s", fileConfigPathEnvVar, configPath)
	}

	if len(defaultConfigContent) > 0 {
		return "", ResolveConfigValue(string(defaultConfigContent)), nil
	}

	return "", "", fmt.Errorf("the config path variable %s has not been set; please set", fileConfigPathEnvVar)
}
*/

func ReadConfig(fileConfigPathEnvVar string, defaultConfigContent []byte, resolveEnv bool) (string, []byte, error) {

	configPath := os.Getenv(fileConfigPathEnvVar)
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			cfgContent, rerr := ReadFileAndResolveEnvVars(configPath)
			if rerr != nil {
				return configPath, nil, err
			}
			return configPath, cfgContent, nil
		}
		return configPath, nil, fmt.Errorf("the %s env variable has been set but no file cannot be found at %s", fileConfigPathEnvVar, configPath)
	}

	if len(defaultConfigContent) > 0 {
		return "", []byte(ResolveConfigValue(string(defaultConfigContent))), nil
	}

	return "", nil, fmt.Errorf("the config path variable %s has not been set; please set", fileConfigPathEnvVar)
}
