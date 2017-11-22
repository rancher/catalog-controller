package parse

import (
	"github.com/rancher/catalog-controller/utils"
	catalogv1 "github.com/rancher/types/apis/catalog.cattle.io/v1"
	yaml "gopkg.in/yaml.v2"
)

func TemplateInfo(contents []byte) (catalogv1.Template, error) {
	var data map[string]interface{}
	if err := yaml.Unmarshal([]byte(contents), &data); err != nil {
		return catalogv1.Template{}, err
	}

	if _, exists := data["projectURL"]; exists {
		data["project_url"] = data["projectURL"]
	}

	if _, exists := data["version"]; exists {
		data["default_version"] = data["version"]
	} else if _, exists := data["defaultVersion"]; exists {
		data["default_version"] = data["defaultVersion"]
	}

	var template catalogv1.Template
	if err := utils.Convert(data, &template); err != nil {
		return catalogv1.Template{}, err
	}

	return template, nil
}

func CatalogInfoFromTemplateVersion(contents []byte) (catalogv1.Version, error) {
	var template catalogv1.Version
	if err := yaml.Unmarshal(contents, &template); err != nil {
		return catalogv1.Version{}, err
	}

	return template, nil
}

func CatalogInfoFromRancherCompose(contents []byte) (catalogv1.Version, error) {
	cfg, err := utils.CreateConfig(contents)
	if err != nil {
		return catalogv1.Version{}, err
	}
	var rawCatalogConfig interface{}

	if cfg.Version == "2" && cfg.Services[".catalog"] != nil {
		rawCatalogConfig = cfg.Services[".catalog"]
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(contents, &data); err != nil {
		return catalogv1.Version{}, err
	}

	if data["catalog"] != nil {
		rawCatalogConfig = data["catalog"]
	} else if data[".catalog"] != nil {
		rawCatalogConfig = data[".catalog"]
	}

	if rawCatalogConfig != nil {
		var template catalogv1.Version
		if err := utils.Convert(rawCatalogConfig, &template); err != nil {
			return catalogv1.Version{}, err
		}
		return template, nil
	}

	return catalogv1.Version{}, nil
}

func CatalogInfoFromCompose(contents []byte) (catalogv1.Version, error) {
	contents = []byte(extractCatalogBlock(string(contents)))
	return CatalogInfoFromRancherCompose(contents)
}
