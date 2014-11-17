package clientpackets

type RequestAuthLogin struct {
	Username string
	Password string
}

func NewRequestAuthLogin(request []byte) RequestAuthLogin {
	var result RequestAuthLogin

	result.Username = string(request[:14])
	result.Password = string(request[14:28])

	return result
}
