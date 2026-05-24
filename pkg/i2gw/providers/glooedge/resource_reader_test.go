/*
Copyright The Kubernetes Authors.

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

package glooedge

import (
	"context"
	"strings"
	"testing"

	"github.com/kgateway-dev/ingress2gateway/pkg/i2gw"
	"k8s.io/apimachinery/pkg/types"
)

func TestReadResourcesFromFileFiltersNamespace(t *testing.T) {
	reader := newResourceReader(&i2gw.ProviderConf{Namespace: "default"})
	storage, err := reader.readResourcesFromFile(strings.NewReader(`
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: default-vs
  namespace: default
spec:
  virtualHost:
    domains:
    - default.example.com
    routes:
    - matchers:
      - prefix: /
      routeAction:
        single:
          upstream:
            name: default-upstream
            namespace: default
---
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: production-vs
  namespace: production
spec:
  virtualHost:
    domains:
    - production.example.com
    routes:
    - matchers:
      - prefix: /
      routeAction:
        single:
          upstream:
            name: production-upstream
            namespace: production
`))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(storage.VirtualServices) != 1 {
		t.Fatalf("Expected 1 VirtualService, got %d: %+v", len(storage.VirtualServices), storage.VirtualServices)
	}
	if _, ok := storage.VirtualServices[types.NamespacedName{Namespace: "default", Name: "default-vs"}]; !ok {
		t.Fatalf("Expected default/default-vs in storage, got %+v", storage.VirtualServices)
	}
	if _, ok := storage.VirtualServices[types.NamespacedName{Namespace: "production", Name: "production-vs"}]; ok {
		t.Fatalf("Expected production/production-vs to be filtered out")
	}
}

func TestReadUpstreamsFromClusterReturnsUnsupportedError(t *testing.T) {
	_, err := readUpstreamsFromCluster(context.Background(), nil)
	if err == nil {
		t.Fatal("Expected unsupported Upstream read error")
	}
	if !strings.Contains(err.Error(), "not supported yet, see issue #112") {
		t.Fatalf("Expected issue #112 unsupported error, got %v", err)
	}
}