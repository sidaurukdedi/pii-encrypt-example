package configs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config is an app configuration.
type Config struct {
	Application struct {
		Name           string
		Port           string
		Environment    string
		AllowedOrigins []string
		Timezone       *time.Location
	}
	BasicAuth struct {
		Username string
		Password string
	}
	Crypto struct {
		Secret string
		Pepper string
	}
	Logger struct {
		Formatter logrus.Formatter
	}
	Redis struct {
		Options *redis.Options
	}
	MariadbReadWrite struct {
		Driver             string
		Host               string
		Port               string
		Username           string
		Password           string
		Database           string
		DSN                string
		MaxOpenConnections int
		MaxIdleConnections int
	}
	MariadbReadOnly struct {
		Driver             string
		Host               string
		Port               string
		Username           string
		Password           string
		Database           string
		DSN                string
		MaxOpenConnections int
		MaxIdleConnections int
	}
	Mongodb struct {
		ClientOptions *options.ClientOptions
		Database      string
	}
	SaramaKafka struct {
		Addresses []string
		Config    *sarama.Config
	}
	Captcha struct {
		Host           string
		Secret         string
		Status         string
		AllowedOrigins []string
	}
	GCPStorage struct {
		AccessID   string
		PrivateKey string
	}
	GCPDataStore struct {
		ProjectID   string
		ProjectCred string
	}
	OTPDuration struct {
		LoginSessionDuration time.Duration
		OTPCodeDuration      time.Duration
		TimeToResendOTP      time.Duration
	}
}

// Load will load the configuration.
func Load() *Config {
	cfg := new(Config)
	cfg.app()
	cfg.basicAuth()
	cfg.crypto()
	cfg.logFormatter()
	cfg.redis()
	cfg.mariadbReadOnly()
	cfg.mariadbReadWrite()
	cfg.mongodb()
	cfg.sarama()
	cfg.captcha()
	cfg.gcpStorage()
	cfg.gcpDatastore()
	cfg.otpDuration()
	return cfg
}

func (cfg *Config) app() {
	appName := os.Getenv("APP_NAME")
	appPort := os.Getenv("APP_PORT")
	appEnvironment := os.Getenv("APP_ENVIRONMENT")

	defaultTimezone, _ := time.LoadLocation("Asia/Jakarta")
	appTimezone := os.Getenv("APP_TIMEZONE")
	if timezone, err := time.LoadLocation(appTimezone); err == nil {
		defaultTimezone = timezone
	}

	appAllowedOrigin := os.Getenv("APP_ALLOWED_ORIGINS")
	rawAllowedOrigins := strings.Trim(appAllowedOrigin, " ")
	allowedOrigins := make([]string, 0)
	if rawAllowedOrigins == "" {
		allowedOrigins = append(allowedOrigins, "*")
	} else {
		allowedOrigins = strings.Split(rawAllowedOrigins, ",")
	}

	cfg.Application.Name = appName
	cfg.Application.Port = appPort
	cfg.Application.Environment = appEnvironment
	cfg.Application.AllowedOrigins = allowedOrigins
	cfg.Application.Timezone = defaultTimezone
}

func (cfg *Config) basicAuth() {
	username := os.Getenv("BASIC_AUTH_USERNAME")
	password := os.Getenv("BASIC_AUTH_PASSWORD")

	cfg.BasicAuth.Username = username
	cfg.BasicAuth.Password = password
}

func (cfg *Config) crypto() {
	secret := os.Getenv("AES_SECRET")
	pepper := os.Getenv("AES_PEPPER")

	cfg.Crypto.Pepper = pepper
	cfg.Crypto.Secret = secret
}

func (cfg *Config) logFormatter() {
	formatter := &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			// _, filename := path.Split(f.File)
			filename := fmt.Sprintf("%s:%d", f.File, f.Line)
			return funcname, filename
		},
	}

	cfg.Logger.Formatter = formatter
}

func (cfg *Config) redis() {
	var tlscfg *tls.Config
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	username := os.Getenv("REDIS_USERNAME")
	password := os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.ParseInt(os.Getenv("REDIS_DATABASE"), 10, 64)
	sslEnable, _ := strconv.ParseBool(os.Getenv("REDIS_SSL_ENABLE"))

	if sslEnable {
		tlscfg = &tls.Config{}
		tlscfg.ServerName = host
	}

	options := &redis.Options{
		Addr:      fmt.Sprintf("%s:%s", host, port),
		Username:  username,
		Password:  password,
		DB:        int(db),
		TLSConfig: tlscfg,
	}

	cfg.Redis.Options = options
}

