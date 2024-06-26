package logparser

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	telemetryv1alpha1 "github.com/kyma-project/telemetry-manager/apis/telemetry/v1alpha1"
	"github.com/kyma-project/telemetry-manager/internal/conditions"
)

func (r *Reconciler) updateStatus(ctx context.Context, parserName string) error {
	var parser telemetryv1alpha1.LogParser
	if err := r.Get(ctx, types.NamespacedName{Name: parserName}, &parser); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}

		return fmt.Errorf("failed to get LogParser: %w", err)
	}

	if parser.DeletionTimestamp != nil {
		return nil
	}

	r.setAgentHealthyCondition(ctx, &parser)
	r.setLegacyConditions(ctx, &parser)

	if err := r.Status().Update(ctx, &parser); err != nil {
		return fmt.Errorf("failed to update LogParser status: %w", err)
	}
	return nil

}

func (r *Reconciler) setAgentHealthyCondition(ctx context.Context, parser *telemetryv1alpha1.LogParser) {
	healthy, err := r.prober.IsReady(ctx, r.config.DaemonSet)
	if err != nil {
		logf.FromContext(ctx).V(1).Error(err, "Failed to probe fluent bit daemonset - set condition as not healthy")
		healthy = false
	}

	status := metav1.ConditionFalse
	reason := conditions.ReasonAgentNotReady
	if healthy {
		status = metav1.ConditionTrue
		reason = conditions.ReasonAgentReady
	}

	condition := metav1.Condition{
		Type:               conditions.TypeAgentHealthy,
		Status:             status,
		Reason:             reason,
		Message:            conditions.MessageForLogPipeline(reason),
		ObservedGeneration: parser.Generation,
	}

	meta.SetStatusCondition(&parser.Status.Conditions, condition)
}

func (r *Reconciler) setLegacyConditions(ctx context.Context, parser *telemetryv1alpha1.LogParser) {
	fluentBitReady, err := r.prober.IsReady(ctx, r.config.DaemonSet)
	if err != nil {
		logf.FromContext(ctx).V(1).Error(err, "Failed to probe fluent bit daemonset")
		fluentBitReady = false
	}

	if !fluentBitReady {
		conditions.HandlePendingCondition(&parser.Status.Conditions, parser.Generation,
			conditions.ReasonFluentBitDSNotReady,
			conditions.MessageForLogPipeline(conditions.ReasonFluentBitDSNotReady))
		return
	}

	conditions.HandleRunningCondition(&parser.Status.Conditions, parser.Generation,
		conditions.ReasonFluentBitDSReady,
		conditions.ReasonFluentBitDSNotReady,
		conditions.MessageForLogPipeline(conditions.ReasonFluentBitDSReady),
		conditions.MessageForLogPipeline(conditions.ReasonFluentBitDSNotReady))
}
