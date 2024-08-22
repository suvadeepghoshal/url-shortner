package google_prov

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortner/controllers/util"
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

	callBackUrl := util.GetCurrDomain(r) + os.Getenv("PORT") + r.RequestURI + "/callback?provider=" + PROVIDER
	slog.Debug("HandleGoogleAuth", "callBackUrl", callBackUrl)

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), callBackUrl),
	)

	gothic.BeginAuthHandler(w, r)

}

func HandleGoogleAuthCallBack(w http.ResponseWriter, r *http.Request) {
	slog.Info("inside HandleGoogleAuthCallBack")

	// TODO: Do I need all of them?
	provider := r.URL.Query().Get("provider")
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	scope := r.URL.Query().Get("scope")
	authuser := r.URL.Query().Get("authuser")
	prompt := r.URL.Query().Get("prompt")

	slog.Debug("HandleGoogleAuthCallBack", "provider", provider)
	slog.Debug("HandleGoogleAuthCallBack", "state", state)
	slog.Debug("HandleGoogleAuthCallBack", "code", code)
	slog.Debug("HandleGoogleAuthCallBack", "scope", scope)
	slog.Debug("HandleGoogleAuthCallBack", "authuser", authuser)
	slog.Debug("HandleGoogleAuthCallBack", "prompt", prompt)

	resp, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		slog.Error("Unable to get context from auth provider", "err", err)
	} else {
		slog.Info("HandleGoogleAuthCallBack", "res", resp)
	}

	// Below is the sample response
	// What can be done?
	// Either use session to store the access_token and use middlewares to validate and then only grant access to protected resources
	// use the user's details to render user details if required
	// {RawData:map[email:ghoshalsuvadeep594@gmail.com id:110078008929846259901 picture:https://lh3.googleusercontent.com/a-/ALV-UjX3GxxuLpQyb_gON6mWadWpDBpJdxPPd1njHiJ8OLc6A2WN8i4s=s96-c verified_email:true] Provider:google Email:ghoshalsuvadeep594@gmail.com Name: FirstName: LastName: NickName: Description: UserID:110078008929846259901 AvatarURL:https://lh3.googleusercontent.com/a-/ALV-UjX3GxxuLpQyb_gON6mWadWpDBpJdxPPd1njHiJ8OLc6A2WN8i4s=s96-c Location: AccessToken:ya29.a0AcM612wkCiFMXmtdaYFBnYaFPDNzm4jZNATKvfTzThBXGHLiz_Qsqy_YtcBR-Vd2kADqD49KNaSu6EQSC8h5aWPTwj59QxTorzEbomnVAJa054IwwiSlKtKF-k1S1TlLxITqVVPlAS7T4EDjUH81cioG91iwfGF7_5jT-j_naCgYKAeQSARASFQHGX2Mi3COaJ94xcnrvdHtv3z9DHA0175 AccessTokenSecret: RefreshToken:1//03ewNqQMhKzTUCgYIARAAGAMSNwF-L9IrxhUDS_xQnRyuyLyq77CSl1_HeqpT2gN09QnV0VFc8eDaXWOadWi3ENVRTCs2TVyrV18 ExpiresAt:2024-08-22 13:41:24.071099 +0100 BST m=+3620.038058292 IDToken:eyJhbGciOiJSUzI1NiIsImtpZCI6ImE0OTM5MWJmNTJiNThjMWQ1NjAyNTVjMmYyYTA0ZTU5ZTIyYTdiNjUiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI0NDMwMDI2MDE0MjQtNmQ1Z2Nra2xxM2Z2ZTE4M2tsZzl0NDYyMTluN2g0OXUuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI0NDMwMDI2MDE0MjQtNmQ1Z2Nra2xxM2Z2ZTE4M2tsZzl0NDYyMTluN2g0OXUuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMTAwNzgwMDg5Mjk4NDYyNTk5MDEiLCJlbWFpbCI6Imdob3NoYWxzdXZhZGVlcDU5NEBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6ImxnNjBnVmFHbnVqRVd6VU9SNHJuUGciLCJpYXQiOjE3MjQzMjY4ODUsImV4cCI6MTcyNDMzMDQ4NX0.GYzuYHqXt65bd1BlUYHvCtdDvejEaTLmQqnc0_3w4VfXQpbVY4X3N4U856g7tl4vXwI9_QpK9OnPG5U7qYCoV1JsfIVPFFnPXUbCOku9b40uUXXdA3wX4Fy4nWDM5AZkG8b8iUEHX8InFNfkDMT9atgEdODMR3QTGAv1ISTguPr4WlZX5SpwQUrC281uVbnbsK_n7x9xgYAQ6mk9h9DBefMRNdL5TmMm1otlEcukCjW8jhTa1unWBpylgxeRmYmIZdYwnP12g46hprR50sumwzNJrh-MaI2ZrXdAv5DP_rWZwxAIxfP-wSf41ZyU_aWaCBQ17k4qHlHMWl6fPmepTQ
	if resp.RawData != nil {
		m := resp.RawData
		authedUser := model.User{
			ID:            m["id"].(string),
			Email:         m["email"].(string),
			Picture:       m["picture"].(string),
			VerifiedEmail: m["verified_email"].(bool),
		}
		slog.Debug("HandleGoogleAuthCallBack", "authedUser", authedUser)
	}
	// 2. In session: we will store the access_token, refresh_token, expires value (may be we will not expose everything)
}
