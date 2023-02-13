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

## Examples

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
found match! mac = '3cecef123456'
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

Encode a SFT-DCMS-SINGLE key with all attributes specified (attributes that are omitted
will be left at their default values):

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
found match! mac = '3cecef123456'
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
- Property (purpose unknown; no samples that have been found had a value set for this)
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

#### SFT-SCM-LIC

Node product key, first seen in a Supermicro Server Management Utilities
brochure dated July 2013. The key was missing from the same brochure dated
April 2015.

This key was probably required to use SCM (Supermicro Command Manager).

#### SFT-DCMS-SINGLE

Node product key, available since at least 2013.

This key generally allows the use of all server management utilities
("Data Center Management Suite") with a system.

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

This key format is not currently supported by this utility.

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

## Glossary

**Node Product Key**: A product key that is activated on a specific system. The key is bound to the MAC address of
the BMC LAN port.

## Contributing

Much of the information here has been compiled from various documentation and may not have
been personally verified. Please report any inaccuracies by opening an issue.

New information would be greatly appreciated! Please open an issue or a discussion as appropriate.
