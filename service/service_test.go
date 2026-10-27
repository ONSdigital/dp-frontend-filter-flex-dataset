package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/filter"
	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/config"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/service"
	"github.com/ONSdigital/dp-frontend-filter-flex-dataset/service/mocks"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-net/v3/request"
	. "github.com/smartystreets/goconvey/convey"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	ctx = context.Background()

	errAddCheckFail = errors.New("failed to register checkers: Error(s) registering checkers for healthcheck")
	errHealthCheck  = errors.New("healthCheck error")
	errServer       = errors.New("HTTP Server error")

	// Health Check Mock
	hcMock = &mocks.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}
	hcMockAddFail = &mocks.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return errAddCheckFail },
		StartFunc:    func(ctx context.Context) {},
	}
	funcDoGetHealthCheckOK = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
		return hcMock, nil
	}
	funcDoGetHealthCheckFail = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
		return nil, errHealthCheck
	}
	funcDoGetHealthAddCheckerFail = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
		return hcMockAddFail, nil
	}

	// Server Mock
	serverWg   = &sync.WaitGroup{}
	serverMock = &mocks.HTTPServerMock{
		ListenAndServeFunc: func() error {
			serverWg.Done()
			return nil
		},
	}
	failingServerMock = &mocks.HTTPServerMock{
		ListenAndServeFunc: func() error {
			serverWg.Done()
			return errServer
		},
	}
	funcDoGetHTTPServerOK = func(bindAddr string, router http.Handler) service.HTTPServer {
		return serverMock
	}
	funcDoGetHTTPServerFail = func(bindAddr string, router http.Handler) service.HTTPServer {
		return failingServerMock
	}

	// Health Client Mock
	funcDoGetHealthClient = func(name string, url string) *health.Client {
		return &health.Client{
			URL:    url,
			Name:   name,
			Client: service.NewMockHTTPClient(&http.Response{}, nil),
		}
	}
)

func TestNew(t *testing.T) {
	Convey("New returns a new uninitialised service", t, func() {
		So(service.New(), ShouldResemble, &service.Service{})
	})
}

func TestInitSuccess(t *testing.T) {
	Convey("Given all dependencies are successfully initialised", t, func() {
		initMock := &mocks.InitialiserMock{
			DoGetHealthClientFunc: funcDoGetHealthClient,
			DoGetHealthCheckFunc:  funcDoGetHealthCheckOK,
			DoGetHTTPServerFunc:   funcDoGetHTTPServerOK,
		}
		mockServiceList := service.NewServiceList(initMock)

		Convey("and valid config and service error channel are provided", func() {
			service.BuildTime = "TestBuildTime"
			service.GitCommit = "TestGitCommit"
			service.Version = "TestVersion"

			cfg, err := config.Get()
			So(err, ShouldBeNil)

			svc := &service.Service{}

			Convey("When Init is called", func() {
				err := svc.Init(ctx, cfg, mockServiceList)

				Convey("Then service is initialised successfully", func() {
					So(svc.Config, ShouldResemble, cfg)
					So(svc.HealthCheck, ShouldResemble, hcMock)
					So(svc.Server, ShouldResemble, serverMock)
					So(svc.ServiceList, ShouldResemble, mockServiceList)

					Convey("And returns no errors", func() {
						So(err, ShouldBeNil)

						Convey("And the checkers are registered and the healthcheck", func() {
							So(mockServiceList.HealthCheck, ShouldBeTrue)
							So(len(hcMock.AddCheckCalls()), ShouldEqual, 1)
							So(len(initMock.DoGetHTTPServerCalls()), ShouldEqual, 1)
							So(initMock.DoGetHTTPServerCalls()[0].BindAddr, ShouldEqual, "localhost:20100")
						})
					})
				})
			})
		})
	})
}

