package commonConfig

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"isp-ctl/bash"
	"isp-ctl/cfg"
	"isp-ctl/command/utils"
	"isp-ctl/flag"
	"isp-ctl/service"
	"os"
)

func Set() cli.Command {
	return cli.Command{
		Name:         "set",
		Usage:        "set common configurations",
		Action:       set.action,
		BashComplete: bash.CommonConfig.GetSetDelete,
	}
}

var set setCommand

type setCommand struct{}

func (g setCommand) action(ctx *cli.Context) {
	if err := flag.CheckGlobal(ctx); err != nil {
		utils.PrintError(err)
		return
	}

	ccName := ctx.Args().First()
	pathObject := ctx.Args().Get(1)
	changeObject := ctx.Args().Get(2)

	if ccName == "" {
		utils.PrintError(errors.New("empty config name"))
		return
	}

	mapConfigByName, err := service.Config.GetMapCommonConfigByName()
	if err != nil {
		utils.PrintError(err)
		return
	}

	config, ok := mapConfigByName[ccName]
	if !ok {
		config = cfg.CommonConfig{
			Name: ccName,
		}
	}

	pathObject, err = utils.CheckPath(pathObject)
	if err != nil {
		utils.PrintError(err)
		return
	}

	if changeObject == "" {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			utils.PrintError(err)
			return
		}
		changeObject = string(bytes)
	}
	if changeObject == "" {
		utils.PrintError(errors.New("expected argument"))
		return
	}

	if pathObject == "" {
		utils.CreateUpdateCommonConfig(changeObject, config)
		return
	} else {
		jsonObject, err := json.Marshal(config.Data)
		if err != nil {
			utils.PrintError(err)
			return
		}

		changeArgument := utils.ParseSetObject(changeObject)
		if stringToChange, err := sjson.Set(string(jsonObject), pathObject, changeArgument); err != nil {
			utils.PrintError(err)
			return
		} else {
			utils.CreateUpdateCommonConfig(stringToChange, config)
			return
		}
	}
}