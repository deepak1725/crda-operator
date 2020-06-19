package controller

import (
	"operator/crda-operator/pkg/controller/codereadyanalytics"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, codereadyanalytics.Add)
}
