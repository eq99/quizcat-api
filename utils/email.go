package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"net/textproto"
	"strings"

	"github.com/jordan-wright/email"

	"quizcat/conf"
)

func SendCaptchaByEmail(captcha, recipient string) error {
	// get infos from config
	from := conf.Conf().GetString("EMAIL_HOST_USER")
	password := conf.Conf().GetString("EMAIL_HOST_PASSWORD")
	to := []string{
		recipient,
	}
	host := conf.Conf().GetString("EMAIL_HOST")
	port := conf.Conf().GetString("EMAIL_PORT")

	// setup html message
	tmpl := fmt.Sprintf(`<table cellpadding="0" cellspacing="0" border="0" align="center">
    <tr>
        <td style="font-size: 15;">学习喵, 您的验证码是:</td>
    </tr>
    <tr>
        <td>
            <h1 style="margin: 5px;">%s</h1>
        </td>
    </tr>
    <tr>
        <td style="font-size: 15px;">十分钟有效。</td>
    </tr>
    </table>`, captcha)

	// Sending email.
	sender := &email.Email{
		To:      to,
		From:    from,
		Subject: "学习喵驾到",
		HTML:    []byte(tmpl),
		Headers: textproto.MIMEHeader{},
	}
	return sender.SendWithTLS(host+":"+port, smtp.PlainAuth("", from, password, host), &tls.Config{ServerName: host})
}

// Only following email domains are valid.
//
// If email is invalid, return true, other wise false
func IsEmailInvalid(email string) bool {
	domains := []string{"qq.com", "163.com", "gmail.com", "sina.com", "outlook.com"}
	split := strings.Split(email, "@")
	domain := split[len(split)-1]

	for _, _domain := range domains {
		if domain == _domain {
			return false
		}
	}

	return true
}
