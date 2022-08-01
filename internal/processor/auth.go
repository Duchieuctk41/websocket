package processor

import (
	"go_chat_tutorial/internal"
	pb "go_chat_tutorial/proto/ws-gateway"
	"strings"
	"time"
)

// authProcessor ...
type authProcessor struct {
	//httpClient *http.Client
	authUrl string
	nowFunc func() time.Time
}

// NewAuth ...
func NewAuth(authUrl string) internal.Processor {
	return &authProcessor{
		authUrl: authUrl,
	}
}

type authRequest struct {
	Token string `json:"token"`
}

func (p *authProcessor) Handle(m *pb.Command, client internal.Client) error {
	//xid := uuid.New().String()
	//tag := "[authProcessor.Handler] "
	//log := logger.WithField("x-request-id", xid)

	if m.GetAuth() == nil {
		return nil
	}
	token := m.GetAuth().GetToken()

	//reqBody := &authRequest{Token: token}
	//reqBodyBytes, err := json.Marshal(reqBody)
	//if err != nil {
	//	log.Errorf(tag+"error while marshalling auth request body: %v", err)
	//	return nil
	//}
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//defer cancel()
	//req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.authUrl, bytes.NewReader(reqBodyBytes))
	//if err != nil {
	//	log.Errorf(tag+"failed to create auth request: %v", err)
	//	return nil
	//}
	//
	//req.Header.Set("content-type", "application/json")
	//req.Header.Set("user-agent", "ws-gateway/v1")
	//req.Header.Set("x-request-id", xid)
	//
	//resp, err := p.httpClient.Do(req)
	//if err != nil {
	//	log.Errorf(tag+"failed to execute http request: %v", err)
	//	return nil
	//}
	//defer func() {
	//	_ = resp.Body.Close()
	//}()

	//respBody, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Errorf(tag+"failed to read auth response body: %v", err)
	//}
	//if resp.StatusCode != 200 {
	//	log.Errorf(tag+"failed to authorize token, auth response: (code: `%d`, body: `%s`)", resp.StatusCode, string(respBody))
	//	return nil
	//}

	userID := p.extractUserID(token)
	client.SetUserID(userID)

	// Todo: is online

	return nil
}

// extractUserID gets UserID part in user block
// on some sub-system, we store userID as a part of user block
// e.g. user123|device1
func (p *authProcessor) extractUserID(input string) string {
	t := strings.Split(input, "|")
	return t[0]
}
