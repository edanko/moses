package license

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/hyperboloide/lk"
)

type lic struct {
	Hostname []string  `json:"hostname"`
	End      time.Time `json:"end"`
}

// IsLicensed checks for license
func IsLicensed() (bool, error) {
	licenseB32, err := ioutil.ReadFile("key.txt")
	if err != nil {
		return false, err
	}

	// the public key b32 encoded from the private key using: lkgen pub my_private_key_file`.
	// It should be hardcoded somewhere in your app.
	const publicKeyBase32 = "***REMOVED***"

	// Unmarshal the public key.
	publicKey, err := lk.PublicKeyFromB32String(publicKeyBase32)
	if err != nil {
		return false, err
	}

	// Unmarshal the customer license.
	license, err := lk.LicenseFromB32String(string(licenseB32))
	if err != nil {
		return false, err
	}

	// validate the license signature.
	if ok, err := license.Verify(publicKey); err != nil {
		return false, err
	} else if !ok {
		return false, errors.New("Invalid license signature")
	}

	result := lic{}

	if err := json.Unmarshal(license.Data, &result); err != nil {
		return false, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	var flag bool

	for _, host := range result.Hostname {
		if host == hostname {
			flag = true
		}
	}

	// Now you just have to check that the end date is after time.Now() then you can continue!
	if result.End.Before(time.Now()) {
		return false, fmt.Errorf("License expired on: %s", result.End.Format("2006-01-02"))
	} else if !flag {
		return false, errors.New("Licensed not for this computer")
	} else {
		return true, nil
	}
}
