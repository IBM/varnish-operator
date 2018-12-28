# Update Procedure

Steps needed to update varnish operator:

* change icm-varnish-version.txt (if varnish image updated)
* change version.txt
* change varnish-operator/values.yaml#operator.varnishImage.tag (if varnish image updated)
* change varnish-operator/values.yaml#operator.controllerImage.tag
* change varnish-operator/Chart.yaml#appVersion
* change varnish-operator/Chart.yaml#version
* change config/samples/icm_v1alpha1_varnishservice.yaml#spec.deployment.varnishImage.tag (if varnish image updated)

Sure wish all of this was automated...