func TestInitFailure(t *testing.T) {
	Convey("Given failure to create healthcheck", t, func() {
		initMock := &mocks.InitialiserMock{
			DoGetHealthClientFunc: funcDoGetHealthClient,
			DoGetHealthCheckFunc:  funcDoGetHealthCheckFail,
		}
		mockServiceList := service.NewServiceList(initMock)

		Convey("and valid config and service error channel are provided", func() {
			service.BuildTime = "TestBuildTime"
			service.GitCommit = "TestGitCommit"
			service.Version = "TestVersion"

			cfg, err := config.Get()
			So(err, ShouldBeNil)

			svc := &service.Service{}

			Convey("When Init is called", func() {
				err := svc.Init(ctx, cfg, mockServiceList)

				Convey("Then service initialisation fails", func() {
					So(svc.Config, ShouldResemble, cfg)
					So(svc.ServiceList, ShouldResemble, mockServiceList)
					So(svc.ServiceList.HealthCheck, ShouldBeFalse)

					// Healthcheck and Server not initialised
					So(svc.HealthCheck, ShouldBeNil)
					So(svc.Server, ShouldBeNil)

					Convey("And returns error", func() {
						So(err, ShouldNotBeNil)
						So(err, ShouldBeError, "failed to create health check: healthCheck error")
					})
				})
			})
		})
	})

	Convey("Given that Checkers cannot be registered", t, func() {
		initMock := &mocks.InitialiserMock{
			DoGetHealthClientFunc: funcDoGetHealthClient,
			DoGetHealthCheckFunc:  funcDoGetHealthAddCheckerFail,
		}
		mockServiceList := service.NewServiceList(initMock)

		Convey("and valid config and service error channel are provided", func() {
			service.BuildTime = "TestBuildTime"
			service.GitCommit = "TestGitCommit"
			service.Version = "TestVersion"

			cfg, err := config.Get()
			So(err, ShouldBeNil)

			svc := &service.Service{}

			Convey("When Init is called", func() {
				err := svc.Init(ctx, cfg, mockServiceList)

				Convey("Then service initialisation fails", func() {
					So(svc.Config, ShouldResemble, cfg)
					So(svc.ServiceList, ShouldResemble, mockServiceList)
					So(svc.ServiceList.HealthCheck, ShouldBeTrue)
					So(svc.HealthCheck, ShouldResemble, hcMockAddFail)

					// Server not initialised
					So(svc.Server, ShouldBeNil)

					Convey("And returns error", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldResemble, errAddCheckFail.Error())

						Convey("And all checks try to register", func() {
							So(mockServiceList.HealthCheck, ShouldBeTrue)
							So(len(hcMockAddFail.AddCheckCalls()), ShouldEqual, 1)
							So(hcMockAddFail.AddCheckCalls()[0].Name, ShouldResemble, "API router")
						})
					})
				})
			})
		})
	})
}

func TestStart(t *testing.T) {
	Convey("Given a correctly initialised Service with mocked dependencies", t, func() {
		initMock := &mocks.InitialiserMock{
			DoGetHealthClientFunc: funcDoGetHealthClient,
			DoGetHealthCheckFunc:  funcDoGetHealthCheckOK,
			DoGetHTTPServerFunc:   funcDoGetHTTPServerOK,
		}
		serverWg.Add(1)
		mockServiceList := service.NewServiceList(initMock)

		svcErrors := make(chan error, 1)

		cfg, err := config.Get()
		So(err, ShouldBeNil)

		svc := &service.Service{
			Config:      cfg,
			HealthCheck: hcMock,
			Server:      serverMock,
			ServiceList: mockServiceList,
		}

		Convey("When service starts", func() {
			svc.Run(ctx, svcErrors)

			Convey("Then healthcheck is started and HTTP server starts listening", func() {
				So(len(hcMock.StartCalls()), ShouldEqual, 1)
				serverWg.Wait() // Wait for HTTP server go-routine to finish
				So(len(serverMock.ListenAndServeCalls()), ShouldEqual, 1)
			})
		})
	})

	Convey("Given that HTTP Server fails", t, func() {
		initMock := &mocks.InitialiserMock{
			DoGetHealthClientFunc: funcDoGetHealthClient,
			DoGetHealthCheckFunc:  funcDoGetHealthCheckOK,
			DoGetHTTPServerFunc:   funcDoGetHTTPServerFail,
		}
		serverWg.Add(1)
		mockServiceList := service.NewServiceList(initMock)

		Convey("and valid config and service error channel are provided", func() {
			service.BuildTime = "TestBuildTime"
			service.GitCommit = "TestGitCommit"
			service.Version = "TestVersion"

			svcErrors := make(chan error, 1)

			cfg, err := config.Get()
			So(err, ShouldBeNil)

			svc := &service.Service{
				Config:      cfg,
				HealthCheck: hcMock,
				Server:      failingServerMock,
				ServiceList: mockServiceList,
			}

			Convey("When service starts", func() {
				svc.Run(ctx, svcErrors)

				Convey("Then service start fails and returns an error in the error channel", func() {
					sErr := <-svcErrors
					So(sErr.Error(), ShouldResemble, errServer.Error())
					So(len(failingServerMock.ListenAndServeCalls()), ShouldEqual, 1)
				})
			})
		})
	})
}

