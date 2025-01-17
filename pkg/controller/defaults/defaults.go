package defaults

import (
	internalv1 "github.com/acorn-io/acorn/pkg/apis/internal.acorn.io/v1"
	v1 "github.com/acorn-io/acorn/pkg/apis/internal.acorn.io/v1"
	"github.com/acorn-io/acorn/pkg/condition"
	"github.com/acorn-io/acorn/pkg/config"
	"github.com/acorn-io/baaah/pkg/router"
)

func Calculate(req router.Request, resp router.Response) (err error) {
	appInstance := req.Object.(*v1.AppInstance)
	status := condition.Setter(appInstance, resp, v1.AppInstanceConditionDefaults)

	defer func() {
		if err == nil {
			status.Success()
		} else {
			status.Error(err)
			// Disable prune because we are short-circuiting the handlers and don't want objects getting cleaned up accidentally.
			resp.DisablePrune()
			// Don't return the error so the other conditions don't get the same information.
			err = nil
		}
	}()
	if appInstance.Generation != appInstance.Status.ObservedGeneration {
		if err = calculate(req, appInstance); err != nil {
			return err
		}
	}

	return nil
}

func calculate(req router.Request, appInstance *internalv1.AppInstance) error {
	cfg, err := config.Get(req.Ctx, req.Client)
	if err != nil {
		return err
	}

	if err := addVolumeClassDefaults(req.Ctx, req.Client, appInstance); err != nil {
		return err
	}

	if err := addDefaultMemory(req, cfg, appInstance); err != nil {
		return err
	}

	return nil
}
