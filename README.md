# Supermicro Product Key Utility

This is a utility for encoding and decoding product keys that are used on Supermicro BMCs.

Inspired by Peter Kleissner's [article](https://peterkleissner.com/2018/05/27/reverse-engineering-supermicro-ipmi/)
about reverse engineering the original SFT-OOB-LIC key format.

## Installation

This utility is a command line program. If Go is installed on your system, the latest version of this utility
can be built and installed by running:

```shell
go install github.com/zsrv/supermicro-product-key@latest
```

Otherwise, a binary can be downloaded from the [releases page](https://github.com/zsrv/supermicro-product-key/releases).

## Usage

Execute the binary without any arguments to see usage instructions.
Reading this entire README is recommended.

Key activation instructions can be found at the
[Supermicro website](https://store.supermicro.com/software/software-license-key-activation-usage)
([PDF version](https://store.supermicro.com/media/wysiwyg/productspecs/Supermicro_Software_License_Key_Activation_User_Guide.pdf)).

## Quick Start

All key types are bound to the MAC address of the BMC LAN port. Replace `3cecef123456` in the following examples
with the BMC MAC address of the system you plan to activate the key on.

### OOB Keys

Encode a new key:

```
$ ./supermicro-product-key oob encode 3cecef123456
CE27-F641-9B04-6B24-5D04-5D32
```

Find the MAC address associated with a key:

```
$ ./supermicro-product-key oob bruteforce CE27-F641-9B04-6B24-5D04-5D32
searching for mac address ...
found match! mac = '3CECEF123456'
```

### Non-JSON Keys

List all key types that are available:

```
$ ./supermicro-product-key nonjson listswid
License SKU         Display Name      ID
-----------         ------------      --
                    Reserved          0
                    SSM               1
                    SD5               2
SFT-SUM-LIC         SUM               3
SFT-SPM-LIC         SPM               4
SFT-SCM-LIC         SCM               5
SFT-DCMS-SINGLE     ALL               6
SFT-DCMS-SITE       SITE              7
SFT-DCMS-CALL-HOME  DCMS-CALL-HOME    8
SFT-DCMS-SVC-KEY    SFT-DCMS-SVC-KEY  9
SFT-SDDC-SINGLE     SFT-SDDC-SINGLE   210
```

Encode a key with no expiration date, and safe default values for the remaining attributes (recommended).
Replace SFT-DCMS-SINGLE with any other SKU from the list above to encode a different key type:

```
$ ./supermicro-product-key nonjson encode --sku SFT-DCMS-SINGLE 3cecef123456
AAYAAAAAAAAAAAAAAAAAAExLCU/N0RxxvG7ZACnE9iz3GCd/is0BF+s/dgeUoIyPgnp3qBgf5iyrcpByGPE9xgkT38mRmt2/R+3S/iXb/8ram5O/cXUJxvkmqZi0ODYUbze7+NLgSZ/YPPM77OBGXTgiLGYq4pihruPxfYLkT64U1vfZQLCJWBwoinbMbdxTlBC9we54hWXSF5vcY9MGNiJsZ+d3vKdoSMgCqSqHNCRLRYlYVjrT4CYmtw2k6qcBvM3Zq0pLwjznsP3tD7CxHKIf5DoIph6zj9M3VTtiPm2yzgJOSJwNGhFJq4bvErw2PqqGURub+lNz5AazeafiRw==
```

For **experimentation purposes**: Encode a key with all attributes specified
(attributes that are omitted will be left at their default values):

```
$ ./supermicro-product-key nonjson encode --sku SFT-DCMS-SINGLE --software-version ABC123 --invoice-number 0123456789 --creation-date 2020-12-30T12:00:00Z --expiration-date 1970-01-01T00:00:00Z --property 01AA02FF 3cecef123456
AAYAAAAAAAAAAAAAAAAAAJOVA97uSfqDCtInPd8H2g4rUdY5PtJ3op7hUYaFWOn2aWeT/f+4ZaMdelxJgFG3NjRPqXIMfJ2mFdeR8tZYfNusG3i6SP7cDBsN5Kbu/Bfj3q6WWlUKavF9j6oulmnvS83CZkzEdjQrOfq7bokHV6HnAUn/83UGKn+b5Db4E7fcGjmHYOCzP9rTTZnmtfUHShcEN48+Aka9ACcyfUDo9kdiJQpFRzTQz0ay5gGii6MB0Yw+G1JECp804aNGzCzl5PsP122dZ0cTFFN7THfy07Se1NDd+GCCoH2AN7UFa8tTKquZBgdwBgU/EOPX31YXJQ==
```

Decode a key:

```
$ ./supermicro-product-key nonjson decode 3cecef123456 AAYAAAAAAAAAAAAAAAAAAJOVA97uSfqDCtInPd8H2g4rUdY5PtJ3op7hUYaFWOn2aWeT/f+4ZaMdelxJgFG3NjRPqXIMfJ2mFdeR8tZYfNusG3i6SP7cDBsN5Kbu/Bfj3q6WWlUKavF9j6oulmnvS83CZkzEdjQrOfq7bokHV6HnAUn/83UGKn+b5Db4E7fcGjmHYOCzP9rTTZnmtfUHShcEN48+Aka9ACcyfUDo9kdiJQpFRzTQz0ay5gGii6MB0Yw+G1JECp804aNGzCzl5PsP122dZ0cTFFN7THfy07Se1NDd+GCCoH2AN7UFa8tTKquZBgdwBgU/EOPX31YXJQ==
{
        "FormatVersion": 0,
        "SoftwareIdentifier": {
                "SKU": "SFT-DCMS-SINGLE",
                "DisplayName": "ALL",
                "ID": 6
        },
        "SoftwareVersion": "ABC123",
        "InvoiceNumber": "0123456789",
        "CreationDate": "2020-12-30T12:00:00Z",
        "ExpirationDate": "1970-01-01T00:00:00Z",
        "Property": "AaoC/w==",
        "SecretData": "MjdiOGUzYWFmN2NlZGU0MzNjNWUzYjUzZjU5YzJhOWI=",
        "Checksum": 24
}
```

Find the MAC address associated with a key:

```
$ ./supermicro-product-key nonjson bruteforce AAYAAAAAAAAAAAAAAAAAAJOVA97uSfqDCtInPd8H2g4rUdY5PtJ3op7hUYaFWOn2aWeT/f+4ZaMdelxJgFG3NjRPqXIMfJ2mFdeR8tZYfNusG3i6SP7cDBsN5Kbu/Bfj3q6WWlUKavF9j6oulmnvS83CZkzEdjQrOfq7bokHV6HnAUn/83UGKn+b5Db4E7fcGjmHYOCzP9rTTZnmtfUHShcEN48+Aka9ACcyfUDo9kdiJQpFRzTQz0ay5gGii6MB0Yw+G1JECp804aNGzCzl5PsP122dZ0cTFFN7THfy07Se1NDd+GCCoH2AN7UFa8tTKquZBgdwBgU/EOPX31YXJQ==
searching for mac address ...
found match! mac = '3CECEF123456'
```

### JSON Keys

Verify the signature of a product key:

```
$ ./supermicro-product-key json verify ac1f6b3ddaec '{"ProductKey":{"Node":{"LicenseID":"1","LicenseName":"SFT-OOB-LIC","CreateDate":"20200921"},"Signature":"OAaLKLy5IEK9WnIdnyA9ew89qTKQrm1eu+Q84CbwjR7XG7JGYccec+3vS3y/kQRRej3DcNVQPWsasX86ROTT+LZFsNY2mIEbQ6+Y/Tmv6+jwYgbQjEN6CjI7ahyKcebN12+3cLvPZyRf3kDqgtcpfuw3Qeg8BbhhyHQk29yNp+NG0XbKn02sHTrskvAGgG0GGlDCT5YmNa0gDSMzsvt/eH9nskb5opQNE3j7MAMXbjpI7xVHRbmB2N5iSu8gQUj0/pmk615ztM/uB54ur3GninJRU74S9Kotz+JunJg4pprGyQW544ggmzklmtr3zCA3GK/d929eZsVk5p8UxXG7wQ=="}}'
signature verified ok
```

Find the MAC address associated with a key:

```
$ ./supermicro-product-key json bruteforce '{"ProductKey":{"Node":{"LicenseID":"1","LicenseName":"SFT-OOB-LIC","CreateDate":"20200921"},"Signature":"OAaLKLy5IEK9WnIdnyA9ew89qTKQrm1eu+Q84CbwjR7XG7JGYccec+3vS3y/kQRRej3DcNVQPWsasX86ROTT+LZFsNY2mIEbQ6+Y/Tmv6+jwYgbQjEN6CjI7ahyKcebN12+3cLvPZyRf3kDqgtcpfuw3Qeg8BbhhyHQk29yNp+NG0XbKn02sHTrskvAGgG0GGlDCT5YmNa0gDSMzsvt/eH9nskb5opQNE3j7MAMXbjpI7xVHRbmB2N5iSu8gQUj0/pmk615ztM/uB54ur3GninJRU74S9Kotz+JunJg4pprGyQW544ggmzklmtr3zCA3GK/d929eZsVk5p8UxXG7wQ=="}}'
searching for mac address ...
found match! mac = 'AC1F6B3DDAEC'
```

List all key types that are available:

```
$ ./supermicro-product-key json listswid
License SKU       ID
-----------       --
SFT-OOB-LIC       1
SFT-DCMS-SINGLE   2
SFT-SPM-LIC       3
SFT-DCMS-SVC-KEY  4
SFT-SDDC-SINGLE   5
```

## Product Key Formats and SKUs

This section describes each product key format, and each license SKU
that is supported by each format.

### OOB

A 24-character hex string split into groups of 4 characters, separated by dashes.

Keys in this format are created using the MAC address of the BMC that the key will be activated on
as input.

#### SFT-OOB-LIC

Used by 11th generation and earlier platforms.
This key is delivered in the JSON key format for 12th generation platforms.

Enables various out-of-band (OOB) management features.

The functionality provided by this key is included in SFT-DCMS-SINGLE.

### Non-JSON

A 344-character base64-encoded string.

Keys in this format are encrypted and decrypted using the MAC address of the BMC that the key
will be activated on as input.

Attributes encoded in the key:
- Key format version byte (always 0)
- Software ID (in the form of a numeric ID byte and a corresponding display name string that is usually
  different from the license SKU)
- Software version string (always "none" in samples from as far back as 2016; a sample from 2014 had
  a null value instead)
- Invoice number string (always "none" in samples from as far back as 2016; a sample from 2014 had 
  a null value instead)
- Creation date (stored as a Unix timestamp with second precision, converted to four bytes)
- Expiration date (stored as a Unix timestamp with second precision, converted to four bytes;
  a value of 0 is treated as "no expiration date")
- Property (purpose unknown; no samples that have been found had a value set for this). **Warning**: Keys with a
  non-null value set for this attribute can be activated on a system but have been shown to be invalid. Only set
  this attribute for experimentation purposes, with a method to remove the key from the system readily available.
- Secret data (calculated using several other attributes as input; keys are validated by the BMC by 
  calculating this and comparing the result to the value stored in the key)
- Checksum (computed using the other attributes as input)

#### SSM (Unknown SKU)

Available since at least 2015.

Purpose unknown. Probably related to Supermicro Server Manager (SSM).

#### SD5 (Unknown SKU)

Available since at least 2015.

Purpose unknown. Probably related to Supermicro SuperDoctor 5 (SD5).

#### SFT-SUM-LIC

Node product key, available since 2015.
Removed from the SUM user guide in 2020.

This key may have been required to use Supermicro Update Manager (SUM).

#### SFT-SPM-LIC

Node product key, available since at least 2013.

This key is required to use Supermicro Power Manager (SPM).

The functionality provided by this key is also included in SFT-DCMS-SINGLE.

##### Samples

```
Source: SPM User Guide (PDF included in SPM download)

MAC Address: 002590FE7A43

Key: AAQAAAAAAAAAAAAAAAAAAGgGV+3RQ/KGDk0jKe/2bnjvnc89ke8Y/Bd1gCtXKH4CjacoAEmPEzff1e5igzxl2RIQns2IlSNCEHk8hi6zEdLRefQgBVGMoyRwHywF/Qid4ugrbX7Q+xE+I4ly0/Hu7QilyJusitzD5SXl6RzkHhCTbrFMOrWP0J9TxeRjpMZGPe2+TVPhuP2Bs2pPMnL5wYFwsoiOCr3fUu9ef7c7DQ13v/00cOAukoP1zxxtbFbwDgMG8ZuMqmTH2d7PCz5g7DXdncu9LyVWxG69jwPiOHXhJVJUfmaw5AhV6FDlXGF+8p/kMvNkzsMddcv+IXfb4Q==

Decoded:
{
        "FormatVersion": 0,
        "SoftwareIdentifier": {
                "SKU": "SFT-SPM-LIC",
                "DisplayName": "SPM",
                "ID": 4
        },
        "SoftwareVersion": "",
        "InvoiceNumber": "",
        "CreationDate": "2014-09-26T08:45:11Z",
        "ExpirationDate": "1970-01-01T00:00:00Z",
        "Property": null,
        "SecretData": null,
        "Checksum": 47
}
```

#### SFT-SCM-LIC

Node product key, first seen in a Supermicro Server Management Utilities
brochure dated July 2013. The key was missing from the same brochure dated
April 2015.

This key was probably required to use Supermicro Command Manager (SCM).

#### SFT-DCMS-SINGLE

Node product key, available since at least 2013.

This key generally allows the use of all server management utilities
("Data Center Management Suite") with a system.

##### Samples

```
Source: https://store.supermicro.com/media/wysiwyg/productspecs/Supermicro_Software_License_Key_Activation_User_Guide.pdf

MAC Address: 0CC47A87AEAA

Key: AAYAAAAAAAAAAAAAAAAAAGSw1zjcokBJrE1va2uIhi1umjMFAUZVqBuP06oZmZkT4y1pOvSYXDuMbJ8S27SyrA8A3S7XgblZeEYbrUF+JzAg6SSAgN6TxAwc8QZS5rIqi7oDARAjJkxMJJKHwYcHiCW1/pSeclAMquDl/g0mTWNxYeG01mGsrtclIKIm3tvN8i6/zpEeV6jfE7AUOAAunFOVWt/OMcD+iTiY7pN3HhQ39cS3t7dBYoWv/a1sPjtlCelA0iLjoXi8TCbqVCxpVv4QQRI7Qv2wH4hNcl5dwndiMkamCYDwU4P7DlN/v+++bAgQ1N0hE/zwLqnxchD+iw==

Decoded:
{
        "FormatVersion": 0,
        "SoftwareIdentifier": {
                "SKU": "SFT-DCMS-SINGLE",
                "DisplayName": "ALL",
                "ID": 6
        },
        "SoftwareVersion": "none",
        "InvoiceNumber": "none",
        "CreationDate": "2021-12-22T05:55:53Z",
        "ExpirationDate": "1970-01-01T00:00:00Z",
        "Property": null,
        "SecretData": "OWJmNjNmYjg2OGYzNDkwMzY2ODRmYjE1NDdmMWJjN2I=",
        "Checksum": 213
}
```

#### SFT-DCMS-SITE

Described as a site license for individual software modules in
a Supermicro Server Management Utilities brochure dated July 2013.
The license was missing from the same brochure dated April 2015.

No other information has been found.

#### SFT-DCMS-CALL-HOME

Node product key introduced in 2017.

This SKU appears to have been superseded by SFT-DCMS-SVC-KEY in 2018.

#### SFT-DCMS-SVC-KEY

Node product key introduced in 2017.

This key (in addition to SFT-DCMS-SINGLE) is required to use the
Service Calls (also known as Call Home) feature available in SSM and SUM.

##### Samples

```
Source: SSM User Guide (PDF included in SSM download)

MAC Address: 0CC47AD57D8D

Key: AAkAAAAAAAAAAAAAAAAAAMjnO7OIeNNWpc63TFto8dp6A5UrXzkBpQdkhtnMrUR/oTFKIdhLPpIi6b32lQJFaoPly7uj2OztgzUxjKy1kdMDrEEFra1KlLDrBoZC88fAWfuVXmnVBhjR7tNKSa4r29owr8M3ETun+GxqerDT8kDa+jafMEkETjDJ2Gln6sk7oRCLA7xVZhG1RfkyjcrO+qyYL4OOHH8GG8CUTDx/dlBCXH8i3TL3g5d7X8U/B2XO/z85JUWOeVgwEzUXxK0eN5l3ub/OGYXVzMAH0fiq0LU6srDV+Qvc82gwckcrUKGpi0c6DUXl/qWUWDsWFrG48w==

Decoded:
{
        "FormatVersion": 0,
        "SoftwareIdentifier": {
                "SKU": "SFT-DCMS-SVC-KEY",
                "DisplayName": "SFT-DCMS-SVC-KEY",
                "ID": 9
        },
        "SoftwareVersion": "none",
        "InvoiceNumber": "none",
        "CreationDate": "2018-12-18T10:09:55Z",
        "ExpirationDate": "2019-04-17T09:09:55Z",
        "Property": null,
        "SecretData": "MTY3YmIyMTg2NGIyOTMxYzkxNGUwY2FkMmY3ZDE1MWE=",
        "Checksum": 147
}
```

#### SFT-SDDC-SINGLE

Node product key introduced in 2021.

This key (in addition to SFT-DCMS-SINGLE) is required to manage the system
with Supermicro SuperCloud Composer (SCC), a "composable cloud management
platform that provides a unified dashboard to administer software-defined
data centers" (hence "SDDC").

This key is also required to use the Attestation command's Compare action
that was introduced in SUM between 2020 and 2022.

### JSON

A variable-length JSON string.

The contents are digitally signed, and the signatures verified using a public key that is
embedded in the BMC firmware.

This key format can only be verified by this utility, not generated.

#### SFT-OOB-LIC

##### Samples

```
Source: https://store.supermicro.com/software/dcms-key-activation-guide

MAC Address: AC1F6B3DDAEC (different from what is displayed in the guide!)

Key: {"ProductKey":{"Node":{"LicenseID":"1","LicenseName":"SFT-OOB-LIC","CreateDate":"20200921"},"Signature":"OAaLKLy5IEK9WnIdnyA9ew89qTKQrm1eu+Q84CbwjR7XG7JGYccec+3vS3y/kQRRej3DcNVQPWsasX86ROTT+LZFsNY2mIEbQ6+Y/Tmv6+jwYgbQjEN6CjI7ahyKcebN12+3cLvPZyRf3kDqgtcpfuw3Qeg8BbhhyHQk29yNp+NG0XbKn02sHTrskvAGgG0GGlDCT5YmNa0gDSMzsvt/eH9nskb5opQNE3j7MAMXbjpI7xVHRbmB2N5iSu8gQUj0/pmk615ztM/uB54ur3GninJRU74S9Kotz+JunJg4pprGyQW544ggmzklmtr3zCA3GK/d929eZsVk5p8UxXG7wQ=="}}
```

#### SFT-DCMS-SINGLE

##### Samples

```
Source: https://github.com/supermicro/redfish/blob/ea5ea99eda0b7ac50d48ca7faa307dfcd3f41e05/Postman_Collections/05_managers.postman_collection.json#L1348

MAC Address: 3CECEF72FC46

Key: {"ProductKey":{"Node":{"LicenseID":"2","LicenseName":"SFT-DCMS-SINGLE","CreateDate":"20220614"},"Signature":"ZKFCkgKEYh9+8MNZW7RfPlt/nRxQJGJ0kLHLkalLt1tpgs4MTLHrXvp/eZzfhSPUb5qMNu9RkFn9MaukK6vNXlOIG7ijbR+vjkxVcdIIkMnhzHFLxE/0ws74/lJyGLkSO1jHRQRaczSDuHgzSgsWivjHejB/tRlSpnAEM7FplgyuBSbisek8pEgSKua5jCf7Zn4sjYXXO7T9rTV4aFq090XgRbEay45eBSGpun9pcyGs8UIeNH93qzqCmlkcjj+bFSNcm3VeucEjScE3fzqG93NMEQQWYEdsYcuJb4a+kWP/ffFvyVRWvqSWvPgD5N+eNqKAmmC4MmjykRy3DWw4fA=="}}
```

## Platform Support

| Platform Generation | OOB Key | Non-JSON Key | JSON Key |
|---------------------|:-------:|:------------:|:--------:|
| 8 and earlier       |   NO    |      NO      |    NO    |
| 9                   |   YES   |      NO      |    NO    |
| 10                  |   YES   |     YES      |    NO    |
| 11                  |   YES   |     YES      |    NO    |
| 12 (select models)  |   YES   |     YES      |    NO    |
| 12                  |   NO    |      NO      |   YES    |

Select 12th generation platform motherboards accept non-JSON keys instead of JSON keys
([source](https://store.supermicro.com/media/wysiwyg/productspecs/Supermicro_Software_License_Key_Activation_User_Guide.pdf)):
- H12DSU-iN
- H12DST-B
- H12SST-PS
- H12SSW-iN
- H12SSW-iNL
- H12SSW-NT
- H12SSW-NTL

## Troubleshooting

### How to remove a non-JSON format key

If you need to remove a non-JSON format key that has been activated on a system, you can use the Supermicro SUM
utility to do this, up to and including version 2.4.0 (dated 2019-12-06). SUM versions after 2.4.0 no longer have
the ability to remove keys.

You may need to do this if you activate a key with a value set for the "property" parameter
(see [issue #2](https://github.com/zsrv/supermicro-product-key/issues/2));
the BMC will accept the key, but it has been reported that BMC functions that require the key you activated will
still report that the system is not licensed.

SUM 2.4.0 is no longer available on Supermicro's website, but it can be found hosted on other websites.

SHA256 checksum for verification ([source](https://aur.archlinux.org/cgit/aur.git/commit/?h=supermicro-update-manager&id=a1a30db1c126a476c0d92339256259ef94adac67)):

```
sum_2.4.0_Linux_x86_64_20191206.tar.gz d0d7203913334e02d4ea1c8493834c5a0e7236040380447e22697b816dc614c7
```

Example of how to remove a key:

```
> ./sum -i 1.2.3.4 -u user -p password -c QueryProductKey
Supermicro Update Manager (for UEFI BIOS) 2.4.0 (2019/12/06) (x86_64)
Copyright(C) 2013-2019 Super Micro Computer, Inc. All rights reserved.

[1] SFT-DCMS-Single, version: ABC123, invoice: 0123456789, creation date: 2020/12/30 12:00:00(Key is good.)
Number of product keys: 1

> ./sum -i 1.2.3.4 -u user -p password -c ClearProductKey --key_index 1
Supermicro Update Manager (for UEFI BIOS) 2.4.0 (2019/12/06) (x86_64)
Copyright(C) 2013-2019 Super Micro Computer, Inc. All rights reserved.

Node product key is deactivated for 1.2.3.4.
```

## Glossary

**Node Product Key**: A product key that is activated on a specific system. The key is bound to the MAC address of
the BMC LAN port.

## Contributing

Much of the information here has been compiled from various documentation and may not have
been personally verified. Please report any inaccuracies by opening an issue.

New information would be greatly appreciated! Please open an issue or a discussion as appropriate.