func (cfg *Config) mariadbReadOnly() {
	host := os.Getenv("MARIADB_RO_HOST")
	port := os.Getenv("MARIADB_RO_PORT")
	username := os.Getenv("MARIADB_RO_USERNAME")
	password := os.Getenv("MARIADB_RO_PASSWORD")
	database := os.Getenv("MARIADB_RO_DATABASE")
	maxOpenConnections, _ := strconv.ParseInt(os.Getenv("MARIADB_RO_MAX_OPEN_CONNECTIONS"), 10, 64)
	maxIdleConnections, _ := strconv.ParseInt(os.Getenv("MARIADB_RO_MAX_IDLE_CONNECTIONS"), 10, 64)
	ca := strings.Replace(os.Getenv("MARIADB_RO_CA_ROOT"), `\n`, "\n", -1)

	clientCert := strings.Replace(os.Getenv("MARIADB_RO_CLIENT_CERT"), `\n`, "\n", -1)
	clientKey := strings.Replace(os.Getenv("MARIADB_RO_CLIENT_KEY"), `\n`, "\n", -1)
	tlsName := os.Getenv("MARIADB_RO_TLS_NAME")

	connVal := url.Values{}
	connVal.Add("parseTime", "true")
	connVal.Add("loc", "Asia/Jakarta")

	if ca != "" {
		tlscfg := tls.Config{}
		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM([]byte(ca)); !ok {
			fmt.Println("mariadb failed append pem")
		}

		tlscfg.RootCAs = pool
		tlscfg.InsecureSkipVerify = true

		if clientCert != "" {
			cert, _ := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
			tlscfg.Certificates = []tls.Certificate{cert}
		}

		mysql.RegisterTLSConfig(tlsName, &tlscfg)

		connVal.Add("tls", tlsName)
	}

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	dsn := fmt.Sprintf("%s?%s", dbConnectionString, connVal.Encode())

	cfg.MariadbReadOnly.Driver = "mysql"
	cfg.MariadbReadOnly.Host = host
	cfg.MariadbReadOnly.Port = port
	cfg.MariadbReadOnly.Username = username
	cfg.MariadbReadOnly.Password = password
	cfg.MariadbReadOnly.Database = database
	cfg.MariadbReadOnly.DSN = dsn
	cfg.MariadbReadOnly.MaxOpenConnections = int(maxOpenConnections)
	cfg.MariadbReadOnly.MaxIdleConnections = int(maxIdleConnections)
}

func (cfg *Config) mariadbReadWrite() {
	host := os.Getenv("MARIADB_RW_HOST")
	port := os.Getenv("MARIADB_RW_PORT")
	username := os.Getenv("MARIADB_RW_USERNAME")
	password := os.Getenv("MARIADB_RW_PASSWORD")
	database := os.Getenv("MARIADB_RW_DATABASE")
	maxOpenConnections, _ := strconv.ParseInt(os.Getenv("MARIADB_RW_MAX_OPEN_CONNECTIONS"), 10, 64)
	maxIdleConnections, _ := strconv.ParseInt(os.Getenv("MARIADB_RW_MAX_IDLE_CONNECTIONS"), 10, 64)
	ca := strings.Replace(os.Getenv("MARIADB_RW_CA_ROOT"), `\n`, "\n", -1)

	clientCert := strings.Replace(os.Getenv("MARIADB_RW_CLIENT_CERT"), `\n`, "\n", -1)
	clientKey := strings.Replace(os.Getenv("MARIADB_RW_CLIENT_KEY"), `\n`, "\n", -1)
	tlsName := os.Getenv("MARIADB_RW_TLS_NAME")

	connVal := url.Values{}
	connVal.Add("parseTime", "true")
	connVal.Add("loc", "Asia/Jakarta")

	if ca != "" {
		tlscfg := tls.Config{}
		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM([]byte(ca)); !ok {
			fmt.Println("mariadb failed append pem")
		}

		tlscfg.RootCAs = pool
		tlscfg.InsecureSkipVerify = true

		if clientCert != "" {
			cert, _ := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
			tlscfg.Certificates = []tls.Certificate{cert}
		}

		mysql.RegisterTLSConfig(tlsName, &tlscfg)

		connVal.Add("tls", tlsName)
	}

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	dsn := fmt.Sprintf("%s?%s", dbConnectionString, connVal.Encode())

	cfg.MariadbReadWrite.Driver = "mysql"
	cfg.MariadbReadWrite.Host = host
	cfg.MariadbReadWrite.Port = port
	cfg.MariadbReadWrite.Username = username
	cfg.MariadbReadWrite.Password = password
	cfg.MariadbReadWrite.Database = database
	cfg.MariadbReadWrite.DSN = dsn
	cfg.MariadbReadWrite.MaxOpenConnections = int(maxOpenConnections)
	cfg.MariadbReadWrite.MaxIdleConnections = int(maxIdleConnections)
}

func (cfg *Config) mongodb() {
	appName := os.Getenv("APP_NAME")
	uri := os.Getenv("MONGODB_URL")
	db := os.Getenv("MONGODB_DATABASE")
	minPoolSize, _ := strconv.ParseUint(os.Getenv("MONGODB_MIN_POOL_SIZE"), 10, 64)
	maxPoolSize, _ := strconv.ParseUint(os.Getenv("MONGODB_MAX_POOL_SIZE"), 10, 64)
	maxConnIdleTime, _ := strconv.ParseInt(os.Getenv("MONGODB_MAX_IDLE_CONNECTION_TIME_MS"), 10, 64)

	opts := options.Client().
		ApplyURI(uri).
		SetAppName(appName).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize).
		SetMaxConnIdleTime(time.Millisecond * time.Duration(maxConnIdleTime)).
		SetMonitor(apmmongo.CommandMonitor())

	cfg.Mongodb.ClientOptions = opts
	cfg.Mongodb.Database = db
}

