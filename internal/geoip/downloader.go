package geoip

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath" // Необходимо для работы с путями
	"time"

	"github.com/sirupsen/logrus"
)

const (
	ASNUrl      = "https://github.com/P3TERX/GeoLite.mmdb/releases/latest/download/GeoLite2-ASN.mmdb"
	CityUrl     = "https://github.com/P3TERX/GeoLite.mmdb/releases/latest/download/GeoLite2-City.mmdb"
	CountryUrl  = "https://github.com/P3TERX/GeoLite.mmdb/releases/latest/download/GeoLite2-Country.mmdb"
	ASNPath     = "geo/GeoLite2-ASN.mmdb"
	CityPath    = "geo/GeoLite2-City.mmdb"
	CountryPath = "geo/GeoLite2-Country.mmdb"
)

// DownloadDB downloads the latest GeoLite2 databases with retries.
func DownloadDB() error {
	if err := downloadWithRetry("ASN", ASNUrl, ASNPath); err != nil {
		return err
	}
	if err := downloadWithRetry("Country", CountryUrl, CountryPath); err != nil {
		return err
	}
	return downloadWithRetry("City", CityUrl, CityPath)
}

func downloadWithRetry(name, url, path string) error {
	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		logrus.Infof("Downloading GeoLite2-%s database (attempt %d/%d)...", name, i+1, maxRetries)
		err := downloadFile(path, url)
		if err == nil {
			logrus.Infof("GeoLite2-%s database downloaded successfully", name)
			return nil
		}
		lastErr = err
		logrus.WithError(err).Warnf("Download attempt %d for %s failed", i+1, name)
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("failed to download GeoLite2-%s database after %d attempts: %w", name, maxRetries, lastErr)
}

func downloadFile(path string, url string) error {
	// 1. Извлекаем директорию из полного пути и создаем её, если её нет
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// 2. Теперь безопасно создаем файл
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// 3. Получаем данные по сети
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// 4. Записываем тело ответа в файл
	_, err = io.Copy(out, resp.Body)
	return err
}
