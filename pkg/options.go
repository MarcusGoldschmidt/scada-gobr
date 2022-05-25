package pkg

import (
	"crypto/tls"
	"github.com/spf13/viper"
)

type ScadagobrOptions struct {
	Address                  string
	Port                     int
	Config                   string
	Logfile                  string
	MaxRecvMsgSize           int
	MetricsServer            bool
	AdminPassword            string
	DevMode                  bool
	TLSConfig                *tls.Config
	PostgresConnectionString string
}

func DefaultOptions() *ScadagobrOptions {
	return &ScadagobrOptions{
		Address:                  "0.0.0.0",
		Port:                     11139,
		Logfile:                  "",
		TLSConfig:                nil,
		MaxRecvMsgSize:           1024 * 1024 * 32, // 32Mb
		MetricsServer:            true,
		DevMode:                  false,
		AdminPassword:            "admin",
		PostgresConnectionString: "scadagobr",
	}
}

func ConfigureFlags(configFile ...string) error {

	viper.SetDefault("postgresConnectionString", "host=localhost user=postgres password=postgres port=5432")

	viper.SetConfigName("scadagobr")

	if len(configFile) != 0 {
		file := configFile[0]
		viper.SetConfigFile(file)
	}

	viper.AutomaticEnv()

	return nil
}

func ParseOptions() (*ScadagobrOptions, error) {

	options := DefaultOptions()

	mtls := viper.GetBool("mtls")
	certificate := viper.GetString("certificate")
	pkey := viper.GetString("pkey")
	clientcas := viper.GetString("clientcas")

	postgresConnectionString := viper.GetString("postgresConnectionString")

	tlsConfig, err := setUpTLS(pkey, certificate, clientcas, mtls)
	if err != nil {
		return nil, err
	}

	options.TLSConfig = tlsConfig

	options.PostgresConnectionString = postgresConnectionString

	return options, nil
}
