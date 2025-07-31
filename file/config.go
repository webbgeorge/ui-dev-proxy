package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/webbgeorge/ui-dev-proxy/domain"
)

func ConfigProvider() domain.ConfigProvider {
	return func(path string) (domain.Config, error) {
		f, err := os.Open(path) // #nosec G304 -- false positive, this is the config file passed to the app by the user
		if err != nil {
			return domain.Config{}, err
		}
		var c domain.Config
		err = json.NewDecoder(f).Decode(&c)
		if err != nil {
			return domain.Config{}, err
		}

		// files referenced by config must be at or below the level of the config file
		fsRoot, err := os.OpenRoot(filepath.Dir(f.Name()))
		if err != nil {
			return domain.Config{}, err
		}

		for _, r := range c.Routes {
			if r.Type != domain.RouteTypeMock {
				if r.Redirect != nil {
					redirectType := r.Redirect.Type
					if redirectType != "permanent" && redirectType != "temporary" {
						return domain.Config{}, fmt.Errorf("invalid redirect type '%s'", redirectType)
					}
				}

				continue
			}

			if r.Mock == nil {
				return domain.Config{}, errors.New("missing mock config on mock type route")
			}

			r.Mock.MatchRequest.Body, err = getBody(r.Mock.MatchRequest.Body, fsRoot)
			if err != nil {
				return domain.Config{}, err
			}

			r.Mock.Response.Body, err = getBody(r.Mock.Response.Body, fsRoot)
			if err != nil {
				return domain.Config{}, err
			}
		}

		return c, nil
	}
}

func getBody(body string, fsRoot *os.Root) (string, error) {
	if !strings.HasSuffix(body, ".json") {
		return body, nil
	}

	f, err := fsRoot.Open(body)
	if err != nil {
		return "", err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
