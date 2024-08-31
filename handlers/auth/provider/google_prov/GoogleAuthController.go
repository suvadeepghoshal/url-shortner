package google_prov

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortner/handlers/util"
	model "url-shortner/model/type"

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

	q := r.URL.Query()
	q.Add("provider", "google")
	r.URL.RawQuery = q.Encode()

	callBackUrl := util.GetCurrDomain(r) + os.Getenv("PORT") + r.RequestURI + "/callback?provider=" + PROVIDER
	slog.Debug("HandleGoogleAuth", "callBackUrl", callBackUrl)

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), callBackUrl),
	)

	gothic.BeginAuthHandler(w, r)

}

func HandleGoogleAuthCallBack(w http.ResponseWriter, r *http.Request) {
	slog.Info("inside HandleGoogleAuthCallBack")

	resp, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		slog.Error("Unable to get context from auth provider", "err", err)
		return
	}

	// do we need a pointer? not setting any data just creating an user object to store val in session. Might need if we consider to store the entire user is session and then modify from somewhere else
	var authedUser model.User

	if resp.RawData != nil {
		m := resp.RawData
		authedUser = model.User{
			ID:            m["id"].(string),
			Email:         m["email"].(string),
			Picture:       m["picture"].(string),
			VerifiedEmail: m["verified_email"].(bool),
		}
		slog.Debug("HandleGoogleAuthCallBack", "authedUser", authedUser)

		if authedUser.VerifiedEmail {
			slog.Info("User is successfully verified")
		}

		session, sErr := gothic.Store.Get(r, "auth-session")
		if sErr != nil {
			slog.Error("Unable to generate auth session", "sErr", sErr)
			return
		}

		session.Values["access_token"] = resp.AccessToken
		session.Values["refresh_token"] = resp.RefreshToken
		session.Values["expires_at"] = resp.ExpiresAt
		session.Values["user_authed"] = authedUser.VerifiedEmail
		err := session.Save(r, w)
		if err != nil {
			slog.Error("Unable to save the session", "err", err)
			return
		}

		slog.Debug("HandleGoogleAuthCallBack", "session", session)

	}

	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//encErr := json.NewEncoder(w).Encode(fmt.Sprintf("%s is successfully logged in!", authedUser.Email))
	//if encErr != nil {
	//	slog.Error("Unable to write response", "err", encErr.Error())
	//	http.Error(w, "Unable to write response", http.StatusInternalServerError)
	//	return
	//}
}
