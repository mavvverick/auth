package accountkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// AuthError is the struct for the error resonse from FB
type AuthError struct {
	Error struct {
		Message   string `json:"message"`
		Type      string `json:"type"`
		Code      int    `json:"code"`
		FbtraceID string `json:"fbtrace_id"`
	} `json:"error"`
}

// Me is the struct for the success response from FB
type Me struct {
	Phone struct {
		Number         string `json:"number"`
		CountryPrefix  string `json:"country_prefix"`
		NationalNumber string `json:"national_number"`
	} `json:"phone"`
	ID          string `json:"id"`
	Application struct {
		ID string `json:"id"`
	} `json:"application"`
}

// GetDetailsFromFB gets details of the user from access token
func GetDetailsFromFB(token string) (me Me, err error) {
	var (
		authErr AuthError
		resp    *http.Response
	)
	// var token = "EMAWcWrmVuLS23gXJZB82qSqQjY3bJqrv4odhAfpxXZCRMT2yyDluBLsfnHrKuO1ZAgSQISK0E6zLMHxRTLM1fmOF7LIcOH5Erg6yAOQDhT4ZAwxnGwGWxbJC9QBO0955UdIXNts51fL43lZArBL361ztfxVHeJOBUZD"
	var url = fmt.Sprintf("https://graph.accountkit.com/v1.3/me?access_token=%s", token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return me, err
	}
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return me, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		_ = json.NewDecoder(resp.Body).Decode(&authErr)
		// fmt.Println("ERROR", authErr.Error.Message)
		err = errors.New(authErr.Error.Message)
		return me, err
	}
	_ = json.NewDecoder(resp.Body).Decode(&me)
	// fmt.Println("RESP", me.Phone.Number)
	return me, err

}
