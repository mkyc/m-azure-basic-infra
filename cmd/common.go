package cmd

import (
	"encoding/json"
	"fmt"
	azbi "github.com/epiphany-platform/e-structures/azbi/v0"
	state "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"
	"io/ioutil"
	"path/filepath"

	terra "github.com/mkyc/go-terraform"
)

type ZeroLogger struct{}

func (z ZeroLogger) Trace(format string, v ...interface{}) {
	logger.
		Trace().
		Msgf(format, v...)
}

func (z ZeroLogger) Debug(format string, v ...interface{}) {
	logger.
		Debug().
		Msgf(format, v...)
}

func (z ZeroLogger) Info(format string, v ...interface{}) {
	logger.
		Info().
		Msgf(format, v...)
}

func (z ZeroLogger) Warn(format string, v ...interface{}) {
	logger.
		Warn().
		Msgf(format, v...)
}

func (z ZeroLogger) Error(format string, v ...interface{}) {
	logger.
		Error().
		Msgf(format, v...)
}

func (z ZeroLogger) Fatal(format string, v ...interface{}) {
	logger.
		Fatal().
		Msgf(format, v...)
}

func (z ZeroLogger) Panic(format string, v ...interface{}) {
	logger.
		Panic().
		Msgf(format, v...)
}

func templateTfVars(c *azbi.Config) error {
	logger.Debug().Msg("templateTfVars")
	tfVarsFile := filepath.Join(ResourcesDirectory, terraformDir, tfVarsFile)
	params := c.Params
	b, err := json.Marshal(&params)
	if err != nil {
		return err
	}
	logger.Info().Msg(string(b))
	err = ioutil.WriteFile(tfVarsFile, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func showModulePlan(c *azbi.Config, s *state.State) error {
	logger.Debug().Msg("showModulePlan")
	futureState := &state.State{}
	err := copier.Copy(futureState, s)
	if err != nil {
		return err
	}
	futureState.AzBI.Config = c
	futureState.AzBI.Status = state.Applied

	//TODO consider adding Output prediction

	diff := cmp.Diff(s, futureState)
	if diff != "" {
		logger.Info().Msg(diff)
		fmt.Println("Planned changes: \n" + diff)
	} else {
		logger.Info().Msg("no changes predicted")
		fmt.Println("No changes predicted.")
	}
	return nil
}

func terraformPlan() string {
	logger.Debug().Msg("terraformPlan")

	options, err := terra.WithDefaultRetryableErrors(&terra.Options{
		TerraformDir: filepath.Join(ResourcesDirectory, terraformDir),
		VarFiles:     []string{filepath.Join(ResourcesDirectory, terraformDir, tfVarsFile)},
		EnvVars: map[string]string{
			"TF_IN_AUTOMATION":    "true",
			"ARM_CLIENT_ID":       clientId,
			"ARM_CLIENT_SECRET":   clientSecret,
			"ARM_SUBSCRIPTION_ID": subscriptionId,
			"ARM_TENANT_ID":       tenantId,
		},
		PlanFilePath:  filepath.Join(SharedDirectory, moduleShortName, applyTfPlanFile),
		StateFilePath: filepath.Join(SharedDirectory, moduleShortName, tfStateFile),
		NoColor:       true,
		Logger:        ZeroLogger{},
	})
	if err != nil {
		logger.Fatal().Err(err)
	}
	s, err := terra.Plan(options)
	if err != nil {
		logger.Fatal().Err(err)
	}
	return s
}

func terraformPlanDestroy() string {
	logger.Debug().Msg("terraformPlanDestroy")

	options, err := terra.WithDefaultRetryableErrors(&terra.Options{
		TerraformDir: filepath.Join(ResourcesDirectory, terraformDir),
		VarFiles:     []string{filepath.Join(ResourcesDirectory, terraformDir, tfVarsFile)},
		EnvVars: map[string]string{
			"TF_IN_AUTOMATION":    "true",
			"ARM_CLIENT_ID":       clientId,
			"ARM_CLIENT_SECRET":   clientSecret,
			"ARM_SUBSCRIPTION_ID": subscriptionId,
			"ARM_TENANT_ID":       tenantId,
		},
		PlanFilePath:  filepath.Join(SharedDirectory, moduleShortName, destroyTfPlanFile),
		StateFilePath: filepath.Join(SharedDirectory, moduleShortName, tfStateFile),
		NoColor:       true,
		Logger:        ZeroLogger{},
	})
	if err != nil {
		logger.Fatal().Err(err)
	}
	s, err := terra.PlanDestroy(options)
	if err != nil {
		logger.Fatal().Err(err)
	}
	return s
}

func terraformApply() (string, error) {
	logger.Debug().Msg("terraformApply")

	options, err := terra.WithDefaultRetryableErrors(&terra.Options{
		TerraformDir: filepath.Join(ResourcesDirectory, terraformDir),
		EnvVars: map[string]string{
			"TF_IN_AUTOMATION":    "true",
			"ARM_CLIENT_ID":       clientId,
			"ARM_CLIENT_SECRET":   clientSecret,
			"ARM_SUBSCRIPTION_ID": subscriptionId,
			"ARM_TENANT_ID":       tenantId,
		},
		PlanFilePath:  filepath.Join(SharedDirectory, moduleShortName, applyTfPlanFile),
		StateFilePath: filepath.Join(SharedDirectory, moduleShortName, tfStateFile),
		NoColor:       true,
		Logger:        ZeroLogger{},
	})
	if err != nil {
		return "", err
	}
	s, err := terra.Apply(options)
	if err != nil {
		return s, err
	}
	return s, nil
}

func getTerraformOutput() (map[string]interface{}, error) {
	logger.Debug().Msg("getTerraformOutput")
	options, err := terra.WithDefaultRetryableErrors(&terra.Options{
		TerraformDir: filepath.Join(ResourcesDirectory, terraformDir),
		EnvVars: map[string]string{
			"TF_IN_AUTOMATION": "true",
		},
		StateFilePath: filepath.Join(SharedDirectory, moduleShortName, tfStateFile),
		Logger:        ZeroLogger{},
	})
	if err != nil {
		return nil, err
	}
	m, err := terra.OutputAll(options)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func terraformDestroy() (string, error) {
	logger.Debug().Msg("terraformDestroy")

	options, err := terra.WithDefaultRetryableErrors(&terra.Options{
		TerraformDir: filepath.Join(ResourcesDirectory, terraformDir),
		EnvVars: map[string]string{
			"TF_IN_AUTOMATION":      "true",
			"TF_WARN_OUTPUT_ERRORS": "1",
			"ARM_CLIENT_ID":         clientId,
			"ARM_CLIENT_SECRET":     clientSecret,
			"ARM_SUBSCRIPTION_ID":   subscriptionId,
			"ARM_TENANT_ID":         tenantId,
		},
		PlanFilePath:  filepath.Join(SharedDirectory, moduleShortName, destroyTfPlanFile),
		StateFilePath: filepath.Join(SharedDirectory, moduleShortName, tfStateFile),
		NoColor:       true,
		Logger:        ZeroLogger{},
	})
	if err != nil {
		return "", err
	}
	s, err := terra.Apply(options)
	if err != nil {
		return s, err
	}
	return s, nil
}

func updateStateAfterDestroy(s *state.State) *state.State {
	logger.Debug().Msg("updateStateAfterDestroy")
	s.AzBI.Output = nil
	s.AzBI.Status = state.Destroyed
	return s
}

func count(output string) (string, error) {
	c, err := terra.Count(output)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Add: %d, Change: %d, Destroy: %d", c.Add, c.Change, c.Destroy), nil
}
