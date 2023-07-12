package observability

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func RecordError(ctx context.Context, span trace.Span, zapLog *otelzap.SugaredLogger, err error, msg string, args ...any) error {
	message := fmt.Sprintf(msg, args...)
	span.AddEvent(message)

	zapLog.Ctx(ctx).Errorw(msg, zap.Error(err), zap.String("span_id", span.SpanContext().SpanID().String()))

	err = errors.Wrap(err, message)
	span.RecordError(err)
	return err
}

func RecordInfo(ctx context.Context, span trace.Span, zapLog *otelzap.SugaredLogger, msg string, args ...any) {
	zapLog.Ctx(ctx).Infow(fmt.Sprintf(msg, args...))
	span.AddEvent(fmt.Sprintf(msg, args...))
}
