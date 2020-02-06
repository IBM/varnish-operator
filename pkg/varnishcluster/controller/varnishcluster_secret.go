package controller

import (
	"context"
	"crypto/rand"
	icmapiv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	vclabels "icm-varnish-k8s-operator/pkg/labels"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/names"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	varnishDefaultSecretKeyName = "secret"
	varnishSecretSize           = 512
)

func (r *ReconcileVarnishCluster) reconcileVarnishSecret(ctx context.Context, instance *icmapiv1alpha1.VarnishCluster) error {
	secretName, secretKey := namesForInstanceSecret(instance)
	logr := logger.FromContext(ctx).With(logger.FieldComponent, icmapiv1alpha1.VarnishComponentSecret)
	logr = logr.With(logger.FieldComponent, secretName)

	secret := &v1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: instance.Namespace}, secret)
	if err != nil && kerrors.IsNotFound(err) {
		logr.Info("Creating Varnish Secret")
		secretLabels := vclabels.CombinedComponentLabels(instance, icmapiv1alpha1.VarnishComponentSecret)
		err = createNewSecret(
			secret,
			secretName,
			secretKey,
			instance.Namespace,
			secretLabels,
		)
		if err != nil {
			return errors.Wrap(err, "can't create secret object")
		}
		//set owner reference for the new created secret to clean up it on a varnish cluster remove.
		if err = controllerutil.SetControllerReference(instance, secret, r.scheme); err != nil {
			return errors.Wrap(err, "could not set secret owner reference")
		}
		if err := r.Create(ctx, secret); err != nil {
			return errors.Wrap(err, "unable to create secret")
		}
		return nil
	} else if err != nil {
		return errors.Wrap(err, "could not get existing varnish secret")
	}
	//check if existing k8s secret has valid password data. Update secret data if required.
	err = validateVarnishSecretContent(secret, secretKey)
	if err != nil {
		logr.Info("Updating Varnish Secret due empty data")
		return r.updateVarnishSecret(ctx, secret, secretKey)
	}
	logr.Debugw("No updates for Varnish Secret")
	return nil
}

func (r *ReconcileVarnishCluster) updateVarnishSecret(ctx context.Context, secret *v1.Secret, varnishSecretKeyName string) error {
	secretData, err := generateSecretData()
	if err != nil {
		return errors.Wrap(err, "could not generate secret data")
	}
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	secret.Data[varnishSecretKeyName] = secretData
	err = r.Update(ctx, secret)
	if err != nil {
		return errors.Wrap(err, "unable to update secret data")
	}
	return nil
}

func validateVarnishSecretContent(secret *v1.Secret, varnishSecretKeyName string) error {
	if len(secret.Data[varnishSecretKeyName]) > 0 {
		return nil
	}
	return errors.New("secret is empty")
}

func createNewSecret(secret *v1.Secret, secretName, varnishSecretKeyName, secretNamespace string, secretLabels map[string]string) error {
	secretData, err := generateSecretData()
	if err != nil {
		return errors.Wrap(err, "could not generate secret data")
	}
	*secret = v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: secretNamespace,
			Labels:    secretLabels,
		},
		Data: map[string][]byte{
			varnishSecretKeyName: secretData,
		},
	}
	return nil
}

func generateSecretData() ([]byte, error) {
	data := make([]byte, varnishSecretSize)
	_, err := rand.Read(data)
	return data, err
}

func namesForInstanceSecret(instance *icmapiv1alpha1.VarnishCluster) (secretName, secretKey string) {
	secretName, secretKey = names.VarnishSecret(instance.Name), varnishDefaultSecretKeyName
	spec := instance.Spec.Varnish.Secret
	if spec != nil && notEmptyString(spec.SecretName) {
		secretName = *spec.SecretName
		if notEmptyString(spec.Key) {
			secretKey = *spec.Key
		}
	}
	return
}

func notEmptyString(s *string) bool {
	return s != nil && *s != ""
}
