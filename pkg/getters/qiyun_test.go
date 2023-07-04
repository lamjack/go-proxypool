package getters

import (
	"encoding/json"
	"go-proxypool/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQiyun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		response := []string{"192.168.0.1:8080", "192.168.0.2:9090"}
		jsonResponse, _ := json.Marshal(response)
		_, _ = rw.Write(jsonResponse)
	}))
	defer server.Close()

	// Expected result.
	expected := []models.Ip{
		{
			Ip:   "192.168.0.1",
			Port: 8080,
		},
		{
			Ip:   "192.168.0.2",
			Port: 9090,
		},
	}

	ips, err := fetchAndParseIps(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(ips) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(ips))
	}

	for i, ip := range ips {
		if ip.Ip != expected[i].Ip || ip.Port != expected[i].Port {
			t.Errorf("expected %v, got %v", expected[i], ip)
		}
	}
}
