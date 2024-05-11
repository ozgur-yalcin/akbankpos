[![license](https://img.shields.io/:license-mit-blue.svg)](https://github.com/ozgur-yalcin/akbankpos.go/blob/master/LICENSE.md)
[![documentation](https://pkg.go.dev/badge/github.com/ozgur-yalcin/akbankpos.go)](https://pkg.go.dev/github.com/ozgur-yalcin/akbankpos.go/src)

# Akbankpos.go
Akbank Virtual POS API with golang


# Installation
```bash
go get github.com/ozgur-yalcin/akbankpos.go
```

# Satış
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	akbankpos "github.com/ozgur-yalcin/akbankpos.go/src"
)

// Pos bilgileri
const (
	envmode    = "TEST"                                                                                                                             // Çalışma ortamı (Production : "PROD" - Test : "TEST")
	merchantid = "2023090417500272654BD9A49CF07574"                                                                                                 // Mağaza numarası
	terminalid = "2023090417500284633D137A249DBBEB"                                                                                                 // Terminal numarası
	secretkey  = "3230323330393034313735303032363031353172675f357637355f3273387373745f7233725f73323333383737335f323272383774767276327672323531355f" // Mağaza anahtarı
)

func main() {
	api, req := akbankpos.Api(merchantid, terminalid, secretkey)
	api.SetMode(envmode)
	req.SetCardNumber("4355093000315232")   // Kart numarası (zorunlu)
	req.SetCardExpiry("11", "35")           // Son kullanma tarihi - AA,YY (zorunlu)
	req.SetCardCode("665")                  // Kart arkasındaki 3 haneli numara (zorunlu)
	req.SetAmount(1.00, "TRY")              // Satış tutarı ve para birimi (zorunlu)
	req.SetInstallment(1)                   // Taksit sayısı (zorunlu)
	req.SetCustomerIPv4("192.168.1.1")      // Müşteri IPv4 adresi (zorunlu)
	req.SetCustomerEmail("test@akbank.com") // Müşteri e-posta adresi (zorunlu)
	ctx := context.Background()
	if res, err := api.Auth(ctx, req); err == nil {
		pretty, _ := json.MarshalIndent(res, " ", " ")
		fmt.Println(string(pretty))
	} else {
		fmt.Println(err)
	}
}
```