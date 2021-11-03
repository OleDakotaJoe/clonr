package core

import (
	"github.com/oledakotajoe/clonr/utils"
	"github.com/robertkrimen/otto"
	"github.com/spf13/cast"
)

func RunScriptAndReturnValue(script string) string {
	vm := otto.New()
	val, err := vm.Run(script)
	utils.ExitIfError(err)
	return cast.ToString(val)
}