func (cfg *Config) sarama() {
	brokers := os.Getenv("KAFKA_BROKERS")
	sslEnable, _ := strconv.ParseBool(os.Getenv("KAFKA_SSL_ENABLE"))
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	ca := strings.Replace(os.Getenv("KAFKA_CA_ROOT"), `\n`, "\n", -1)
	clientCert := strings.Replace(os.Getenv("KAFKA_CLIENT_CERT"), `\n`, "\n", -1)
	clientKey := strings.Replace(os.Getenv("KAFKA_CLIENT_KEY"), `\n`, "\n", -1)

	sc := sarama.NewConfig()
	sc.Version = sarama.V3_2_0_0
	if username != "" {
		sc.Net.SASL.User = username
		sc.Net.SASL.Password = password
		sc.Net.SASL.Enable = true
	}
	sc.Net.TLS.Enable = sslEnable

	if clientCert != "" {
		tlscfg := tls.Config{}
		cert, _ := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		tlscfg.Certificates = []tls.Certificate{cert}

		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM([]byte(ca)); !ok {
			fmt.Println("kafka fail append pem")
		}

		tlscfg.RootCAs = pool
		sc.Net.TLS.Config = &tlscfg
	}

	// consumer config
	sc.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	sc.Consumer.Offsets.Initial = sarama.OffsetOldest
	// producer config
	sc.Producer.Retry.Backoff = time.Millisecond * 500

	cfg.SaramaKafka.Addresses = strings.Split(brokers, ",")
	cfg.SaramaKafka.Config = sc
}

func (cfg *Config) captcha() {
	host := os.Getenv("GOOGLE_CAPTCHA_HOST")
	secret := os.Getenv("GOOGLE_CAPTCHA_SECRET")
	status := os.Getenv("GOOGLE_CAPTCHA_STATUS")
	rawAllowedOrigins := strings.Trim(os.Getenv("GOOGLE_ALLOWED_ORIGINS"), " ")

	allowedOrigins := make([]string, 0)
	if rawAllowedOrigins == "" {
		allowedOrigins = append(allowedOrigins, "*")
	} else {
		allowedOrigins = strings.Split(rawAllowedOrigins, ",")
	}

	cfg.Captcha.Host = host
	cfg.Captcha.Secret = secret
	cfg.Captcha.Status = status
	cfg.Captcha.AllowedOrigins = allowedOrigins
}

func (cfg *Config) gcpStorage() {
	accessID := os.Getenv("GCP_ACCESS_ID")
	privateKey := os.Getenv("GCP_PRIVATE_KEY")

	cfg.GCPStorage.AccessID = accessID
	cfg.GCPStorage.PrivateKey = string(privateKey)
}

func (cfg *Config) gcpDatastore() {
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	cfg.GCPDataStore.ProjectID = projectID

	projectCred := os.Getenv("DATASTORE_PROJECT_CRED")
	cfg.GCPDataStore.ProjectCred = projectCred
}

func (cfg *Config) otpDuration() {
	defaultLoginSessionDuration := time.Hour * 2
	defaultOtpCodeDuration := time.Hour * 1
	defaultTimeToResendOTP := time.Minute * 3

	cfg.OTPDuration.LoginSessionDuration = defaultLoginSessionDuration
	cfg.OTPDuration.OTPCodeDuration = defaultOtpCodeDuration
	cfg.OTPDuration.TimeToResendOTP = defaultTimeToResendOTP

	loginSessionDuration := os.Getenv("OTP_LOGIN_SESSION_DURATION")
	otpCodeDuration := os.Getenv("OTP_CODE_DURATION")
	timeToResendOTP := os.Getenv("OTP_TIME_TO_RESEND")

	if loginSessionDuration != "" {
		loginSessionDurationInSecond, err := strconv.Atoi(loginSessionDuration)
		if err == nil {
			cfg.OTPDuration.LoginSessionDuration = time.Second * time.Duration(loginSessionDurationInSecond)
		}
	}

	if otpCodeDuration != "" {
		otpCodeDurationInSecond, err := strconv.Atoi(otpCodeDuration)
		if err == nil {
			cfg.OTPDuration.OTPCodeDuration = time.Second * time.Duration(otpCodeDurationInSecond)
		}
	}

	if timeToResendOTP != "" {
		timeToResendOTPInSecond, err := strconv.Atoi(timeToResendOTP)
		if err == nil {
			cfg.OTPDuration.TimeToResendOTP = time.Second * time.Duration(timeToResendOTPInSecond)
		}
	}
}
