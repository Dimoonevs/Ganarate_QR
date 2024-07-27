package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/wedelivery123/Wedel-ganrate-qrcode/genqr"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

func main() {

	flag.Parse()

	configPath := filepath.Join("utils/config", *configFileName+".ini")

	port, err := parseConfig(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println(port)
	http.HandleFunc("/qr-codes", qrCodesHandler)
	addr := fmt.Sprintf(":%s", port)

	logrus.Infof("Server started on %s...", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		logrus.Fatalf("Error starting server: %s", err)
	}

}

func parseConfig(filePath string) (string, error) {
	// Load Ini file
	cfg, err := ini.Load(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось загрузить файл: %v", err)
	}
	section := cfg.Section("")
	port := section.Key("port").String()

	return port, nil
}

func qrCodesHandler(w http.ResponseWriter, r *http.Request) {
	var contact contactData

	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	vCard := fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nFN:%s %s\nTEL:%s\nEMAIL:%s\nEND:VCARD",
		contact.FirstName, contact.LastName, contact.Phone, contact.Email)
	png, err := genqr.GenerateQrCode(vCard)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(png)))

	_, err = w.Write(png)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
