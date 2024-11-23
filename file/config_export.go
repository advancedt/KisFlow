package file

import (
	"KisFlow/common"
	"KisFlow/kis"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// ConfigExportYaml 将flow配置输出，且存储本地
func ConfigExportYaml(flow kis.Flow, savePath string) error {
	if data, err := yaml.Marshal(flow.GetConfig()); err != nil {
		return nil
	} else {
		// flow
		flowFilePath := filepath.Join(savePath, common.KisIdTypeFlow+"-"+flow.GetName()+".yaml")
		err := os.WriteFile(flowFilePath, data, 0644)
		if err != nil {
			return err
		}

		// function
		for _, fp := range flow.GetConfig().Flows {
			fConf := flow.GetFuncConfigByName(fp.FuncName)
			if fConf == nil {
				return errors.New(fmt.Sprintf("function name = %s config is nil ", fp.FuncName))
			}

			if fdata, err := yaml.Marshal(fConf); err != nil {
				return err
			} else {
				funcFilePath := filepath.Join(savePath, common.KisIdTypeFlow+"-"+fp.FuncName+".yaml")
				if err := os.WriteFile(funcFilePath, fdata, 0644); err != nil {
					return err
				}
			}
			// Connector
			if fConf.Option.CName != "" {
				cConf, err := fConf.GetConnConfig()
				if err != nil {
					return err
				}
				if cdata, err := yaml.Marshal(cConf); err != nil {
					return err
				} else {
					connFilePath := filepath.Join(savePath, common.KisIdTypeConnnector+"-"+cConf.CName+".yaml")
					if err := os.WriteFile(connFilePath, cdata, 0644); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
