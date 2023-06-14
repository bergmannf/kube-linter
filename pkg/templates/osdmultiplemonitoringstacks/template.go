package osdmultiplemonitoringstacks

import (
	"fmt"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/osdmultiplemonitoringstacks/internal/params"
)

func init() {
	templates.Register(check.Template{
		HumanName:              "Multiple Monitoring Operators in OSD cluster",
		Key:                    "osdmultiplemonitoringstacks",
		Description:            "Flag prometheus operator pods that are not running openshift-monitoring",
		SupportedObjectKinds:   config.ObjectKindsDesc{ObjectKinds: []string{objectkinds.DeploymentLike}},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return func(lintCtx lintcontext.LintContext, object lintcontext.Object) []diagnostic.Diagnostic {
				kind := object.K8sObject.GetObjectKind().GroupVersionKind().Kind
				// TODO: Might be better to check the images instead
				name := object.K8sObject.GetName()
				namespace := object.K8sObject.GetNamespace()
				if strings.Contains(name, "prometheus-operator") && namespace != "openshift-monitoring" {
					return []diagnostic.Diagnostic{{
						Message: fmt.Sprintf("%s: %s seems to run a prometheus-operator but is running in namespace '%s'. Multiple prometheus-operators are not supported in OSD.", kind, name, namespace),
					}}
				}
				return nil
			}, nil
		}),
	})
}
