package redirect

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

func RedirectController(writer http.ResponseWriter, request *http.Request) {
	slog.Info("inside RedirectController")

	hash := chi.URLParam(request, "hash")

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(fmt.Sprintf("The hash ID received is: %s", hash))
	if err != nil {
		slog.Error("Unable to write response: ", "err", err.Error())
		http.Error(writer, "Unable to write response", http.StatusInternalServerError)
		return
	}
}
