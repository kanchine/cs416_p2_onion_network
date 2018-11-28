package server_test

import (
	"../DataServer"
	"testing"
)

func TestServerInit(t *testing.T) {
	server, err := DataServer.Initialize("test.json")

	if err != nil {
		t.Errorf("Server initializetion failed")
	}

	key1 := "a"
	valExpected := "test1"

	if val, ok := server.DataBase[key1]; !ok {
		t.Errorf("Server database key %s not found", key1)
	} else {
		if val != valExpected {
			t.Errorf("Server database value for key %s actual: %s  Expected: %s", key1, val, valExpected)
		}
	}

	key2 := "b"
	valExpected = "test2"

	if val, ok := server.DataBase[key2]; !ok {
		t.Errorf("Server database key %s not found", key2)
	} else {
		if val != valExpected {
			t.Errorf("Server database value for key %s actual: %s  Expected: %s", key1, val, valExpected)
		}
	}

	key3 := "c"
	valExpected = "test3"

	if val, ok := server.DataBase[key3]; !ok {
		t.Errorf("Server database key %s not found", key3)
	} else {
		if val != valExpected {
			t.Errorf("Server database value for key %s actual: %s  Expected: %s", key1, val, valExpected)
		}
	}

	expectedIpPort := "127.0.0.1:8080"

	if server.IpPort != expectedIpPort {
		t.Errorf("Server ip port mismatch.")
	}
}
