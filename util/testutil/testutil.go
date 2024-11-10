package testutil

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/fileutil"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type StorageAccountKeys struct {
	AccountName string `yaml:"account-name"`
	AccountKey  string `yaml:"account-key"`
}

func ReadConfig(fileName string) (StorageAccountKeys, error) {

	stg := StorageAccountKeys{}

	wd, _ := os.Getwd()
	log.Info().Str("wd", wd).Msg("working dir")

	var b []byte

	configPath := fileutil.FindFileInClosestDirectory(".", fileName)
	if configPath == "" {
		return stg, fmt.Errorf("cannot find config file of name %s", fileName)
	}

	log.Info().Str("file-name", configPath).Msg("found config file")

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return stg, err
	}

	err = yaml.Unmarshal(b, &stg)
	if err != nil {
		return stg, err
	}

	if stg.AccountName == "" || stg.AccountKey == "" {
		return stg, fmt.Errorf("config file %s does not contain storage info", configPath)
	}

	return stg, nil
}
