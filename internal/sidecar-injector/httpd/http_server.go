package httpd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"tracing/internal/sidecar-injector/admission"
	"tracing/internal/sidecar-injector/webhook"

	"github.com/gofiber/fiber/v2/log"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type HTTPServer struct {
	Local    bool
	Port     int
	CertFile string
	KeyFile  string
	Patcher  webhook.SidecarInjectorPatcher
	Debug    bool
}

func (s *HTTPServer) Start() error {
	k8sClient, err := s.CreateClient()
	if err != nil {
		return err
	}

	s.Patcher.K8sClient = k8sClient
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", s.Port),
	}

	mux := http.NewServeMux()
	server.Handler = mux

	admissionHandler := &admission.Handler{
		Controller: &admission.PodAdmissionRequestController{
			PodHandler: &s.Patcher,
		},
	}
	mux.HandleFunc("/healthz", webhook.HealthCheckHandler)
	mux.HandleFunc("/mutate", admissionHandler.HandleAdmission)

	if s.Local {
		return server.ListenAndServe()
	}
	return server.ListenAndServeTLS(s.CertFile, s.KeyFile)
}

// CreateClient Create the server
func (s *HTTPServer) CreateClient() (*kubernetes.Clientset, error) {
	config, err := s.buildConfig()

	if err != nil {
		return nil, errors.Wrapf(err, "error setting up cluster config")
	}

	return kubernetes.NewForConfig(config)
}

func (s *HTTPServer) buildConfig() (*rest.Config, error) {
	if s.Local {
		log.Debug("Using local kubeconfig.")
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	log.Debug("Using in cluster kubeconfig.")
	return rest.InClusterConfig()
}
