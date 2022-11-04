package command

import (
	"github.com/MarcusGoldschmidt/scadagobr/pkg"
	"github.com/spf13/cobra"
)

func ConfigureSetupFlags(cmd *cobra.Command, opt *pkg.ScadagobrOptions) {
	cmd.Flags().String("postgresConnectionString", opt.PostgresConnectionString, "postgres connection string")
	cmd.Flags().String("address", opt.Address, "address to listen on")
	cmd.Flags().Int("port", opt.Port, "port to listen on")
	cmd.Flags().String("logfile", opt.Logfile, "logfile to write to")
	cmd.Flags().Int("maxRecvMsgSize", opt.MaxRecvMsgSize, "max message size")
	cmd.Flags().Bool("metricsServer", opt.MetricsServer, "enable metrics server")
	cmd.Flags().Bool("dev", opt.DevMode, "dev mode")
	cmd.Flags().String("adminPassword", opt.AdminPassword, "default admin password")
	cmd.Flags().String("adminEmail", opt.AdminEmail, "default admin email")
	cmd.Flags().String("timezone", opt.Timezone, "timezone")

	cmd.Flags().Duration("expiration", opt.Expiration, "expiration time for jwt")
	cmd.Flags().Duration("refreshExpiration", opt.RefreshExpiration, "expiration time for refresh jwt")
	cmd.Flags().Duration("shutdownWait", opt.ShutdownWait, "shutdown wait time")

	//viper.SetDefault("address", "0.0.0.0")
	//viper.SetDefault("port", 11139)
	//viper.SetDefault("logfile", "")
	//// 32Mb
	//viper.SetDefault("maxRecvMsgSize", 1024*1024*32)
	//viper.SetDefault("metricsServer", true)
	//viper.SetDefault("devMode", false)
	//viper.SetDefault("adminPassword", "admin")
	//viper.SetDefault("refreshExpiration", 15*24*time.Hour)
	//viper.SetDefault("expiration", 15*time.Minute)
	//viper.SetDefault("shutdownWait", 30*time.Minute)
	//
	//viper.SetConfigName("scadagobr")
	//
	//viper.ReadInConfig()
	//
	//viper.AutomaticEnv()
}
