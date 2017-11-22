Examples
========

## Run the controller

`./bin/catalog-controller`

## Create Custom Resource Definition

`kubectl create -f catalog-crd.yaml`

`kubectl create -f template-crd.yaml`

## Create catalog

`kubectl create -f catalog.yaml`

## View templates

`kubectl get templates`

```$xslt
$ kubectl get templates
NAME                                KIND
test-convoy-nfs                     Template.v1.catalog.cattle.io
test-infra-container-crontab        Template.v1.catalog.cattle.io
test-infra-ebs                      Template.v1.catalog.cattle.io
test-infra-ecr                      Template.v1.catalog.cattle.io
test-infra-efs                      Template.v1.catalog.cattle.io
test-infra-healthcheck              Template.v1.catalog.cattle.io
test-infra-ipsec                    Template.v1.catalog.cattle.io
test-infra-k8s                      Template.v1.catalog.cattle.io
test-infra-l2-flat                  Template.v1.catalog.cattle.io
test-infra-netapp-eseries           Template.v1.catalog.cattle.io
test-infra-netapp-ontap-nas         Template.v1.catalog.cattle.io
test-infra-netapp-ontap-san         Template.v1.catalog.cattle.io
test-infra-netapp-solidfire         Template.v1.catalog.cattle.io
test-infra-network-diagnostics      Template.v1.catalog.cattle.io
test-infra-network-policy-manager   Template.v1.catalog.cattle.io
test-infra-network-services         Template.v1.catalog.cattle.io
test-infra-nfs                      Template.v1.catalog.cattle.io
test-infra-per-host-subnet          Template.v1.catalog.cattle.io
test-infra-portworx                 Template.v1.catalog.cattle.io
test-infra-route53                  Template.v1.catalog.cattle.io
test-infra-scheduler                Template.v1.catalog.cattle.io
test-infra-secrets                  Template.v1.catalog.cattle.io
test-infra-vxlan                    Template.v1.catalog.cattle.io
test-infra-windows                  Template.v1.catalog.cattle.io
test-k8s                            Template.v1.catalog.cattle.io
test-kubernetes                     Template.v1.catalog.cattle.io
test-project-cattle                 Template.v1.catalog.cattle.io
test-project-kubernetes             Template.v1.catalog.cattle.io
test-project-mesos                  Template.v1.catalog.cattle.io
test-project-swarm                  Template.v1.catalog.cattle.io
test-project-windows                Template.v1.catalog.cattle.io
test-route53                        Template.v1.catalog.cattle.io
```