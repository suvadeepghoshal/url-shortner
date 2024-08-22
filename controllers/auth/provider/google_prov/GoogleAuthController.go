package google_prov

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortner/controllers/util"

	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const PROVIDER = "google"

func HandleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	slog.Info("inside HandleGoogleAuth")

	// TODO: Check if the godotenv can be load once in the main and shared across as a state
	envErr := godotenv.Load(".env")
	if envErr != nil {
		slog.Error("Error loading env file", "err", envErr)
		return
	}

	callBackUrl := util.GetCurrDomain(r) + os.Getenv("PORT") + r.RequestURI + "/callback?provider=" + PROVIDER
	slog.Debug("HandleGoogleAuth", "callBackUrl", callBackUrl)

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), callBackUrl),
	)

	gothic.BeginAuthHandler(w, r)
}
