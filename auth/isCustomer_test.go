// Copyright (c) Derek Author(s) 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package auth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_isCustomer_Yes(t *testing.T) {
	os.Setenv("validate_customers", "true")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `alex
richard`)
	})

	server := httptest.NewServer(handler)
	defer server.Close()
	os.Setenv("customers_url", server.URL+"/CUSTOMERS")

	owner := "alex"
	isCustomer, err := IsCustomer(owner, server.Client())
	if err != nil {
		t.Errorf("want no error, but got one: %s", err)
		t.Fail()
	}

	want := true
	if isCustomer != want {
		t.Errorf("want %s customer value %t but got %t", owner, want, isCustomer)
		t.Fail()
	}
}

func Test_isCustomer_No_EmptyFile(t *testing.T) {
	os.Setenv("validate_customers", "true")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, ``)
	})

	server := httptest.NewServer(handler)
	defer server.Close()
	os.Setenv("customers_url", server.URL+"/CUSTOMERS")

	owner := "alex"
	isCustomer, err := IsCustomer(owner, server.Client())
	if err != nil {
		t.Errorf("want no error, but got one: %s", err)
		t.Fail()
	}

	want := false
	if isCustomer != want {
		t.Errorf("want %s customer value %t but got %t", owner, want, isCustomer)
		t.Fail()
	}
}

func Test_isCustomer_No_NotInList(t *testing.T) {
	os.Setenv("validate_customers", "true")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `alex
richard`)
	})

	server := httptest.NewServer(handler)
	defer server.Close()
	os.Setenv("customers_url", server.URL+"/CUSTOMERS")

	owner := "johmmcabe"
	isCustomer, err := IsCustomer(owner, server.Client())
	if err != nil {
		t.Errorf("want no error, but got one: %s", err)
		t.Fail()
	}
	want := false
	if isCustomer != want {
		t.Errorf("want %s customer value %t but got %t", owner, want, isCustomer)
		t.Fail()
	}
}
