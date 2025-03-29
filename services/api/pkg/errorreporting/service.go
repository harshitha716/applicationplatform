package errorreporting

import (
	"context"

	constants "github.com/Zampfi/application-platform/services/api/helper/constants"
	apicontext "github.com/Zampfi/application-platform/services/api/helper/context"
	"github.com/getsentry/sentry-go"
)

type Tag struct {
	Key   string
	Value string
}

func InitializeErrorReporting(dsn string, environment string) error {
	if environment == constants.ENVPRODUCTION || environment == constants.ENVSTAGING || environment == constants.ENVDEVELOPMENT {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:         dsn,
			Environment: environment,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func captureException(err error) {
	sentry.CaptureException(err)
}

func captureMessage(message string) {
	sentry.CaptureMessage(message)
}

func CaptureMessage(message string, ctx context.Context, tags ...Tag) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetContext("context", apicontext.GetAllCtxVariablesFromCtx(ctx))
		insertTags(scope, tags)
		captureMessage(message)
	})
}

func CaptureException(err error, ctx context.Context, tags ...Tag) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetContext("context", apicontext.GetAllCtxVariablesFromCtx(ctx))
		insertTags(scope, tags)
		captureException(err)
	})
}

func insertTags(scope *sentry.Scope, tags []Tag) {
	if tags != nil {
		for _, tag := range tags {
			scope.SetTag(tag.Key, tag.Value)
		}
	}
}
