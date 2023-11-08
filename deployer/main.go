package main

import (
	"context"
	"fmt"
	"github.com/auth0/go-auth0/management"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type Action struct {
	Id           string `mapstructure:"id"`
	Name         string `mapstructure:"name"`
	CodeFilePath string `mapstructure:"codeFilePath"`
	// Optional
	Dependencies []struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version,omitempty"`
	} `mapstructure:"dependencies,omitempty"`
	Secrets []struct {
		Key    string `mapstructure:"key"`
		Value  string `mapstructure:"value,omitempty"`
		EnvKey string `mapstructure:"envKey,omitempty"`
	} `mapstructure:"secrets,omitempty"`
}

type ActionsConfig struct {
	PostLogin            []Action `mapstructure:"postLogin"`
	CredentialsExchange  []Action `mapstructure:"credentialsExchange"`
	PostChallenge        []Action `mapstructure:"postChallenge"`
	PreUserRegistration  []Action `mapstructure:"preUserRegistration"`
	PostUserRegistration []Action `mapstructure:"postUserRegistration"`
	SendPhoneMessage     []Action `mapstructure:"sendPhoneMessage"`
}

var ActionVersionsMap = map[string]string{
	management.ActionTriggerPostLogin: "v3",
	"credentials-exchange":            "v2",
	"post-challenge":                  "v2",
	"pre-user-registration":           "v2",
	"post-user-registration":          "v2",
	"send-phone-message":              "v2",
}

var ActionsRuntime = "node18-actions"

func main() {
	viper.AutomaticEnv()

	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("viper.ReadInConfig: %w", err))
	}

	auth0API, err := management.New(
		viper.GetString("AUTH0_TENANT_DOMAIN"),
		management.WithClientCredentials(
			context.TODO(),
			viper.GetString("AUTH0_CLIENT_ID"),
			viper.GetString("AUTH0_CLIENT_SECRET")),
	)
	if err != nil {
		log.Fatal(fmt.Errorf("management.New: %w", err))
	}

	var actionsConfig ActionsConfig
	err = viper.UnmarshalKey("actions", &actionsConfig)

	for trigger, actions := range map[string][]Action{
		management.ActionTriggerPostLogin: actionsConfig.PostLogin,
		"credentials-exchange":            actionsConfig.CredentialsExchange,
		"post-challenge":                  actionsConfig.PostChallenge,
		"pre-user-registration":           actionsConfig.PreUserRegistration,
		"post-user-registration":          actionsConfig.PostUserRegistration,
		"send-phone-message":              actionsConfig.SendPhoneMessage,
	} {
		fmt.Printf("Processing trigger: %s\n", trigger)
		for _, action := range actions {
			fmt.Printf("\tUpdating action: %s\n", action.Name)

			var code []byte
			code, err = os.ReadFile(action.CodeFilePath)
			if err != nil {
				log.Fatal(fmt.Errorf("os.ReadFile: %w", err))
			}
			codeStr := string(code)

			dependencies := make([]management.ActionDependency, 0)
			for i := range action.Dependencies {
				dep := action.Dependencies[i]

				dependency := management.ActionDependency{
					Name: &dep.Name,
				}

				version := "latest"
				if dep.Version != "" {
					version = dep.Version
				}
				dependency.Version = &version

				dependencies = append(dependencies, dependency)
			}

			secrets := make([]management.ActionSecret, 0)
			for i := range action.Secrets {
				sec := action.Secrets[i]

				secret := management.ActionSecret{
					Name: &sec.Key,
				}

				var value string
				if sec.Value != "" {
					value = sec.Value
				} else if sec.EnvKey != "" {
					value = viper.GetString(sec.EnvKey)
				}
				secret.Value = &value

				secrets = append(secrets, secret)
			}

			triggerVersion := ActionVersionsMap[trigger]

			auth0Action := &management.Action{
				Name: &action.Name,
				SupportedTriggers: []management.ActionTrigger{{
					ID:      &trigger,
					Version: &triggerVersion,
				}},
				Code:         &codeStr,
				Dependencies: &dependencies,
				Secrets:      &secrets,
				Runtime:      &ActionsRuntime,
			}

			//auth0ActionStr, _ := json.Marshal(auth0Action)
			//fmt.Println(string(auth0ActionStr))

			err = auth0API.Action.Update(context.TODO(), action.Id, auth0Action)
			if err != nil {
				log.Fatal(fmt.Errorf("auth0API.Action.Update: %w", err))
			}

			for i := 0; i < 10; i++ {
				auth0Action, err = auth0API.Action.Read(context.TODO(), action.Id)
				if err != nil {
					log.Fatal(fmt.Errorf("auth0API.Action.Read: %w", err))
				}

				if *auth0Action.Status == management.ActionStatusBuilt {
					break
				}

				time.Sleep(2 * time.Second)
			}
			if *auth0Action.Status != management.ActionStatusBuilt {
				log.Fatal(fmt.Errorf("action status is not built"))
			}

			_, err = auth0API.Action.Deploy(context.TODO(), action.Id)
			if err != nil {
				log.Fatal(fmt.Errorf("auth0API.Action.Deploy: %w", err))
			}
		}
	}
}
