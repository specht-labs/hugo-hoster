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
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	//+kubebuilder:scaffold:imports

	pageClient "github.com/cedi/hugo-hoster/pkg/client"
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
	var settingsName string

	flag.StringVar(&settingsName, "settingName", "settings", "The name of the hugo-hoster/Setting resource used to configure this instance of hugo-hoster")
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&debug, "debug", false, "Turn on debug logging")

	flag.Parse()

	// Initialize Logging
	otelLogger, undo := observability.InitLogging(debug)
	defer undo()
	defer otelLogger.Sync()

	ctrl.SetLogger(zapr.NewLogger(otelzap.L().Logger))

	initLog := otelzap.L().Sugar()

	// Initialize Tracing (OpenTelemetry)
	traceProvider, tracer, err := observability.InitTracer(serviceName, serviceVersion)
	if err != nil {
		initLog.Errorw("failed initializing tracing",
			zap.Error(err),
		)
		os.Exit(1)
	}

	ctx, span := tracer.Start(context.Background(), "main.startManager")
	log := initLog.Ctx(ctx)

	defer func() {
		if err := traceProvider.Shutdown(ctx); err != nil {
			observability.RecordError(&log, span, err, "Error shutting down tracer provider")
		}
	}()

	_, span = tracer.Start(context.Background(), "main.startManager")

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: metricsAddr,
		},
		HealthProbeBindAddress:        probeAddr,
		LeaderElection:                false,
		LeaderElectionID:              "9123ff57.cedi.dev",
		LeaderElectionReleaseOnCancel: false,
	})

	if err != nil {
		observability.RecordError(&log, span, err, "Unable to start hugo-hoster")
		os.Exit(1)
	}

	hugoPageClient := pageClient.NewHugoPageClient(
		mgr.GetClient(),
		tracer,
	)

	settingClient := pageClient.NewSettingsClient(
		mgr.GetClient(),
		tracer,
	)

	hugoPageController := controllers.NewHugoPageReconciler(
		mgr.GetClient(),
		hugoPageClient,
		settingClient,
		settingsName,
		mgr.GetScheme(),
		tracer,
	)

	if err = hugoPageController.SetupWithManager(mgr); err != nil {
		observability.RecordError(&log, span, err, "Unable to create controller")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	span.End()

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		observability.RecordError(&log, span, err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		observability.RecordError(&log, span, err, "unable to set up ready check")
		os.Exit(1)
	}

	log.Infof("starting Hugo-Hoster %s", serviceVersion)
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		observability.RecordError(&log, span, err, "unable running manager")
		os.Exit(1)
	}
}
