# Update Procedure

Steps needed to update varnish operator:

* change version.txt
* change varnish-operator/values.yaml#operator.varnishImage.tag
* change varnish-operator/values.yaml#operator.controllerImage.tag
* change varnish-operator/Chart.yaml#appVersion
* change varnish-operator/Chart.yaml#version
* change config/samples/icm_v1alpha1_varnishservice.yaml#spec.deployment.varnishImage.tag

Sure wish all of this was automated...
