# L(ean)C(loud)

[![GoDoc](https://godoc.org/github.com/bcho/lc?status.svg)](https://godoc.org/github.com/bcho/lc)
[![Build Status](https://travis-ci.org/bcho/lc.svg)](https://travis-ci.org/bcho/lc)

Customized LeanCloud (internal) API client.

**This package shares no relationship with LeanCloud**

## Usage

### Client

```go
clientAuth, err := client.NewClientAuthFromLogin("email", "password")
checkErr(err)

client := client.NewClient(clientAuth)
```

### Get recent SMS records

```go
appId := "leancloud_app_id"
getOptions := &sms.GetRecentRecordsOptions{Limit: 10}
total, records, _, err := sms.GetRecentRecords(client, appId, getOptions)
checkErr(err)

fmt.Printf("total: %d\n", total)
for _, record := range records {
        fmt.Printf("%s %s %s %s\n", *record.Phone, *record.Status, *record.Message, *record.Created)
}
```

## LICENSE

[MIT](https://hbc.mit-license.org)
