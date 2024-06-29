package menu

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/PuerkitoBio/goquery"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/gtuk/discordwebhook"
)

func init() {
	functions.CloudEvent("CloudEventFunction", cloudEventFunction)
}

func cloudEventFunction(ctx context.Context, e event.Event) error {
	return Menu()
}

func Menu() error {
	// Request the HTML page.
	res, err := http.Get("https://www.sas.ulisboa.pt/ementas")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	var content = ""

	// Find the Menu
	doc.Find(".menus").Each(func(i int, s1 *goquery.Selection) {
		title := s1.Find(".title_1").Text()

		if strings.Contains(title, "Cantina Velha") {
			s1.Find("p").Each(func(i int, s *goquery.Selection) {
				if strings.Contains(s.Text(), "Cantina Velha") {
					day := strings.ReplaceAll(s.Text(), "Cantina Velha", "")
					content = content + fmt.Sprintf("\n :stew: **%s**\n", day)
				} else {
					if strings.Contains(s.Text(), "Almoço") ||
						strings.Contains(s.Text(), "Jantar") ||
						strings.Contains(s.Text(), "Snack Bar") ||
						strings.Contains(s.Text(), "Linha") ||
						strings.Contains(s.Text(), "Macrobiótica") {
						content = content + fmt.Sprintf("**%s**\n", s.Text())
					} else {
						content = content + s.Text() + "\n"
					}
				}
			})

		}
	})

	if content == "" {
		return fmt.Errorf("no menu found")
	} else {
		content = "Hi! :wave::skin-tone-3: :wave::skin-tone-3: Here is the menu for the week  \n" + content
	}

	var username = "Menu Bot"
	var url = os.Getenv("DISCORD_WEBHOOK_URL")

	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}

	err = discordwebhook.SendMessage(url, message)
	if err != nil {
		return err
	}

	return nil
}
