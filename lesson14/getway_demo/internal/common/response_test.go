package common

import "testing"

func TestInfoResponseContainsCoreFields(t *testing.T) {
	resp := BuildInfoResponse("api", "api-1", "/ping", "api.localtest.me", map[string]string{
		"X-Forwarded-Proto": "http",
	})

	if resp.Service != "api" {
		t.Fatalf("expected service api, got %s", resp.Service)
	}
	if resp.Instance != "api-1" {
		t.Fatalf("expected instance api-1, got %s", resp.Instance)
	}
	if resp.Path != "/ping" {
		t.Fatalf("expected path /ping, got %s", resp.Path)
	}
	if resp.Host != "api.localtest.me" {
		t.Fatalf("expected host api.localtest.me, got %s", resp.Host)
	}
	if resp.Forwarded["X-Forwarded-Proto"] != "http" {
		t.Fatalf("expected forwarded proto http, got %s", resp.Forwarded["X-Forwarded-Proto"])
	}
}

func TestBuildRequestInfoResponseContainsRuntimeFields(t *testing.T) {
	resp := BuildRequestInfoResponse(
		"echo",
		"echo-1",
		"GET",
		"/echo/inspect",
		"traefik.localtest.me",
		"127.0.0.1",
		map[string]string{
			"X-Forwarded-For":   "127.0.0.1",
			"X-Forwarded-Proto": "http",
		},
	)

	if resp.Service != "echo" {
		t.Fatalf("expected service echo, got %s", resp.Service)
	}
	if resp.Instance != "echo-1" {
		t.Fatalf("expected instance echo-1, got %s", resp.Instance)
	}
	if resp.Method != "GET" {
		t.Fatalf("expected method GET, got %s", resp.Method)
	}
	if resp.Path != "/echo/inspect" {
		t.Fatalf("expected path /echo/inspect, got %s", resp.Path)
	}
	if resp.Host != "traefik.localtest.me" {
		t.Fatalf("expected host traefik.localtest.me, got %s", resp.Host)
	}
	if resp.ClientIP != "127.0.0.1" {
		t.Fatalf("expected client ip 127.0.0.1, got %s", resp.ClientIP)
	}
	if resp.Forwarded["X-Forwarded-For"] != "127.0.0.1" {
		t.Fatalf("expected forwarded for 127.0.0.1, got %s", resp.Forwarded["X-Forwarded-For"])
	}
}
