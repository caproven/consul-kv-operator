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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kvv1alpha1 "github.com/caproven/consul-kv-operator/api/v1alpha1"
)

const defaultRefresh = 60 * time.Second

// KVSecretReconciler reconciles a KVSecret object
type KVSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=consul-kv.caproven.info,resources=kvsecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=consul-kv.caproven.info,resources=kvsecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=consul-kv.caproven.info,resources=kvsecrets/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *KVSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	l.Info("Reconciling KVSecret")

	// TODO should use finalizers to stop syncing goroutine?

	kvSecret := &kvv1alpha1.KVSecret{}
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

	refresh := defaultRefresh
	if kvSecret.Spec.RefreshInterval != nil {
		refresh = time.Second * time.Duration(*kvSecret.Spec.RefreshInterval)
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
		Data: map[string][]byte{},
	}

	_, err = ctrl.CreateOrUpdate(ctx, r.Client, output, func() error {
		if err := ctrl.SetControllerReference(kvSecret, output, r.Scheme); err != nil {
			return err
		}
		for k, v := range data {
			output.Data[k] = v
		}
		return nil
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	l.Info("Finished reconciling KVSecret")

	return ctrl.Result{RequeueAfter: refresh}, nil
}

func secretData(cs *kvv1alpha1.KVSecret) (map[string][]byte, error) {
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

func lookupConsulKV(cs *kvv1alpha1.KVSecret, key string) ([]byte, error) {
	addr := fmt.Sprintf("%s:%d/v1/kv/%s", cs.Spec.Source.Host, cs.Spec.Source.Port, key)
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
	if len(body) == 0 {
		return []byte{}, nil
	}

	r := []kvResponse{}
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	if len(r) != 1 {
		return nil, fmt.Errorf("only expected 1 value, got %d", len(r))
	}

	val, err := base64.StdEncoding.DecodeString(r[0].Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode value from Consul: %v", err)
	}

	return val, nil
}

type kvResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SetupWithManager sets up the controller with the Manager.
func (r *KVSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kvv1alpha1.KVSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
