package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/wedelivery123/Wedel-ganrate-qrcode/genqr"
	"gopkg.in/ini.v1"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	configFileName = flag.String("config", "local", "Имя конфигурационного файла (без расширения .ini)")
	verbose        = flag.Bool("verbose", false, "Включить режим детального вывода")
)

type contactData struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type configData struct {
	Port     string
	LogoPath string
}

func main() {

	flag.Parse()

	configPath := filepath.Join("utils/config", *configFileName+".ini")

	config, err := parseConfig(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println(config.Port)
	log.Println(config.LogoPath)
	http.HandleFunc("/qr-codes", qrCodesHandler(config.LogoPath))
	addr := fmt.Sprintf(":%s", config.Port)

	logrus.Infof("Server started on %s...", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		logrus.Fatalf("Error starting server: %s", err)
	}

}

func parseConfig(filePath string) (*configData, error) {
	config := &configData{}
	// Load Ini file
	cfg, err := ini.Load(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить файл: %v", err)
	}
	section := cfg.Section("")
	config.Port = section.Key("port").String()
	config.LogoPath = section.Key("pathToLogo").String()

	return config, nil
}

func qrCodesHandler(logo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contact contactData

		err := json.NewDecoder(r.Body).Decode(&contact)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		vCard := fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s %s\nTEL:%s\nEMAIL:%s\nEND:VCARD",
			contact.FirstName, contact.LastName, contact.Phone, contact.Email)

		logoFile, err := os.Open(logo + "Logo.svg")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to open logo file: " + err.Error()))
			return
		}
		defer logoFile.Close()
		pngQr, err := genqr.GenerateQrCode(vCard, logoFile)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "image/png")

		err = png.Encode(w, pngQr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
