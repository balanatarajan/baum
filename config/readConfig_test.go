package config

import "testing"

func TestConfig(t *testing.T) {

	m, err := ReadConfig("no.cdb")
	if m != nil {
		t.Errorf("Opened a non-existing CDB")
	}

	m, err = ReadConfig("unit.cdb")

	defer m.Close()

	if m == nil {
		t.Fatal("Unable to open CDB", err)
	}

	str, err := m.Find("ip")

	if err != nil || str != "1.2.3.4" {
		t.Fatal("Unable to find key ip", err)
	}

	str, err = m.Find("ipv6")

	if err != nil || str != "xvert:xhoriz:1234" {
		t.Fatal("Unable to find key ip", err)
	}

}
