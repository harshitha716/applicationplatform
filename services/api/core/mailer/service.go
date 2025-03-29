package mailer

import (
	"context"
	"fmt"

	emailtemplates "github.com/Zampfi/application-platform/services/api/core/mailer/email_templates"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/Zampfi/application-platform/services/api/pkg/sparkpost"
	"go.uber.org/zap"
)

type MailerService interface {
	SendInvitationEmail(ctx context.Context, data InvitationEmailData) error
}

type mailerService struct {
	sparkpostClient    sparkpost.SparkPostClient
	emailTemplatesPath string
	fromEmail          string
	templater          emailtemplates.Templater
}

func NewMailerService(sparkpostClient sparkpost.SparkPostClient, fromEmail string, emailTemplatesPath string) MailerService {
	return &mailerService{
		sparkpostClient:    sparkpostClient,
		fromEmail:          fromEmail,
		emailTemplatesPath: emailTemplatesPath,
		templater:          emailtemplates.NewTemplater(emailTemplatesPath),
	}
}

func (m mailerService) SendInvitationEmail(ctx context.Context, data InvitationEmailData) error {

	ctxlogger := apicontext.GetLoggerFromCtx(ctx)

	invitationTemplate, err := m.templater.GetTemplate("send_membership_invitation", m.emailTemplatesPath, map[string]string{
		"invited_by_first_name": data.InvitedByFirstName,
		"organization_name":     data.OrganizationName,
		"invitation_link":       data.InvitationLink,
	})

	ctxlogger.Debug("invitation template", zap.String("template", invitationTemplate))

	if err != nil {
		ctxlogger.Error("failed to get invitation template", zap.Error(err))
		return err
	}

	err = m.sparkpostClient.SendEmail(ctx, m.fromEmail, fmt.Sprintf("%s invited you to join Zamp", data.InvitedByFirstName), invitationTemplate, []string{data.RecipientEmail})
	if err != nil {
		ctxlogger.Error("failed to send invitation email", zap.Error(err))
		return err
	}

	return nil
}
