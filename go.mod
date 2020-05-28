module github.com/nifcloud-lab/cluster-api-provider-nifcloud

go 1.13

require (
	github.com/aokumasan/nifcloud-sdk-go-v2 v0.0.5
	github.com/aws/aws-sdk-go v1.29.1
	github.com/aws/aws-sdk-go-v2 v0.15.0
	github.com/bramvdbogaerde/go-scp v0.0.0-20200119201711-987556b8bdd7
	github.com/chyeh/pubip v0.0.0-20170203095919-b7e679cf541c
	github.com/go-logr/logr v0.1.0
	github.com/golang/mock v1.4.1
	github.com/google/go-cmp v0.3.1
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/pkg/errors v0.9.1
	go.uber.org/multierr v1.1.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	k8s.io/api v0.0.0-20190918195907-bd6ac527cfd2
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20190918200256-06eb1244587a
	k8s.io/klog v0.4.0
	k8s.io/utils v0.0.0-20190809000727-6c36bc71fc4a
	sigs.k8s.io/cluster-api v0.2.7
	sigs.k8s.io/controller-runtime v0.4.0
)

replace (
	github.com/aokumasan/nifcloud-sdk-go-v2 v0.0.5 => github.com/donkomura/nifcloud-sdk-go-v2 v0.0.6-0.20200119061616-b429a3c824ff
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655 => k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/apimachinery v0.17.1 => k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible => k8s.io/client-go v0.0.0-20190918200256-06eb1244587a
	sigs.k8s.io/controller-runtime v0.4.0 => sigs.k8s.io/controller-runtime v0.3.0
)
