package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Colorspace 자료구조는 ocio.config 파일내 colorspaces: 자료구조이다.
type Colorspace struct {
	Name     string `yaml:"name"`
	Family   string `yaml:"family"`
	Bitdepth string `yaml:"bitdepth"`
}

// Displays 자료구조는 ocio.config 파일내 displays: ACES: 자료구조이다.
type Displays struct {
	ACES []View `yaml:"ACES"`
}

// View 자료구조는 ocio.config 파일내 displays: ACES: view 자료구조이다.
type View struct {
	Name       string `yaml:"name"`
	Colorspace string `yaml:"colorspace"`
}

// OCIOConfig 자료구조는 config.ocio 파일 자료구조이다.
type OCIOConfig struct {
	OCIOProfileVersion string       `yaml:"ocio_profile_version"`
	Colorspaces        []Colorspace `yaml:"colorspaces"`
	Displays           `yaml:"displays"`
	Roles              map[string]string `yaml:"roles"`
}

// loadOCIOConfig 함수는 OpenColorIO ocio.config 파일을 분석하여 OCIOConfig 자료구조를 반환한다.
func loadOCIOConfig(configPath string) (OCIOConfig, error) {
	var oc OCIOConfig
	// 파일이 존재하는지 체크한다.
	_, err := os.Stat(configPath)
	if err != nil {
		return oc, err
	}
	// OCIO.config 파일을 불러온다.
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return oc, err
	}
	// 존재하면, 해당파일을 파싱한다.
	err = yaml.Unmarshal(data, &oc)
	if err != nil {
		return oc, err
	}
	return oc, nil
}