func TestCloseSuccess(t *testing.T) {
	Convey("Given a correctly initialised service", t, func() {
		ctx := context.Background()

		cfg, err := config.Get()
		So(err, ShouldBeNil)

		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcCloseMock := &mocks.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverCloseMock := &mocks.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("Server stopped before healthcheck")
				}
				return nil
			},
		}

		serviceList := service.NewServiceList(nil)
		serviceList.HealthCheck = true
		svc := service.Service{
			Config:      cfg,
			HealthCheck: hcCloseMock,
			Server:      serverCloseMock,
			ServiceList: serviceList,
		}

		Convey("When closing service", func() {
			err = svc.Close(ctx)

			Convey("Then it results in all the dependencies being closed in the expected order", func() {
				So(err, ShouldBeNil)
				So(len(hcCloseMock.StopCalls()), ShouldEqual, 1)
				So(len(serverCloseMock.ShutdownCalls()), ShouldEqual, 1)
			})
		})
	})
}

func TestCloseFailure(t *testing.T) {
	Convey("Given if service fails to stop", t, func() {
		failingServerCloseMock := &mocks.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				return errors.New("Failed to stop http server")
			},
		}

		Convey("And given a correctly initialised service", func() {
			ctx := context.Background()

			cfg, err := config.Get()
			So(err, ShouldBeNil)

			serviceList := service.NewServiceList(nil)
			serviceList.HealthCheck = true
			svc := service.Service{
				Config:      cfg,
				HealthCheck: hcMock,
				Server:      failingServerCloseMock,
				ServiceList: serviceList,
			}

			Convey("When closing the service", func() {
				err = svc.Close(ctx)

				Convey("Then Close operation tries to close all dependencies", func() {
					So(len(hcMock.StopCalls()), ShouldEqual, 1)
					So(len(failingServerCloseMock.ShutdownCalls()), ShouldEqual, 1)

					Convey("And returns an error", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldResemble, "failed to shutdown gracefully")
					})
				})
			})
		})
	})

	Convey("Given that a dependency takes more time to close than the graceful shutdown timeout", t, func() {
		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcShutdownCloseMock := &mocks.HealthCheckerMock{
			StopFunc: func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverShutdownCloseMock := &mocks.HTTPServerMock{
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("Server was stopped before healthcheck")
				}
				return nil
			},
		}

		serverShutdownCloseMock.ShutdownFunc = func(ctx context.Context) error {
			time.Sleep(20 * time.Millisecond)
			return nil
		}

		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.GracefulShutdownTimeout = 1 * time.Millisecond

		Convey("And given a correctly initialised service", func() {
			ctx := context.Background()

			serviceList := service.NewServiceList(nil)
			serviceList.HealthCheck = true
			svc := service.Service{
				Config:      cfg,
				HealthCheck: hcShutdownCloseMock,
				Server:      serverShutdownCloseMock,
				ServiceList: serviceList,
			}

			Convey("When closing the service", func() {
				err = svc.Close(ctx)

				Convey("Then closing the service fails with context.DeadlineExceeded error and no further dependencies are attempted to close", func() {
					So(err, ShouldResemble, context.DeadlineExceeded)
					So(len(hcShutdownCloseMock.StopCalls()), ShouldEqual, 1)
					So(len(serverShutdownCloseMock.ShutdownCalls()), ShouldEqual, 1)
				})
			})
		})
	})
}

type backendRequest struct {
	method  string
	path    string
	query   string
	headers http.Header
	body    []byte
}

