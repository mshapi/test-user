package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	test_user "test-user"
	"test-user/pkg/repository"
	"test-user/pkg/service"
)

type Config struct {
	Server struct {
		Host string `ini:"host"`
		Port uint   `ini:"port"`
	} `ini:"server"`
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

func writeServerError(w http.ResponseWriter, status int, err error) {
	if status > 0 {
		w.WriteHeader(status)
	}
	var errStr string
	if err != nil {
		errStr = err.Error()
	} else {
		errStr = "unknown error"
	}
	fmt.Fprintf(w, "{\"error\": \"%s\"}", errStr)
}

func main() {
	var confPath string

	flag.StringVar(&confPath, "conf", "conf.ini", "Configuration file path")
	flag.Parse()

	config, err := getConfig(confPath)
	if err != nil {
		log.Fatalf("Ошибка чтения конфига '%s': %s\n", confPath, err)
	}

	repo := repository.NewMemoryRepository()
	userService := service.NewUserService(repo)

	handler := http.NewServeMux()

	userHandler := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json")

		// костыль
		url := r.URL.Path
		if url == "/user" {
			url = "/user/"
		}
		// ::end

		id := url[len("/user/"):]

		if r.Method != http.MethodGet && (id != "") == (r.Method == http.MethodPost) {
			writeServerError(w, http.StatusBadRequest, errors.New("bad request"))
			return
		}

		var data interface{}
		var err error

		switch {
		case r.Method == http.MethodGet && id == "":
			data, _ = userService.GetUsers()
		case r.Method == http.MethodGet && id != "":
			data, err = userService.GetUser(id)
		case r.Method == http.MethodDelete:
			data, err = userService.DeleteUser(id)
		case r.Method == http.MethodPost || r.Method == http.MethodPut:
			var user test_user.User
			decoder := json.NewDecoder(r.Body)
			err = decoder.Decode(&user)
			if err != nil {
				break
			}

			if r.Method == http.MethodPost {
				user, err = userService.CreateUser(user)
			} else {
				user.ID = id
				user, err = userService.UpdateUser(user)
			}
			if err == nil {
				data = user
			}

		default:
			writeServerError(w, http.StatusBadRequest, errors.New("bad request"))
			return
		}

		if err == nil {
			encoder := json.NewEncoder(w)
			err = encoder.Encode(data)
		} else if _, ok := err.(*repository.UserNotFoundError); ok {
			writeServerError(w, http.StatusNotFound, err)
			return
		}

		if err != nil {
			log.Printf("http request error: %s", err)
			writeServerError(w, http.StatusInternalServerError, err)
		}
	}

	handler.HandleFunc("/user", userHandler)
	handler.HandleFunc("/user/", userHandler)

	go func() {
		err = http.ListenAndServe(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port), handler)
		if err != nil {
			log.Fatalf("Ошибка запуска сервера: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
