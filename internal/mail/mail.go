package mail

import (
	"crypto/tls"
	"net/smtp"
)

// SendEmail 发送邮件
func SendEmail(from string, to string, password string, title string, body string) error {
	conn, err := tls.Dial("tcp", "smtp.qq.com:465", nil)
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, "smtp.qq.com")
	if err != nil {
		return err
	}
	a := smtp.PlainAuth("", from, password, "smtp.qq.com")
	err = c.Auth(a)
	if err != nil {
		return err
	}
	err = c.Mail(from)
	if err != nil {
		return err
	}
	err = c.Rcpt(to)
	if err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}

	header := make(map[string]string)
	header["From"] = "HTU_YDXG_Checkin"
	header["TO"] = to
	header["Subject"] = title
	header["Content-Type"] = "text/html;chartset=UTF-8"

	var smtpMsg string
	for k, v := range header {
		smtpMsg += k + ":" + v + "\r\n"
	}

	// 将正文拼接
	smtpMsg += "\r\n" + body

	_, err = w.Write([]byte(smtpMsg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	err = c.Quit()
	if err != nil {
		err = c.Close()
		if err != nil {
			return err
		}
		return err
	}
	return nil
}
