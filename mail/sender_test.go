package mail

import (
	"testing"

	"github.com/dbracic21-foi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test subject"
	content := `
	<h1>Test content</h1>
	<p>Test paragraph</p>`
	to := []string{"dariobracic5@gmail.com"}
	attachFiles := []string{"../README.md"}
	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)

}
