package auth

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func setupSecret(secretName string) (string, error) {
	//find temp location
	dir := os.TempDir()
	//make secrets folder there
	secFolder := filepath.Join(dir, "derekTestSecrets")

	if _, err := os.Stat(secFolder); os.IsNotExist(err) {
		os.Mkdir(secFolder, 0777)
	}
	//write the file
	secFile := filepath.Join(secFolder, secretName)

	secretContents := fmt.Sprintf("Secret:%s", secretName)
	fileErr := ioutil.WriteFile(secFile, []byte(secretContents), 0777)
	if fileErr != nil {
		return "", fileErr
	}
	return secFolder, nil
}

func Test_getSecretKey_notFound(t *testing.T) {

	secretName := "derek-secret"

	secDir, err := setupSecret(secretName)

	if err != nil {
		t.Errorf("secret setup failed: %s", err)
		t.Fail()
	}

	defer os.RemoveAll(secDir)

	_, err = getSecret(secretName)

	if err == nil {
		t.Errorf("getSecret should error as secret doesn't exist")
		t.Fail()
	}

}

func Test_getSecret_usingPath(t *testing.T) {

	secretName := "derek-secret"
	expected := fmt.Sprintf("Secret:%s", secretName)

	secDir, err := setupSecret(secretName)

	if err != nil {
		t.Errorf("secret setup failed: %s", err)
		t.Fail()
	}

	defer os.RemoveAll(secDir)

	secretLocation := filepath.Join(secDir, secretName)

	secretKey, err := getSecret(secretLocation)
	if err != nil {
		t.Errorf("getSecretKey failed: %s", err)
		t.Fail()
	}

	if string(secretKey) != expected {
		t.Errorf("want '%s' but got '%s'", expected, secretKey)
		t.Fail()
	}
}