func TestSubmitBackendForwarding(t *testing.T) {
	const filterID = "filter-123"
	const requestID = "caller-123"
	const userToken = "synthetic-user-token"

	backend, received := newSubmitTestBackend(t)
	handler := newSubmitTestHandler(t, backend.URL)
	incoming := httptest.NewRequest(http.MethodPost, "/filters/"+filterID+"/submit?choice=one&choice=two", strings.NewReader("choice=from-form"))
	incoming.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	incoming.Header.Set("X-Request-Id", requestID)
	incoming.Header.Set("X-Florence-Token", userToken)
	incoming.Header.Set("Authorization", "Bearer unrelated-service-token")
	incoming.Header.Set("Collection-Id", "collection-123")
	incoming.Header.Set("X-Client-Only", "not-for-backend")
	incoming.AddCookie(&http.Cookie{Name: "session", Value: "synthetic-session"})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, incoming)
	if w.Code != http.StatusFound {
		t.Fatalf("submit status: got %d, want %d: %s", w.Code, http.StatusFound, w.Body.String())
	}

	lookup := <-received
	checkSubmitBackendRequest(t, lookup, http.MethodGet, "/v1/filters/"+filterID, requestID, userToken)
	if got := lookup.headers.Get("Collection-Id"); got != "collection-123" {
		t.Errorf("lookup collection ID: got %q, want %q", got, "collection-123")
	}
	if len(lookup.body) != 0 {
		t.Errorf("lookup body: got %q, want empty", lookup.body)
	}

	submit := <-received
	checkSubmitBackendRequest(t, submit, http.MethodPost, "/v1/filters/"+filterID+"/submit", requestID, userToken)
	if got := submit.headers.Get("If-Match"); got != "etag-123" {
		t.Errorf("submit If-Match: got %q, want %q", got, "etag-123")
	}
	if got := submit.headers.Get("Collection-Id"); got != "" {
		t.Errorf("submit collection ID: got %q, want empty", got)
	}
	if got := submit.headers.Get("Content-Type"); got == "application/x-www-form-urlencoded" {
		t.Error("submit forwarded the caller's form content type")
	}
	var payload filter.SubmitFilterRequest
	if err := json.Unmarshal(submit.body, &payload); err != nil {
		t.Fatalf("decode backend submit body: %v", err)
	}
	if payload.FilterID != filterID || payload.PopulationType != "type-123" || len(payload.Dimensions) != 0 {
		t.Errorf("backend submit body: got %+v, want filter ID and population type from lookup", payload)
	}
}

func TestSubmitBackendTraceWithOTel(t *testing.T) {
	t.Setenv("OTEL_PROPAGATORS", "tracecontext,baggage")
	previousProvider := otel.GetTracerProvider()
	previousPropagator := otel.GetTextMapPropagator()
	provider := sdktrace.NewTracerProvider()
	propagator := autoprop.NewTextMapPropagator()
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagator)
	t.Cleanup(func() {
		otel.SetTracerProvider(previousProvider)
		otel.SetTextMapPropagator(previousPropagator)
		if err := provider.Shutdown(context.Background()); err != nil {
			t.Errorf("shut down tracer provider: %v", err)
		}
	})

	backend, received := newSubmitTestBackend(t)
	handler := newSubmitTestHandler(t, backend.URL)
	const traceparent = "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01"
	const requestID = "caller-123"
	const userToken = "synthetic-user-token"
	incoming := httptest.NewRequest(http.MethodPost, "/filters/filter-123/submit", http.NoBody)
	incoming.Header.Set("traceparent", traceparent)
	incoming.Header.Set("X-Request-Id", requestID)
	incoming.Header.Set("X-Florence-Token", userToken)
	incomingSpan := trace.SpanContextFromContext(propagator.Extract(context.Background(), propagation.HeaderCarrier(incoming.Header)))
	if !incomingSpan.IsValid() {
		t.Fatalf("invalid test traceparent: %q", traceparent)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, incoming)
	if w.Code != http.StatusFound {
		t.Fatalf("submit status: got %d, want %d: %s", w.Code, http.StatusFound, w.Body.String())
	}

	for _, expected := range []struct{ method, path string }{
		{http.MethodGet, "/v1/filters/filter-123"},
		{http.MethodPost, "/v1/filters/filter-123/submit"},
	} {
		got := <-received
		checkSubmitBackendRequest(t, got, expected.method, expected.path, requestID, userToken)
		backendSpan := trace.SpanContextFromContext(propagator.Extract(context.Background(), propagation.HeaderCarrier(got.headers)))
		if !backendSpan.IsValid() {
			t.Errorf("backend %s traceparent: got %q, want a valid continuation of %q", expected.method, got.headers.Get("traceparent"), traceparent)
			continue
		}
		if backendSpan.TraceID() != incomingSpan.TraceID() {
			t.Errorf("backend %s trace ID: got %s, want %s", expected.method, backendSpan.TraceID(), incomingSpan.TraceID())
		}
		if backendSpan.SpanID() == incomingSpan.SpanID() {
			t.Errorf("backend %s span ID: got %s, want a new outbound span", expected.method, backendSpan.SpanID())
		}
	}
}

