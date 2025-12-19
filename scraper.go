package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Kullanım: go run scraper.go <URL>")
		return
	}

	url := os.Args[1]

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP isteği başarısız:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("HTTP hata kodu:", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("HTML okunamadı:", err)
		return
	}

	err = os.WriteFile("site_data.html", body, 0644)
	if err != nil {
		fmt.Println("Dosyaya yazılamadı:", err)
		return
	}

	fmt.Println("HTML kaydedildi: site_data.html")

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var screenshot []byte

	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
		chromedp.FullScreenshot(&screenshot, 100),
	)

	if err != nil {
		fmt.Println("Screenshot alınamadı:", err)
		return
	}

	err = os.WriteFile("screenshot.png", screenshot, 0644)
	if err != nil {
		fmt.Println("Screenshot dosyaya yazılamadı:", err)
		return
	}

	fmt.Println("FULL PAGE screenshot kaydedildi: screenshot.png")
}
