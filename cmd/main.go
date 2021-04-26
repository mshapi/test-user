package main

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	Server struct {
		Host string `ini:"host"`
		Port uint   `ini:"port"`
	} `ini:"server"`
	Storage struct {
		UseFile bool   `ini:"use_file"`
		Path    string `ini:"path"` // json storage file path
	} `ini:"storage"`
}

func (c *Config) load(file string) error {
	cfg, err := ini.Load(file)
	if cfg != nil {
		err = cfg.MapTo(c)
	}
	return err
}

func getConfig(path string) (*Config, error) {
	conf := new(Config)
	err := conf.load(path)
	if err != nil && conf.Server.Port == 0 {
		conf.Server.Port = 8000 // set default port
	}
	return conf, err
}

func main() {
	var confPath string

	flag.StringVar(&confPath, "conf", "conf.ini", "Configuration file path")
	flag.Parse()

	config, err := getConfig(confPath)
	if err != nil {
		log.Fatalf("Ошибка конфига '%s': %s\n", confPath, err)
	}

	handler := http.NewServeMux()

	go func() {
		err = http.ListenAndServe(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port), handler)
		if err != nil {
			log.Fatalf("Ошибка запуска сервера: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	fmt.Printf("Config path: %s\n", confPath)
	fmt.Printf("Config: %#v\n", config)
}
