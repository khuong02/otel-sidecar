package admission

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

var _ AdmissionController = &PodAdmissionRequestController{}

// PodAdmissionRequestController PodAdmissionRequest handler
type PodAdmissionRequestController struct {
	PodHandler PodPatcher
}

func (handler *PodAdmissionRequestController) handleAdmissionCreate(ctx context.Context, request *admissionv1.AdmissionRequest) ([]PatchOperation, error) {
	pod, err := unmarshalPod(request.Object.Raw)
	if err != nil {
		return nil, err
	}
	return handler.PodHandler.PatchPodCreate(ctx, request.Namespace, pod)
}

func (handler *PodAdmissionRequestController) handleAdmissionUpdate(ctx context.Context, request *admissionv1.AdmissionRequest) ([]PatchOperation, error) {
	oldPod, err := unmarshalPod(request.OldObject.Raw)
	if err != nil {
		return nil, err
	}
	newPod, err := unmarshalPod(request.Object.Raw)
	if err != nil {
		return nil, err
	}
	return handler.PodHandler.PatchPodUpdate(ctx, request.Namespace, oldPod, newPod)
}

func (handler *PodAdmissionRequestController) handleAdmissionDelete(ctx context.Context, request *admissionv1.AdmissionRequest) ([]PatchOperation, error) {
	pod, err := unmarshalPod(request.OldObject.Raw)
	if err != nil {
		return nil, err
	}
	return handler.PodHandler.PatchPodDelete(ctx, request.Namespace, pod)
}

func unmarshalPod(rawObject []byte) (corev1.Pod, error) {
	var pod corev1.Pod
	err := json.Unmarshal(rawObject, &pod)
	return pod, errors.Wrapf(err, "error unmarshalling object")
}
