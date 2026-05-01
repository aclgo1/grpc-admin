package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	ApiPort           string `mapstructure:"API_PORT"`
	DbDriver          string `mapstructure:"DB_DRIVER"`
	DbURI             string `mapstructure:"DB_URI"`
	DevMode           bool   `mapstructure:"DEV_MODE"`
	Encoding          string `mapstructure:"ENCODING"`
	LogLevel          string `mapstructure:"LOG_LEVEL"`
	TracerExporterURL string `mapstructure:"TRACER_EXPORTER_URL"`
	MetricExporterURL string `mapstructure:"METER_EXPORTER_URL"`
	RedisAddr         string `mapstructure:"REDIS_ADDR"`
	RedisDB           int    `mapstructure:"REDIS_DB"`
	RedisPass         string `mapstructure:"REDIS_PASS"`
}

func NewConfig(path string) *Config {
	var cfg Config

	_, ok := os.LookupEnv("PROD")
	if !ok {

		log.Printf("env variable PROD (NOT) set, reading from file")

		err := cfg.LoadFromFile(path)
		if err != nil {
			log.Fatalf("NewConfig: %v", err)
		}
		return &cfg
	}

	log.Printf("env variable PROD set, reading from system")

	err := cfg.LoadFromEnv()
	if err != nil {
		log.Fatalf("LoadFromEnv: %v", err)
	}

	return &cfg

}

func (c *Config) LoadFromEnv() error {
	elems := reflect.ValueOf(c).Elem()
	tp := elems.Type()

	errs := make([]error, elems.NumField())

	for i := 0; i < elems.NumField(); i++ {
		tag := tp.Field(i).Tag.Get("mapstrucuture")
		env := os.Getenv(tag)

		switch elems.Field(i).Type().Name() {
		case "string":
			elems.Field(i).SetString(env)
		case "int":
			parse, err := strconv.ParseInt(env, 10, 64)
			if err == nil {
				elems.Field(i).SetInt(parse)
			}

			errs[i] = err

		case "bool":
			parse, err := strconv.ParseBool(env)
			if err == nil {
				elems.Field(i).SetBool(parse)
			}

			errs[i] = err
		}
	}

	return errors.Join(errs...)
}

func (c *Config) LoadFromFile(path string) error {

	viper.SetConfigName("app")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("LoadFile.ReadInConfig: %v", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		return fmt.Errorf("LoadFile.Unmarshal: %v", err)
	}

	return nil
}
