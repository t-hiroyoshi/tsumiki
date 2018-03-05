package version

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type PackageInfo struct {
	Tag        string `json:"tag"`
	DockerHost string `json:"docker_host"`
	ImageID    string `json:"image_id"`
}

const PKG_FILE = "packages.json"

func readPackages() ([]PackageInfo, error) {
	raw, err := ioutil.ReadFile(PKG_FILE)
	if err != nil {
		return []PackageInfo{}, err
	}

	var packages []PackageInfo

	err = json.Unmarshal(raw, &packages)
	if err != nil {
		return []PackageInfo{}, nil
	}

	return packages, nil
}

func writePackages(pkgs []PackageInfo) error {
	bytes, err := json.Marshal(pkgs)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(PKG_FILE, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func GetPackage(tag string) (*PackageInfo, error) {
	pkgs, err := readPackages()
	if err != nil {
		return nil, err
	}

	for _, p := range pkgs {
		if p.Tag == tag {
			return &p, nil
		}
	}
	return nil, nil
}

func GetPackages() ([]PackageInfo, error) {
	return readPackages()
}

func AddPackage(pkg PackageInfo) error {
	pkgs, err := GetPackages()
	if err != nil {
		return err
	}

	err = writePackages(append(pkgs, pkg))
	if err != nil {
		return err
	}

	return nil
}

func RemovePackage(pkg PackageInfo, host string) error {
	pkgs, err := readPackages()
	if err != nil {
		return err
	}

	removed := []PackageInfo{}
	for _, p := range pkgs {
		if p.DockerHost != host || p.Tag != pkg.Tag {
			removed = append(removed, p)
		}
	}

	err = writePackages(removed)
	if err != nil {
		return err
	}

	return nil
}
