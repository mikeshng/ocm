// Copyright Contributors to the Open Cluster Management project
package operator

import (
	"context"
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"

	operatorapiv1 "open-cluster-management.io/api/operator/v1"
)

var _ = ginkgo.Describe("Klusterlet hosting-cluster report validation", func() {
	createKlusterlet := func(spec operatorapiv1.KlusterletSpec) error {
		_, err := operatorClient.OperatorV1().Klusterlets().Create(context.Background(), &operatorapiv1.Klusterlet{
			ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("hosting-report-validation-%s", rand.String(5))},
			Spec:       spec,
		}, metav1.CreateOptions{})
		return err
	}

	ginkgo.It("rejects reporting without clusterName in Default mode", func() {
		err := createKlusterlet(operatorapiv1.KlusterletSpec{
			DeployOption: operatorapiv1.KlusterletDeployOption{
				Mode:                 operatorapiv1.InstallModeDefault,
				ReportHostingCluster: operatorapiv1.ReportHostingClusterModeEnable,
			},
		})
		gomega.Expect(apierrors.IsInvalid(err)).To(gomega.BeTrue())
		gomega.Expect(err).To(gomega.MatchError(gomega.ContainSubstring("spec.clusterName is required")))
	})

	ginkgo.It("rejects reporting without managementClusterName in Hosted mode", func() {
		err := createKlusterlet(operatorapiv1.KlusterletSpec{
			DeployOption: operatorapiv1.KlusterletDeployOption{
				Mode:                 operatorapiv1.InstallModeHosted,
				ReportHostingCluster: operatorapiv1.ReportHostingClusterModeEnable,
				Hosted:               &operatorapiv1.KlusterletHostedConfiguration{},
			},
		})
		gomega.Expect(apierrors.IsInvalid(err)).To(gomega.BeTrue())
		gomega.Expect(err).To(gomega.MatchError(gomega.ContainSubstring("spec.deployOption.hosted.managementClusterName is required")))
	})

	ginkgo.It("rejects reporting when ClusterClaim is disabled", func() {
		err := createKlusterlet(operatorapiv1.KlusterletSpec{
			ClusterName: "spoke-a",
			DeployOption: operatorapiv1.KlusterletDeployOption{
				Mode:                 operatorapiv1.InstallModeDefault,
				ReportHostingCluster: operatorapiv1.ReportHostingClusterModeEnable,
			},
			RegistrationConfiguration: &operatorapiv1.RegistrationConfiguration{
				FeatureGates: []operatorapiv1.FeatureGate{{
					Feature: "ClusterClaim",
					Mode:    operatorapiv1.FeatureGateModeTypeDisable,
				}},
			},
		})
		gomega.Expect(apierrors.IsInvalid(err)).To(gomega.BeTrue())
		gomega.Expect(err).To(gomega.MatchError(gomega.ContainSubstring("ClusterClaim cannot be disabled")))
	})
})
