package config

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"testing"
)

func TestConfigShouldParse(t *testing.T) {
	t.Run("single urls", func(t *testing.T) {
		env := "OUTPUT_ADDRS"
		envsToPersist := []string{"OUTPUT_ADDRS", "SERVICE_TYPE"}
		ta := "http://127.0.0.1:8080/input/a"
		tu, err := url.Parse(ta)
		if err != nil {
			t.Fatal(err)
		}

		expected_urls := []url.URL{*tu}
		for _, pe := range envsToPersist {
			oe := os.Getenv(pe)
			if oe == "" {
				defer os.Unsetenv(pe)
			} else {
				defer os.Setenv(pe, oe)
			}
		}

		os.Setenv(env, ta)
		os.Setenv("SERVICE_TYPE", "not")

		c, err := New()
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(c.OutputAddrs, expected_urls) {
			t.Errorf("wanted %v, got %v", expected_urls, c.OutputAddrs)
		}
	})

	t.Run("multiple urls", func(t *testing.T) {
		env := "OUTPUT_ADDRS"
		envsToPersist := []string{"OUTPUT_ADDRS", "SERVICE_TYPE"}
		ta1 := "http://127.0.0.1:8080/input/a"
		ta2 := "http://127.0.0.1:8080/input/b"
		tu1, err := url.Parse(ta1)
		if err != nil {
			t.Fatal(err)
		}
		tu2, err := url.Parse(ta2)
		if err != nil {
			t.Fatal(err)
		}

		expected_urls := []url.URL{*tu1, *tu2}
		for _, pe := range envsToPersist {
			oe := os.Getenv(pe)
			if oe == "" {
				defer os.Unsetenv(pe)
			} else {
				defer os.Setenv(pe, oe)
			}
		}

		os.Setenv(env, fmt.Sprintf("%s,%s", ta1, ta2))
		os.Setenv("SERVICE_TYPE", "not")

		c, err := New()
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(c.OutputAddrs, expected_urls) {
			t.Errorf("wanted %v, got %v", expected_urls, c.OutputAddrs)
		}
	})
}
