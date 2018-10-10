package config

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestNewConfig_NoSecretPath(t *testing.T) {

	os.Setenv("secret_path", "")

	_, err := NewConfig()

	if err == nil {
		t.Fail()
	}

	want := "secret_path env-var not set, this should be /var/openfaas/secrets or /run/secrets"
	if err.Error() != want {
		t.Errorf("want %q, got %q", want, err.Error())
		t.Fail()
	}
}

func TestNewConfig_ValidSecretPath_WithApplicationID(t *testing.T) {
	privateWant := "private"
	secretWant := "secret"
	appIDWant := "321"
	tmpDir := os.TempDir()

	ioutil.WriteFile(path.Join(tmpDir, "derek-private-key"), []byte(privateWant), 0600)
	ioutil.WriteFile(path.Join(tmpDir, "derek-secret-key"), []byte(secretWant), 0600)

	defer os.RemoveAll(path.Join(tmpDir, "derek-private-key"))
	defer os.RemoveAll(path.Join(tmpDir, "derek-secret-key"))

	os.Setenv("secret_path", tmpDir)
	os.Setenv("application_id", appIDWant)

	cfg, err := NewConfig()

	if err != nil {
		t.Errorf("%s", err.Error())
		t.Fail()
		return
	}

	if cfg.SecretKey != secretWant {
		t.Errorf("want %q, got %q", secretWant, cfg.SecretKey)
		t.Fail()
	}

	if cfg.PrivateKey != privateWant {
		t.Errorf("want %q, got %q", privateWant, cfg.PrivateKey)
		t.Fail()
	}

	if cfg.ApplicationID != appIDWant {
		t.Errorf("want %q, got %q", appIDWant, cfg.ApplicationID)
		t.Fail()
	}
}

func Test_getFirstLine(t *testing.T) {
	var exampleSecrets = []struct {
		secret       string
		expectedByte string
	}{
		{
			secret:       "New-line \n",
			expectedByte: "New-line ",
		},
		{
			secret: `Newline and text 
			`,
			expectedByte: "Newline and text ",
		},
		{
			secret:       `Example secret2 `,
			expectedByte: `Example secret2 `,
		},
		{
			secret:       "\n",
			expectedByte: "",
		},
		{
			secret:       "",
			expectedByte: "",
		},
	}
	for _, test := range exampleSecrets {

		t.Run(string(test.secret), func(t *testing.T) {
			stringNoLines := getFirstLine([]byte(test.secret))
			if test.expectedByte != string(stringNoLines) {
				t.Errorf("String after removal - wanted: \"%s\", got \"%s\"", test.expectedByte, test.secret)
			}
		})
	}
}
