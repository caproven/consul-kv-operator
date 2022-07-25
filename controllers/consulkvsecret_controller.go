/*
Copyright 2022.

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	consulkvv1alpha1 "github.com/caproven/consul-kv-operator/api/v1alpha1"
)

// ConsulKVSecretReconciler reconciles a ConsulKVSecret object
type ConsulKVSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=consul-kv.caproven.info,resources=consulkvsecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=consul-kv.caproven.info,resources=consulkvsecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=consul-kv.caproven.info,resources=consulkvsecrets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ConsulKVSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO should use finalizers to stop syncing goroutine?

	kvSecret := &consulkvv1alpha1.ConsulKVSecret{}
	err := r.Get(ctx, req.NamespacedName, kvSecret)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !kvSecret.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}

	secretName := kvSecret.Spec.Secret.Name
	if secretName == "" {
		secretName = kvSecret.GetName()
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConsulKVSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&consulkvv1alpha1.ConsulKVSecret{}).
		Complete(r)
}