func newSubmitTestBackend(t *testing.T) (backend *httptest.Server, backendRequests <-chan backendRequest) {
	t.Helper()
	received := make(chan backendRequest, 2)
	backend = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read backend request: %v", err)
			http.Error(w, "could not read request", http.StatusBadRequest)
			return
		}
		received <- backendRequest{r.Method, r.URL.Path, r.URL.RawQuery, r.Header.Clone(), body}

		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("ETag", "etag-123")
			if _, err := io.WriteString(w, `{"filter_id":"filter-123","population_type":"type-123","dataset":{"id":"dataset-123","edition":"2021","version":1}}`); err != nil {
				t.Errorf("write filter response: %v", err)
			}
		case http.MethodPost:
			w.WriteHeader(http.StatusAccepted)
			if _, err := io.WriteString(w, `{"filter_output_id":"output-123"}`); err != nil {
				t.Errorf("write submit response: %v", err)
			}
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
		}
	}))
	t.Cleanup(backend.Close)
	return backend, received
}

func newSubmitTestHandler(t *testing.T, backendURL string) http.Handler {
	t.Helper()
	initialiser := &mocks.InitialiserMock{
		DoGetHealthClientFunc: new(service.Init).DoGetHealthClient,
		DoGetHealthCheckFunc: func(*config.Config, string, string, string) (service.HealthChecker, error) {
			return hcMock, nil
		},
		DoGetHTTPServerFunc: func(string, http.Handler) service.HTTPServer {
			return serverMock
		},
	}
	cfg := &config.Config{APIRouterURL: backendURL + "/v1", OTServiceName: "filter-flex-test"}
	if err := service.New().Init(context.Background(), cfg, service.NewServiceList(initialiser)); err != nil {
		t.Fatalf("initialise service: %v", err)
	}

	// The HTTP server adds request-ID middleware when it starts, outside Service.Init.
	return request.HandlerRequestID(16)(initialiser.DoGetHTTPServerCalls()[0].Router)
}

func checkSubmitBackendRequest(t *testing.T, got backendRequest, method, path, requestID, userToken string) {
	t.Helper()
	if got.method != method || got.path != path || got.query != "" {
		t.Errorf("backend request: got %s %s?%s, want %s %s with no caller query", got.method, got.path, got.query, method, path)
	}
	if value := got.headers.Get("X-Florence-Token"); value != userToken {
		t.Errorf("backend Florence token: got %q, want %q", value, userToken)
	}
	if value := got.headers.Get("Authorization"); value != "Bearer "+userToken {
		t.Errorf("backend authorization: got %q, want bearer user token", value)
	}
	if value := got.headers.Get("X-Request-Id"); !strings.HasPrefix(value, requestID+",") || len(value) != len(requestID)+1+len(requestID)/2 {
		t.Errorf("backend request ID: got %q, want %q followed by a generated ID", value, requestID+",")
	}
	if value := got.headers.Get("Cookie"); value != "" {
		t.Errorf("backend cookie: got %q, want none", value)
	}
	if value := got.headers.Get("X-Client-Only"); value != "" {
		t.Errorf("backend caller-only header: got %q, want none", value)
	}
}
