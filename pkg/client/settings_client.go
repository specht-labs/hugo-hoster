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

// SettingsClient is a Kubernetes client for easy CRUD operations
type SettingsClient struct {
	client client.Client
	tracer trace.Tracer
}

// NewSettingsClient creates a new Setting Client
func NewSettingsClient(client client.Client, tracer trace.Tracer) *SettingsClient {
	return &SettingsClient{
		client: client,
		tracer: tracer,
	}
}

// Get returns a Setting in the current namespace
func (c *SettingsClient) Get(ct context.Context, name string) (*v1alpha1.Setting, error) {
	ctx, span := c.tracer.Start(ct, "SettingsClient.Get", trace.WithAttributes(attribute.String("name", name)))
	defer span.End()

	// try to read the namespace from /var/run
	namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		span.RecordError(err)
		return nil, errors.Wrap(err, "Unable to read current namespace")
	}

	return c.GetNamespaced(ctx, types.NamespacedName{Name: name, Namespace: string(namespace)})
}

// GetNameNamespace returns a Setting for a given name in a given namespace
func (c *SettingsClient) GetNameNamespace(ct context.Context, name, namespace string) (*v1alpha1.Setting, error) {
	ctx, span := c.tracer.Start(ct, "SettingsClient.GetNameNamespace", trace.WithAttributes(attribute.String("name", name), attribute.String("namespace", namespace)))
	defer span.End()

	return c.GetNamespaced(ctx, types.NamespacedName{Name: name, Namespace: namespace})
}

// Get returns a Setting
func (c *SettingsClient) GetNamespaced(ct context.Context, nameNamespaced types.NamespacedName) (*v1alpha1.Setting, error) {
	ctx, span := c.tracer.Start(
		ct, "SettingsClient.GetNamespaced",
		trace.WithAttributes(
			attribute.String("name", nameNamespaced.Name),
			attribute.String("namespace", nameNamespaced.Namespace),
		),
	)
	defer span.End()

	Setting := &v1alpha1.Setting{}

	if err := c.client.Get(ctx, nameNamespaced, Setting); err != nil {
		span.RecordError(err)
		return nil, err
	}

	return Setting, nil
}

// List returns a list of all Settings in the current namespace
func (c *SettingsClient) List(ct context.Context) (*v1alpha1.SettingList, error) {
	ctx, span := c.tracer.Start(ct, "SettingsClient.List")
	defer span.End()

	// try to read the namespace from /var/run
	namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		span.RecordError(err)
		return nil, errors.Wrap(err, "Unable to read current namespace")
	}

	return c.ListNamespaced(ctx, string(namespace))
}

// ListNamespaced returns a list of all Settings in a namespace
func (c *SettingsClient) ListNamespaced(ct context.Context, namespace string) (*v1alpha1.SettingList, error) {
	ctx, span := c.tracer.Start(ct, "SettingsClient.ListNamespaced", trace.WithAttributes(attribute.String("namespace", namespace)))
	defer span.End()

	Settings := &v1alpha1.SettingList{}

	if err := c.client.List(ctx, Settings, &client.ListOptions{Namespace: namespace}); err != nil {
		span.RecordError(err)
		return nil, err
	}

	return Settings, nil
}

func (c *SettingsClient) Update(ct context.Context, Setting *v1alpha1.Setting) error {
	ctx, span := c.tracer.Start(ct, "SettingsClient.Update", trace.WithAttributes(attribute.String("Setting", Setting.ObjectMeta.Name), attribute.String("namespace", Setting.ObjectMeta.Namespace)))
	defer span.End()

	if err := c.client.Update(ctx, Setting); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (c *SettingsClient) Delete(ct context.Context, Setting *v1alpha1.Setting) error {
	ctx, span := c.tracer.Start(ct, "SettingsClient.Delete", trace.WithAttributes(attribute.String("name", Setting.Name), attribute.String("namespace", Setting.Namespace)))
	defer span.End()

	if err := c.client.Delete(ctx, Setting); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (c *SettingsClient) Create(ct context.Context, Setting *v1alpha1.Setting) error {
	ctx, span := c.tracer.Start(ct, "SettingsClient.Create", trace.WithAttributes(attribute.String("Setting", Setting.ObjectMeta.Name), attribute.String("namespace", Setting.ObjectMeta.Namespace)))
	defer span.End()

	if Setting.Namespace == "" {
		// try to read the namespace from /var/run
		namespace, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			span.RecordError(err)
			return errors.Wrap(err, "Unable to read current namespace")
		}

		Setting.Namespace = string(namespace)
	}

	// if not exists, create a new one
	if err := c.client.Create(ctx, Setting); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
