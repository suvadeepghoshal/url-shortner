package home

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"url-shortner/controllers/util"
	TYPE "url-shortner/model/type"
	"url-shortner/view/home"
)

func HomeController(writer http.ResponseWriter, request *http.Request) error {
	slog.Info("inside HomeController")
	resp, err := http.Get("http://localhost:1323/init")

	if err != nil {
		http.Error(writer, "Failed to call Init API", http.StatusInternalServerError)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("inside HomeController :: not able to clean up network resources: ", err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		http.Error(writer, "inside HomeController :: init API returned error : ", http.StatusInternalServerError)
		return err
	}

	var initResponse TYPE.CommonResponse
	if err := json.NewDecoder(resp.Body).Decode(&initResponse); err != nil {
		http.Error(writer, "Failed to decode init API response", http.StatusInternalServerError)
		return err
	}

	slog.Info("Init API is decoded")

	writer.WriteHeader(http.StatusOK)
	return util.Render(writer, request, home.Index(initResponse))
}
