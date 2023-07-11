package mail

import (
	"testing"

	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("../../")
	require.NoError(t, err)

	sender := NewGmailSender(
		config.EmailSenderName,
		config.EmailSenderAddress,
		config.EmailSenderAddress)

	subject := "A test email"
	content := `
		<h1>Hello world</h1>
		<p>This is a test message from <a href="http://github.com/ak-karimzai">Ahmad Khalid Karimzai</a></p>
	`
	to := []string{"ak.karimzai@mail.ru", "ak.karimzai1@mail.ru"}
	attachFiles := []string{"../../readme.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
