package plek

import (
	"os"
	"testing"
)

type FindExample struct {
	GovukAppDomain string
	ServiceName    string
	ExpectedURL    string
	Environ        map[string]string
}

var findExamples = []FindExample{
	{
		GovukAppDomain: "example.com",
		ServiceName:    "foo",
		ExpectedURL:    "https://foo.example.com",
	},
	{
		GovukAppDomain: "example.com",
		ServiceName:    "foo.bar",
		ExpectedURL:    "https://foo.bar.example.com",
	},
	{ // dev.gov.uk domains should magically return http
		GovukAppDomain: "dev.gov.uk",
		ServiceName:    "foo",
		ExpectedURL:    "http://foo.dev.gov.uk",
	},
}

func TestFind(t *testing.T) {
	for i, ex := range findExamples {
		testFind(t, i, ex)
	}
}

func testFind(t *testing.T, i int, ex FindExample) {
	actual := New(ex.GovukAppDomain).Find(ex.ServiceName)
	expected := ex.ExpectedURL
	if actual != expected {
		t.Errorf("Example %d: expected %s, got %s", i, expected, actual)
	}
}

type FindURLExample struct {
	GovukAppDomain    string
	ServiceName       string
	ExpectedURLScheme string
	ExpectedURLHost   string
}

var findURLExamples = []FindURLExample{
	{
		GovukAppDomain:    "example.com",
		ServiceName:       "foo",
		ExpectedURLScheme: "https",
		ExpectedURLHost:   "foo.example.com",
	},
	{
		GovukAppDomain:    "example.com",
		ServiceName:       "foo.bar",
		ExpectedURLScheme: "https",
		ExpectedURLHost:   "foo.bar.example.com",
	},
	{ // dev.gov.uk domains should magically return http
		GovukAppDomain:    "dev.gov.uk",
		ServiceName:       "foo",
		ExpectedURLScheme: "http",
		ExpectedURLHost:   "foo.dev.gov.uk",
	},
}

func TestFindURL(t *testing.T) {
	for i, ex := range findURLExamples {
		testFindURL(t, i, ex)
	}
}

func testFindURL(t *testing.T, i int, ex FindURLExample) {
	actual := New(ex.GovukAppDomain).FindURL(ex.ServiceName)
	if actual.Host != ex.ExpectedURLHost {
		t.Errorf("Example %d: expected URL with host %s, got %s", i, ex.ExpectedURLHost, actual.Host)
	}
	if actual.Scheme != ex.ExpectedURLScheme {
		t.Errorf("Example %d: expected URL with scheme %s, got %s", i, ex.ExpectedURLScheme, actual.Scheme)
	}
}

var packageFindExamples = []FindExample{
	{
		GovukAppDomain: "example.com",
		ServiceName:    "foo",
		ExpectedURL:    "https://foo.example.com",
	},
	{
		GovukAppDomain: "",
		ServiceName:    "foo",
		ExpectedURL:    "http://foo.dev.gov.uk",
	},
	// Setting a hostname prefix
	{
		GovukAppDomain: "",
		ServiceName:    "foo",
		ExpectedURL:    "http://draft-foo.dev.gov.uk",
		Environ: map[string]string{
			"PLEK_HOSTNAME_PREFIX": "draft-",
		},
	},
	// Overriding a specific service URL with an ENV var.
	{
		GovukAppDomain: "foo.com",
		ServiceName:    "foo",
		ExpectedURL:    "http://foo.example.com",
		Environ:        map[string]string{"PLEK_SERVICE_FOO_URI": "http://foo.example.com"},
	},
	{
		GovukAppDomain: "foo.com",
		ServiceName:    "foo-bar",
		ExpectedURL:    "http://anything.example.com",
		Environ:        map[string]string{"PLEK_SERVICE_FOO_BAR_URI": "http://anything.example.com"},
	},
	{
		GovukAppDomain: "", // Should not be required when using overrides
		ServiceName:    "foo",
		ExpectedURL:    "http://foo.example.com",
		Environ:        map[string]string{"PLEK_SERVICE_FOO_URI": "http://foo.example.com"},
	},
	{
		GovukAppDomain: "foo.com",
		ServiceName:    "foo",
		ExpectedURL:    "http://invalid%hostname.com",
		Environ:        map[string]string{"PLEK_SERVICE_FOO_URI": "http://invalid%hostname.com"},
	},
	// PLEK_SERVICE_FOO_BAR_URI overrides the hostname prefix
	{
		GovukAppDomain: "foo.com",
		ServiceName:    "foo-bar",
		ExpectedURL:    "http://anything.example.com",
		Environ: map[string]string{
			"PLEK_SERVICE_FOO_BAR_URI": "http://anything.example.com",
			"PLEK_HOSTNAME_PREFIX":     "draft-",
		},
	},
}

func TestPackageFind(t *testing.T) {
	for i, ex := range packageFindExamples {
		testPackageFind(t, i, ex)
	}
}

func testPackageFind(t *testing.T, i int, ex FindExample) {
	os.Clearenv()
	for k, v := range ex.Environ {
		os.Setenv(k, v)
	}
	os.Setenv("GOVUK_APP_DOMAIN", ex.GovukAppDomain)

	actual, expected := Find(ex.ServiceName), ex.ExpectedURL
	if actual != expected {
		t.Errorf("Example %d: expected %s, got %s", i, expected, actual)
	}
}

func TestWebsiteRoot(t *testing.T) {
	os.Clearenv()
	os.Setenv("GOVUK_WEBSITE_ROOT", "https://www.gov.uk")

	actual, err := WebsiteRoot()
	if err != nil {
		t.Fatalf("Received unexpected error %v", err)
	}
	expected := "https://www.gov.uk"
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestWebsiteRootMissing(t *testing.T) {
	os.Clearenv()

	_, err := WebsiteRoot()
	if err == nil {
		t.Fatal("Expected error, received none")
	}
	errMissing, ok := err.(*EnvVarMissing)
	if !ok {
		t.Fatalf("Expected error to be a *EnvVarMissing, got %T", err)
	}
	if errMissing.EnvVar != "GOVUK_WEBSITE_ROOT" {
		t.Errorf("Expected error relating to GOVUK_WEBSITE_ROOT, got %s", errMissing.EnvVar)
	}
}

func TestAssetRoot(t *testing.T) {
	os.Clearenv()
	os.Setenv("GOVUK_ASSET_ROOT", "https://www.gov.uk")

	actual, err := AssetRoot()
	if err != nil {
		t.Fatalf("Received unexpected error %v", err)
	}
	expected := "https://www.gov.uk"
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestAssetRootMissing(t *testing.T) {
	os.Clearenv()

	_, err := AssetRoot()
	if err == nil {
		t.Fatal("Expected error, received none")
	}
	errMissing, ok := err.(*EnvVarMissing)
	if !ok {
		t.Fatalf("Expected error to be a *EnvVarMissing, got %T", err)
	}
	if errMissing.EnvVar != "GOVUK_ASSET_ROOT" {
		t.Errorf("Expected error relating to GOVUK_ASSET_ROOT, got %s", errMissing.EnvVar)
	}
}
