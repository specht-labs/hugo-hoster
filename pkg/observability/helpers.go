package observability

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func NewZapLoggerWithCtxSpanPageName(name string, ctx context.Context, span trace.Span, pageName string) otelzap.SugaredLoggerWithCtx {
	return otelzap.L().Sugar().With(
		zap.String("name", name),
		zap.String("span_id", span.SpanContext().SpanID().String()),
		zap.String("page_name", pageName),
	).Ctx(ctx)
}

func RecordPanic(zapLog *otelzap.SugaredLoggerWithCtx, span trace.Span, err error, msg string, args ...any) {
	span.RecordError(err)
	zapLog.Panicw(fmt.Sprintf(msg, args...), zap.Error(err), zap.String("span_id", span.SpanContext().SpanID().String()))
}

func RecordError(zapLog *otelzap.SugaredLoggerWithCtx, span trace.Span, err error, msg string, args ...any) error {
	message := fmt.Sprintf(msg, args...)

	zapLog.Errorw(message, zap.Error(err), zap.String("span_id", span.SpanContext().SpanID().String()), zap.String("trace_id", span.SpanContext().TraceID().String()))

	err = errors.Wrap(err, message)
	span.RecordError(err)

	return err
}

func RecordInfo(zapLog *otelzap.SugaredLoggerWithCtx, span trace.Span, msg string, args ...any) {
	message := fmt.Sprintf(msg, args...)

	span.AddEvent(message)
	zapLog.Infow(message, zap.String("span_id", span.SpanContext().SpanID().String()))
}
