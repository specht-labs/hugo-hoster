package client

import (
	"context"
	"os"

	"github.com/cedi/hugo-hoster/api/v1alpha1"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HugoPageClient is a Kubernetes client for easy CRUD operations
type HugoPageClient struct {
	client client.Client
	tracer trace.Tracer
}

// NewHugoPageClient creates a new HugoPage Client
func NewHugoPageClient(client client.Client, tracer trace.Tracer) *HugoPageClient {
	return &HugoPageClient{
		client: client,
		tracer: tracer,
	}
}

// Get returns a HugoPage in the current namespace
func (c *HugoPageClient) Get(ct context.Context, name string) (*v1alpha1.HugoPage, error) {
	ctx, span := c.tracer.Start(ct, "HugoPageClient.Get", trace.WithAttributes(attribute.String("name", name)))
	defer span.End()

	// try to read the namespace from /var/run
	namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		span.RecordError(err)
		return nil, errors.Wrap(err, "Unable to read current namespace")
	}

	return c.GetNamespaced(ctx, types.NamespacedName{Name: name, Namespace: string(namespace)})
}

// GetNameNamespace returns a HugoPage for a given name in a given namespace
func (c *HugoPageClient) GetNameNamespace(ct context.Context, name, namespace string) (*v1alpha1.HugoPage, error) {
	ctx, span := c.tracer.Start(ct, "HugoPageClient.GetNameNamespace", trace.WithAttributes(attribute.String("name", name), attribute.String("namespace", namespace)))
	defer span.End()

	return c.GetNamespaced(ctx, types.NamespacedName{Name: name, Namespace: namespace})
}

// Get returns a HugoPage
func (c *HugoPageClient) GetNamespaced(ct context.Context, nameNamespaced types.NamespacedName) (*v1alpha1.HugoPage, error) {
	ctx, span := c.tracer.Start(
		ct, "HugoPageClient.GetNamespaced",
		trace.WithAttributes(
			attribute.String("name", nameNamespaced.Name),
			attribute.String("namespace", nameNamespaced.Namespace),
		),
	)
	defer span.End()

	HugoPage := &v1alpha1.HugoPage{}

	if err := c.client.Get(ctx, nameNamespaced, HugoPage); err != nil {
		span.RecordError(err)
		return nil, err
	}

	return HugoPage, nil
}

// List returns a list of all HugoPages in the current namespace
func (c *HugoPageClient) List(ct context.Context) (*v1alpha1.HugoPageList, error) {
	ctx, span := c.tracer.Start(ct, "HugoPageClient.List")
	defer span.End()

	// try to read the namespace from /var/run
	namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		span.RecordError(err)
		return nil, errors.Wrap(err, "Unable to read current namespace")
	}

	return c.ListNamespaced(ctx, string(namespace))
}

// ListNamespaced returns a list of all HugoPages in a namespace
func (c *HugoPageClient) ListNamespaced(ct context.Context, namespace string) (*v1alpha1.HugoPageList, error) {
	ctx, span := c.tracer.Start(ct, "HugoPageClient.ListNamespaced", trace.WithAttributes(attribute.String("namespace", namespace)))
	defer span.End()

	HugoPages := &v1alpha1.HugoPageList{}

	if err := c.client.List(ctx, HugoPages, &client.ListOptions{Namespace: namespace}); err != nil {
		span.RecordError(err)
		return nil, err
	}

	return HugoPages, nil
}

func (c *HugoPageClient) Update(ct context.Context, HugoPage *v1alpha1.HugoPage) error {
	ctx, span := c.tracer.Start(ct, "HugoPageClient.Update", trace.WithAttributes(attribute.String("HugoPage", HugoPage.ObjectMeta.Name), attribute.String("namespace", HugoPage.ObjectMeta.Namespace)))
	defer span.End()

	if err := c.client.Update(ctx, HugoPage); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (c *HugoPageClient) UpdateStatus(ct context.Context, HugoPage *v1alpha1.HugoPage) error {
	ctx, span := c.tracer.Start(ct, "HugoPageClient.UpdateStatus", trace.WithAttributes(attribute.String("HugoPage", HugoPage.ObjectMeta.Name), attribute.String("namespace", HugoPage.ObjectMeta.Namespace)))
	defer span.End()

	err := c.client.Status().Update(ctx, HugoPage)
	if err != nil {
		span.RecordError(err)
	}

	return err
}

func (c *HugoPageClient) Delete(ct context.Context, HugoPage *v1alpha1.HugoPage) error {
	ctx, span := c.tracer.Start(ct, "HugoPageClient.Delete", trace.WithAttributes(attribute.String("name", HugoPage.Name), attribute.String("namespace", HugoPage.Namespace)))
	defer span.End()

	if err := c.client.Delete(ctx, HugoPage); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (c *HugoPageClient) Create(ct context.Context, HugoPage *v1alpha1.HugoPage) error {
	ctx, span := c.tracer.Start(ct, "HugoPageClient.Create", trace.WithAttributes(attribute.String("HugoPage", HugoPage.ObjectMeta.Name), attribute.String("namespace", HugoPage.ObjectMeta.Namespace)))
	defer span.End()

	if HugoPage.Namespace == "" {
		// try to read the namespace from /var/run
		namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			span.RecordError(err)
			return errors.Wrap(err, "Unable to read current namespace")
		}

		HugoPage.Namespace = string(namespace)
	}

	// if not exists, create a new one
	if err := c.client.Create(ctx, HugoPage); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
