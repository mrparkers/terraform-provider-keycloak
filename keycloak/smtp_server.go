package keycloak

type SmtpServer struct {
	From               string `json:"from"`
	FromDisplayName    string `json:"fromDisplayName"`
	Host               string `json:"host"`
	ReplyTo            string `json:"replyTo"`
	ReplyToDisplayName string `json:"replyToDisplayName"`
	EnvelopeFrom       string `json:"envelopeFrom"`
	SSL                bool   `json:"ssl"`
	StartTLS           bool   `json:"starttls"`
	Port               int    `json:"port"`
	Authentication     bool   `json:"auth"`
	User               string `json:"user"`
	Password           string `json:"password"`
}
