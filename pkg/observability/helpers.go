package observability

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func RecordPanic(ctx context.Context, span trace.Span, zapLog *otelzap.SugaredLogger, err error, msg string, args ...any) {
	span.RecordError(err)
	zapLog.Ctx(ctx).Panicw(fmt.Sprintf(msg, args...), zap.Error(err), zap.String("span_id", span.SpanContext().SpanID().String()))
}

func RecordError(ctx context.Context, span trace.Span, zapLog *otelzap.SugaredLogger, err error, msg string, args ...any) error {
	message := fmt.Sprintf(msg, args...)

	zapLog.Ctx(ctx).Errorw(message, zap.Error(err), zap.String("span_id", span.SpanContext().SpanID().String()))

	err = errors.Wrap(err, message)
	span.RecordError(err)

	return err
}

func RecordInfo(ctx context.Context, span trace.Span, zapLog *otelzap.SugaredLogger, msg string, args ...any) {
	message := fmt.Sprintf(msg, args...)

	span.AddEvent(message)
	zapLog.Ctx(ctx).Infow(message, zap.String("span_id", span.SpanContext().SpanID().String()))
}
