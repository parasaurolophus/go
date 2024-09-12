// Copryright 2024 Kirk Rader

package powerview

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Invoke the API exposed by the PowerView hub at the specified address.
func ActivateScene(address string, sceneId int) (response any, err error) {

	url := fmt.Sprintf(`http://%s/scenes?sceneId=%d`, address, sceneId)

	var resp *http.Response
	if resp, err = http.DefaultClient.Get(url); err != nil {
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		return
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	return
}
