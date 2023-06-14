package osdmultiplemonitoringstacks

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/osdmultiplemonitoringstacks/internal/params"
	v1 "k8s.io/api/apps/v1"
)

func TestOsdMultipleMonitoringStacks(t *testing.T) {
	suite.Run(t, new(OsdMultipleMonitoringStacksTestSuite))
}

type OsdMultipleMonitoringStacksTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *OsdMultipleMonitoringStacksTestSuite) SetupTest() {
	s.Init("osdmultiplemonitoringstacks")
	s.ctx = mocks.NewMockContext()
}

func (s *OsdMultipleMonitoringStacksTestSuite) addDeployment(name, namespace string) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *v1.Deployment) {
		deployment.Namespace = namespace
	})
}

func (s *OsdMultipleMonitoringStacksTestSuite) TestRouteTLS() {
	const (
		validNamespace   = "openshift-monitoring"
		invalidNamespace = "smm-system"
	)

	s.addDeployment("prometheus-operator-a", validNamespace)
	s.addDeployment("prometheus-operator-b", invalidNamespace)

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				"prometheus-operator-a": {},
				"prometheus-operator-b": {
					{Message: "Deployment: prometheus-operator-b seems to run a prometheus-operator but is running in namespace 'smm-system'. Multiple prometheus-operators are not supported in OSD."},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}
