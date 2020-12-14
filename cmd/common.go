package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	azbi "github.com/epiphany-platform/e-structures/azbi/v0"
	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"
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

func templateTfVars(config *azbi.Config) error {
	logger.Debug().Msg("templateTfVars")
	tfVarsFile := filepath.Join(ResourcesDirectory, terraformDir, tfVarsFile)
	params := config.Params
	bytes, err := json.Marshal(&params)
	if err != nil {
		return err
	}
	logger.Info().Msg(string(bytes))
	err = ioutil.WriteFile(tfVarsFile, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func showModulePlan(config *azbi.Config, state *st.State) error {
	logger.Debug().Msg("showModulePlan")
	futureState := &st.State{}
	err := copier.Copy(futureState, state)
	if err != nil {
		return err
	}
	futureState.AzBI.Config = config
	futureState.AzBI.Status = st.Applied

	//TODO consider adding Output prediction

	diff := cmp.Diff(state, futureState)
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
	output, err := terra.Plan(options)
	if err != nil {
		logger.Fatal().Err(err)
	}
	return output
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
	output, err := terra.PlanDestroy(options)
	if err != nil {
		logger.Fatal().Err(err)
	}
	return output
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
	output, err := terra.Apply(options)
	if err != nil {
		return output, err
	}
	return output, nil
}

func getTerraformOutputMap() (map[string]interface{}, error) {
	logger.Debug().Msg("getTerraformOutputMap")
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
	outputMap, err := terra.OutputAll(options)
	if err != nil {
		return nil, err
	}
	return outputMap, nil
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
	output, err := terra.Apply(options)
	if err != nil {
		return output, err
	}
	return output, nil
}

func updateStateAfterDestroy(state *st.State) *st.State {
	logger.Debug().Msg("updateStateAfterDestroy")
	state.AzBI.Output = nil
	state.AzBI.Status = st.Destroyed
	return state
}

func count(output string) (string, error) {
	resourceCount, err := terra.Count(output)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Add: %d, Change: %d, Destroy: %d", resourceCount.Add, resourceCount.Change, resourceCount.Destroy), nil
}
