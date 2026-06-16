//go:build ignore

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	data, err := os.ReadFile("Bank of Singapore - Invitation To Connect .html")
	if err != nil {
		fmt.Println("Error reading HTML file:", err)
		os.Exit(1)
	}

	html := string(data)

	// Match all base64 image data in src attributes (both <img> and <v:imagedata>)
	re := regexp.MustCompile(`(?:src)="data:image/(png|jpeg|jpg);base64,([^"]+)"`)
	matches := re.FindAllStringSubmatch(html, -1)

	if len(matches) == 0 {
		fmt.Println("No base64 images found.")
		return
	}

	names := []string{"header_logo", "footer_banner_vml", "footer_banner"}
	for i, match := range matches {
		ext := match[1]
		b64 := strings.ReplaceAll(match[2], "\n", "")
		b64 = strings.ReplaceAll(b64, "\r", "")
		b64 = strings.ReplaceAll(b64, " ", "")

		imgData, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			fmt.Printf("[%d] Error decoding base64: %v\n", i, err)
			continue
		}

		name := fmt.Sprintf("image_%d.%s", i, ext)
		if i < len(names) {
			name = fmt.Sprintf("%s.%s", names[i], ext)
		}

		err = os.WriteFile(name, imgData, 0644)
		if err != nil {
			fmt.Printf("[%d] Error writing file: %v\n", i, err)
			continue
		}

		fmt.Printf("[%d] Saved: %s (%d bytes)\n", i, name, len(imgData))
	}
}
