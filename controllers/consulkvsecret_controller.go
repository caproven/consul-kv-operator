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
	"errors"
	"fmt"
	"io"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	l := log.FromContext(ctx)

	l.Info("Reconciling ConsulKVSecret")

	// TODO should use finalizers to stop syncing goroutine?

	kvSecret := &consulkvv1alpha1.ConsulKVSecret{}
	err := r.Get(ctx, req.NamespacedName, kvSecret)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !kvSecret.GetDeletionTimestamp().IsZero() {
		return ctrl.Result{}, nil
	}

	secretName := kvSecret.Spec.Output.Name
	if secretName == "" {
		secretName = kvSecret.GetName()
	}

	data, err := secretData(kvSecret)
	if err != nil {
		l.Error(err, "Failed to fetch secret data from Consul")
		return ctrl.Result{}, err
	}

	output := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: kvSecret.GetNamespace(),
		},
		Data: data,
	}
	ctrl.SetControllerReference(kvSecret, output, r.Scheme)

	err = r.Create(ctx, output)
	if k8serrors.IsAlreadyExists(err) {
		err = r.Update(ctx, output)
		if err != nil {
			l.Error(err, "Failed to update secret")
			return ctrl.Result{}, err
		}
	} else {
		l.Error(err, "Failed to create secret")
		return ctrl.Result{}, err
	}

	l.Info("Finished reconciling ConsulKVSecret")

	return ctrl.Result{}, nil
}

func secretData(cs *consulkvv1alpha1.ConsulKVSecret) (map[string][]byte, error) {
	d := make(map[string][]byte)
	for _, value := range cs.Spec.Values {
		v, err := lookupConsulKV(cs, value.SourceKey)
		if err != nil {
			return nil, err
		}
		d[value.Key] = v
	}
	return d, nil
}

func lookupConsulKV(cs *consulkvv1alpha1.ConsulKVSecret, key string) ([]byte, error) {
	// TODO how to enable https?
	addr := fmt.Sprintf("http://%s:%d/v1/kv/%s", cs.Spec.Source.Host, cs.Spec.Source.Port, key)
	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return nil, errors.New("failed to get data from Consul")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConsulKVSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&consulkvv1alpha1.ConsulKVSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
