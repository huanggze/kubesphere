package v1alpha2

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/kiali/kiali/handlers"
	"io/ioutil"
	"k8s.io/klog"
	"kubesphere.io/kubesphere/pkg/api"
	"net/http"
)

var JaegerQueryUrl = "http://jaeger-query.istio-system.svc:16686/jaeger"

// Get app metrics
func getAppMetrics(request *restful.Request, response *restful.Response) {
	handlers.AppMetrics(request, response)
}

// Get workload metrics
func getWorkloadMetrics(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	workload := request.PathParameter("workload")

	if len(namespace) > 0 && len(workload) > 0 {
		request.Request.URL.RawQuery = fmt.Sprintf("%s&namespaces=%s&workload=%s", request.Request.URL.RawQuery, namespace, workload)
	}

	handlers.WorkloadMetrics(request, response)
}

// Get service metrics
func getServiceMetrics(request *restful.Request, response *restful.Response) {
	handlers.ServiceMetrics(request, response)
}

// Get namespace metrics
func getNamespaceMetrics(request *restful.Request, response *restful.Response) {
	handlers.NamespaceMetrics(request, response)
}

// Get service graph for namespace
func getNamespaceGraph(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")

	if len(namespace) > 0 {
		request.Request.URL.RawQuery = fmt.Sprintf("%s&namespaces=%s", request.Request.URL.RawQuery, namespace)
	}

	handlers.GetNamespaceGraph(request, response)
}

// Get service graph for namespaces
func getNamespacesGraph(request *restful.Request, response *restful.Response) {
	handlers.GraphNamespaces(request, response)
}

// Get namespace health
func getNamespaceHealth(request *restful.Request, response *restful.Response) {
	handlers.NamespaceHealth(request, response)
}

// Get workload health
func getWorkloadHealth(request *restful.Request, response *restful.Response) {
	handlers.WorkloadHealth(request, response)
}

// Get app health
func getAppHealth(request *restful.Request, response *restful.Response) {
	handlers.AppHealth(request, response)
}

// Get service health
func getServiceHealth(request *restful.Request, response *restful.Response) {
	handlers.ServiceHealth(request, response)
}

func getServiceTracing(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	service := request.PathParameter("service")

	serviceName := fmt.Sprintf("%s.%s", service, namespace)

	url := fmt.Sprintf("%s/api/traces?%s&service=%s", JaegerQueryUrl, request.Request.URL.RawQuery, serviceName)

	resp, err := http.Get(url)

	if err != nil {
		klog.Errorf("query jaeger faile with err %v", err)
		api.HandleInternalError(response, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		klog.Errorf("read response error : %v", err)
		api.HandleInternalError(response, err)
		return
	}

	// need to set header for proper response
	response.Header().Set("Content-Type", "application/json")
	_, err = response.Write(body)

	if err != nil {
		klog.Errorf("write response failed %v", err)
	}
}