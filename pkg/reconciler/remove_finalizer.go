package reconciler

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func sFnRemoveFinalizer(ctx context.Context, r *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	controllerutil.RemoveFinalizer(&s.instance, r.Finalizer)

	err := r.Update(ctx, &s.instance)
	if client.IgnoreNotFound(err) != nil {
		return nil, &ctrl.Result{Requeue: true}, err
	}

	r.log.Debug("finalizer removed")
	return nil, nil, nil
}
