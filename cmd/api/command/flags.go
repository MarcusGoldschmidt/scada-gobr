package command

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg"
	"github.com/spf13/cobra"
)

func ConfigureSetupFlags(cmd *cobra.Command, opt *pkg.ScadagobrOptions) {
	cmd.Flags().String("postgres-connection-string", opt.PostgresConnectionString, "postgres connection string")
	cmd.Flags().String("address", opt.Address, "address to listen on")
	cmd.Flags().Int("port", opt.Port, "port to listen on")
	cmd.Flags().String("logfile", opt.Logfile, "logfile to write to")
	cmd.Flags().Int("max-recv-msg-size", opt.MaxRecvMsgSize, "max message size")
	cmd.Flags().Bool("metrics-server", opt.MetricsServer, "enable metrics server")
	cmd.Flags().Bool("dev", opt.DevMode, "dev mode")
	cmd.Flags().String("admin-password", opt.AdminPassword, "default admin password")
	cmd.Flags().String("admin-email", opt.AdminEmail, "default admin email")
	cmd.Flags().String("timezone", opt.Timezone, "timezone")

	cmd.Flags().Duration("expiration", opt.Expiration, "expiration time for jwt")
	cmd.Flags().Duration("refresh-expiration", opt.RefreshExpiration, "expiration time for refresh jwt")
	cmd.Flags().Duration("shutdown-wait", opt.ShutdownWait, "shutdown wait time")
}
