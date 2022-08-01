package ws_gateway

// client uses this message to authorized with socket server
type AuthBody struct {
	Token string `json:"token,omitempty"`
}

type isCommand_Body interface {
	isCommand_Body()
}

type Command_Auth struct {
	Auth *AuthBody `json:"auth,omitempty"`
}

func (*Command_Auth) isCommand_Body() {}
func (*Command_Echo) isCommand_Body() {}

type Command struct {
	Body isCommand_Body `json:"body"`
}

func (m *Command) GetBody() isCommand_Body {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *Command) GetAuth() *AuthBody {
	if x, ok := m.GetBody().(*Command_Auth); ok {
		return x.Auth
	}
	return nil
}

func (m *AuthBody) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type Command_Echo struct {
	Echo *EchoBody `json:"echo,omitempty"`
}

func (m *Command) GetEcho() *EchoBody {
	if x, ok := m.GetBody().(*Command_Echo); ok {
		return x.Echo
	}
	return nil
}

func (m *EchoBody) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

type EchoBody struct {
	Body string `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}
