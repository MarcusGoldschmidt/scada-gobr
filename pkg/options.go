package pkg

import (
	"crypto/tls"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/persistence"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/providers"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type ScadagobrOptions struct {
	Address                  string
	Port                     int
	Config                   string
	Logfile                  string
	MaxRecvMsgSize           int
	MetricsServer            bool
	AdminPassword            string
	AdminEmail               string
	DevMode                  bool
	TLSConfig                *tls.Config
	PostgresConnectionString string

	// Jwt
	Expiration        time.Duration
	RefreshExpiration time.Duration
	RefreshKey        []byte
	Key               []byte

	// Shutdown 30 seconds
	ShutdownWait time.Duration

	// must populate location
	Timezone string
	Location *time.Location
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
		AdminEmail:               "admin@localhost",
		PostgresConnectionString: "host=localhost user=postgres password=postgres port=5432",
		RefreshExpiration:        15 * 24 * time.Hour,
		Expiration:               15 * time.Minute,
		ShutdownWait:             30 * time.Second,
		Timezone:                 time.UTC.String(),
	}
}

func ParseOptions() (*ScadagobrOptions, error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	options := DefaultOptions()

	mtls := viper.GetBool("mtls")
	certificate := viper.GetString("certificate")
	pkey := viper.GetString("pkey")
	clientcas := viper.GetString("clientcas")

	options.Address = viper.GetString("address")
	options.Port = viper.GetInt("port")
	options.Logfile = viper.GetString("logfile")
	options.MaxRecvMsgSize = viper.GetInt("max-receive-message-size")
	options.MetricsServer = viper.GetBool("metrics-server")
	options.DevMode = viper.GetBool("dev-mode")
	options.AdminPassword = viper.GetString("admin-password")
	options.RefreshExpiration = viper.GetDuration("refresh-expiration")
	options.Expiration = viper.GetDuration("expiration")
	options.ShutdownWait = viper.GetDuration("shutdown-wait")
	options.PostgresConnectionString = viper.GetString("postgres-connection-string")

	location, err := time.LoadLocation(viper.GetString("timezone"))
	if err != nil {
		return nil, err
	}

	options.Location = location

	tlsConfig, err := setUpTLS(pkey, certificate, clientcas, mtls)
	if err != nil {
		return nil, err
	}

	options.TLSConfig = tlsConfig

	return options, nil
}

func SetupJwtHandler(opt *ScadagobrOptions, persistence persistence.UserPersistence) *auth.JwtHandler {
	handler := &auth.JwtHandler{
		Key:               opt.Key,
		Expiration:        opt.Expiration,
		RefreshExpiration: opt.RefreshExpiration,
		RefreshKey:        opt.RefreshKey,
		TimeProvider:      providers.TimeProviderFromTimeZone(opt.Location),
		UserPersistence:   persistence,
	}

	return handler
}
