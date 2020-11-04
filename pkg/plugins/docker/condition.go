package docker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/olblak/updateCli/pkg/core/scm"
)

// ConditionFromSCM returns an error because it's not supported
func (d *Docker) ConditionFromSCM(source string, scm scm.Scm) (bool, error) {
	return false, fmt.Errorf("SCM configuration is not supported for dockerRegistry condition, aborting")
}

// Condition checks if a docker image with a specific tag is published
func (d *Docker) Condition(source string) (bool, error) {
	URL := ""

	if d.Tag != "" {
		fmt.Printf("Tag %v, already defined from configuration file\n", d.Tag)
	} else {
		d.Tag = source
	}

	if ok, err := d.Check(); !ok {
		return false, err
	}

	if d.isDockerHub() {
		URL = fmt.Sprintf("https://%s/v2/repositories/%s/tags/%s/",
			d.URL,
			d.Image,
			d.Tag)

	} else {
		if ok, err := d.IsDockerRegistry(); !ok {
			return false, err
		}
		URL = fmt.Sprintf("https://%s/v2/%s/manifests/%s",
			d.URL,
			d.Image,
			d.Tag)
	}

	req, err := http.NewRequest("GET", URL, nil)

	if err != nil {
		return false, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return false, err
	}

	if res.StatusCode == 200 && !d.isDockerHub() {
		fmt.Printf("\u2714 %s:%s available on the Docker Registry\n", d.Image, d.Tag)
		return true, nil

	} else if d.isDockerHub() {

		data := map[string]string{}

		json.Unmarshal(body, &data)

		if val, ok := data["message"]; ok && strings.Contains(val, "not found") {
			fmt.Printf("\u2717 %s:%s doesn't exist on the Docker Registry \n", d.Image, d.Tag)
			return false, nil
		}

		if val, ok := data["name"]; ok && val == d.Tag {
			fmt.Printf("\u2714 %s:%s available on the Docker Registry\n", d.Image, d.Tag)
			return true, nil
		}

	} else {

		fmt.Printf("\u2717Something went wrong on URL: %s\n", URL)
	}

	return false, fmt.Errorf("something went wrong %s", URL)
}