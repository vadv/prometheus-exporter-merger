package cmd

import (
	"os"
	"testing"
	"time"
)

func Test_parseEnv(t *testing.T) {
	os.Setenv("LISTEN", ":9090")
	os.Setenv("SCRAPE_TIMEOUT", "120s")
	os.Setenv("URL_8080", "http://127.0.0.1:8080/metrics,keyUrl1_1:valueUrl1_1,keyUrl1_2:valueUrl1_2")
	os.Setenv("URL_8081", "http://127.0.0.1:8081/metrics,keyUrl2_1:valueUrl2_1")
	os.Setenv("URL_8082", "http://127.0.0.1:8082/url3")
	c, err := parseConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if c.Listen != ":9090" {
		t.Fatalf("listen: %s\n", c.Listen)
	}
	if c.ScrapeTimeout != 120*time.Second {
		t.Fatalf("timeout: %s\n", c.ScrapeTimeout)
	}
	for _, s := range c.Sources {
		switch s.Url {
		case "http://127.0.0.1:8080/metrics":
			if s.Labels[`keyUrl1_1`] != `valueUrl1_1` || s.Labels[`keyUrl1_2`] != `valueUrl1_2` {
				t.Fatalf("labels: %v", s.Labels)
			}
		case "http://127.0.0.1:8081/metrics":
			if s.Labels[`keyUrl2_1`] != `valueUrl2_1` {
				t.Fatalf("labels: %v", s.Labels)
			}
		case "http://127.0.0.1:8082/url3":
			if len(s.Labels) > 0 {
				t.Fatalf("labels: %v", s.Labels)
			}
		default:
			t.Fatalf("unknown url: %s", s.Url)
		}
	}
}
