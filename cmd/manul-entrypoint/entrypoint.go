package main

import (
	"github.com/argcv/manul/model"
	"github.com/argcv/manul/version"
	"github.com/argcv/webeh/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func main() {
	viper.SetConfigName("manul.project")
	viper.AddConfigPath(".")
	if conf := os.Getenv("MANUL_PROJECT_CFG"); conf != "" {
		viper.SetConfigFile(conf)
	}

	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.Infof("Manul Loader Started. Version: %s (%s) Built At: %v", version.Version, version.GitHash, version.BuildDate)
			if verbose, err := cmd.Flags().GetBool("verbose"); err == nil {
				if verbose {
					log.Verbose()
					log.Debug("verbose mode: ON")
				}
			}
			if conf, _ := cmd.Flags().GetString("config"); conf != "" {
				viper.SetConfigFile(conf)
			}
			err := viper.ReadInConfig()
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok && err != nil {
				return err
			}
			if conf := viper.ConfigFileUsed(); conf != "" {
				log.Debugf("using config file: %s", conf)
			} else {
				return errors.New("configure file not found!!")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, e := model.LoadProjectConfig("manul.project.yml")

			if e != nil {
				return e
			}

			bse := cfg.ToBashScriptsExecutor()

			bse.Id = "j1"

			log.Infof("Scripts: [%s]", bse.EncodedScript())
			out, err := bse.Perform()
			log.Infof(".... out[\n%s\n] err::[%v]", string(out), err)

			//envs := viper.GetStringSlice("envs")
			//
			//log.Infof("[[[[[%v]]]]]", envs)
			//
			//outputDir, _ := cmd.Flags().GetString("out")
			//
			//log.Infof("output folder: %v", outputDir)
			//
			//os.MkdirAll(outputDir, 0700)
			//
			//bse := helpers.NewBashScriptsExecutor("123")
			//
			//bse.SetEnv(envs...)
			//for _, stage := range stages {
			//	log.Info("stage:", stage)
			//	script := viper.GetString(fmt.Sprintf("jobs.%s.script", stage))
			//	scripts := viper.GetStringSlice(fmt.Sprintf("jobs.%s.scripts", stage))
			//
			//	bse.AddEnv(envs...)
			//	if len(script) > 0 {
			//		scripts = append(scripts, script)
			//	}
			//	bse.AddScriptsInStage(stage, scripts...)
			//}
			//log.Infof("Scripts: [%s]", bse.EncodedScript())
			//out, err := bse.Perform()
			//log.Infof(".... out[\n%s\n] err::[%v]",  string(out), err)

			return nil
		},
	}
	cmd.PersistentFlags().StringP("config", "c", "", "explicit assign a configuration file")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "log verbose")
	cmd.PersistentFlags().StringP("out", "o", "/tmp", "output folder")
	log.Infof("start.[%v]", os.Args[0])
	if err := cmd.Execute(); err != nil {
		log.Infof("%v: %v", os.Args[0], err)
		os.Exit(1)
	}
}
