package facebook

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// FBError is the struct for the error response from facebook
type FBError struct {
	Error struct {
		Message      string `json:"message"`
		Type         string `json:"type"`
		Code         int64  `json:"code"`
		ErrorSubcode int64  `json:"error_subcode"`
		FbtraceID    string `json:"fbtrace_id"`
	} `json:"error"`
}

// FBResp is the struct for the success response from facebook
type FBResp struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// GetDetailsFromFB gets details of the user from access token
func GetDetailsFromFB(token string) (fbResp FBResp, err error) {
	var (
		fbErr FBError
		resp  *http.Response
	)
	// var token = "EMAWcWrmVuLS23gXJZB82qSqQjY3bJqrv4odhAfpxXZCRMT2yyDluBLsfnHrKuO1ZAgSQISK0E6zLMHxRTLM1fmOF7LIcOH5Erg6yAOQDhT4ZAwxnGwGWxbJC9QBO0955UdIXNts51fL43lZArBL361ztfxVHeJOBUZD"
	var url = fmt.Sprintf("https://graph.facebook.com/me?fields=id,email&access_token=%s", token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fbResp, err
	}
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return fbResp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		_ = json.NewDecoder(resp.Body).Decode(&fbErr)
		// fmt.Println("ERROR", authErr.Error.Message)
		err = errors.New(fbErr.Error.Message)
		return fbResp, err
	}
	_ = json.NewDecoder(resp.Body).Decode(&fbResp)
	// fmt.Println("RESP", me.Phone.Number)
	return fbResp, err

}
