/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"go.uber.org/zap"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	hugohosterv1alpha1 "github.com/cedi/hugo-hoster/api/v1alpha1"
	"github.com/cedi/hugo-hoster/controllers"
	"github.com/cedi/hugo-hoster/pkg/observability"
	"github.com/go-logr/zapr"

	//+kubebuilder:scaffold:imports

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var (
	scheme         = runtime.NewScheme()
	serviceName    = "hugo-hoster"
	serviceVersion = "1.0.0"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(hugohosterv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var debug bool

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&debug, "debug", false, "Turn on debug logging")

	flag.Parse()

	// Initialize Logging
	otelLogger, undo := observability.InitLogging(debug)
	defer otelLogger.Sync()
	defer undo()

	ctrl.SetLogger(zapr.NewLogger(otelzap.L().Logger))

	// Initialize Tracing (OpenTelemetry)
	traceProvider, tracer, err := observability.InitTracer(serviceName, serviceVersion)
	if err != nil {
		otelzap.L().Sugar().Errorw("failed initializing tracing",
			zap.Error(err),
		)
		os.Exit(1)
	}

	ctx, span := tracer.Start(context.Background(), "main.startManager")

	defer func() {
		if err := traceProvider.Shutdown(ctx); err != nil {
			otelzap.L().Sugar().Errorw("Error shutting down tracer provider",
				zap.Error(err),
			)
		}
	}()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "9123ff57.cedi.dev",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		span.RecordError(err)
		otelzap.L().Sugar().Errorw("unable to start urlshortener",
			zap.Error(err),
		)
		os.Exit(1)
	}

	if err = (&controllers.HugoSiteReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		span.RecordError(err)
		otelzap.L().Sugar().Errorw("unable to create controller",
			zap.Error(err),
			zap.String("controller", "HugoSite"),
		)
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		otelzap.L().Sugar().Errorw("unable to set up health check",
			zap.Error(err),
		)
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		otelzap.L().Sugar().Errorw("unable to set up ready check",
			zap.Error(err),
		)
		os.Exit(1)
	}

	otelzap.L().Info("starting urlshortener")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		otelzap.L().Sugar().Errorw("unable running manager",
			zap.Error(err),
		)
		os.Exit(1)
	}
}
