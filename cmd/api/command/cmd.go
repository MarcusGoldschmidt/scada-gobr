package command

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
)

func Execute(ctx context.Context) error {
	cmd, err := newCommand(ctx)

	if err != nil {
		return err
	}

	err = cmd.Execute()

	return err
}

func newCommand(ctx context.Context) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "scadagobr",
		Short: "ScadaGoBR building your own SCADA",
		Long:  "ScadaGoBR is a system that allows you to build your own SCADA.",
		RunE: func(cmd *cobra.Command, args []string) error {
			opt, err := pkg.ParseOptions()

			if err != nil {
				return err
			}

			scada, err := pkg.DefaultScadagobr(opt)
			if err != nil {
				return err
			}

			err = scada.SetupAndRun(ctx)
			if err != nil {
				return err
			}

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c

			ctx, cancel := context.WithTimeout(ctx, opt.ShutdownWait)
			defer cancel()

			scada.Shutdown(ctx)

			return nil
		},
	}

	cmd.AddCommand(versionCmd())

	ConfigureSetupFlags(cmd, pkg.DefaultOptions())

	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of ScadaGoBR",
		Long:  "All software has versions. This is ScadaGoBR's",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("VERSION: %s\n", pkg.Version)
			cmd.Printf("COMMIT: %s\n", pkg.Commit)

			return nil
		},
	}
}
