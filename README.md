[![CircleCI](https://circleci.com/gh/Financial-Times/image-resolver/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/image-resolver/tree/master)
[![Coverage Status](https://coveralls.io/repos/github/Financial-Times/image-resolver/badge.svg)](https://coveralls.io/github/Financial-Times/image-resolver)
# image-resolver

Image resolver is an internally used API for retrieving unrolled images and leadimages. It receive a content and return the content plus images unrolled.

## Usage
### Install
`go get -u github.com/Financial-Times/image-resolver`

## Running locally
To run the service locally, you will need to run the following commands first to get the vendored dependencies for this project:
  `go get github.com/kardianos/govendor` and
  `govendor sync`

```
Usage: image-resolver [OPTIONS]

Options:
  --port="8080"                               Port to listen on ($PORT)
  --cprHost="content-public-read"             The host to connect to content-public-read API
  --routerAddress="localhost:8080"            Vulcan host
  --graphiteTCPAddress=""                     Graphite TCP address, e.g. graphite.ft.com:2003. Leave as default if you d
o NOT want to output to graphite (e.g. if running locally) ($GRAPHITE_ADDRESS)
  --graphitePrefix=""                         Prefix to use. Should start with content, include the environment, and the
 host name. e.g. coco.pre-prod.image-resolver.1 ($GRAPHITE_PREFIX)
  --logMetrics=false                          Whether to log metrics. Set to true if running locally and you want metric
s output ($LOG_METRICS)

```

## Endpoints

### Application specific endpoints:

* /content
* /internalcontent

### Admin specific endpoints:

* /ping
* /build-info
* /__ping
* /__build-info
* /__health
* /__gtg


## Example 1 (main image)
POST: /content
Body:
```
{
  "id": "http://www.ft.com/thing/22c0d426-1466-11e7-b0c1-37e417ee6c76",
  "type": "http://www.ft.com/ontology/content/Article",
  "bodyXML": "<body></body>",
  "title": "Brexit begins as Theresa May triggers Article 50",
  "alternativeTitles": {
    "promotionalTitle": "Brexit begins as Theresa May triggers Article 50"
  },
  "standfirst": "Prime minister sets out Britain’s negotiating stance in statement to MPs",
  "alternativeStandfirsts": {},
  "byline": "George Parker and Kate Allen in London and Arthur Beesley in Brussels",
  "firstPublishedDate": "2017-03-29T11:07:52.000Z",
  "publishedDate": "2017-03-30T06:54:02.000Z",
  "identifiers": [
    {
      "authority": "http://api.ft.com/system/FTCOM-METHODE",
      "identifierValue": "22c0d426-1466-11e7-b0c1-37e417ee6c76"
    }
  ],
  "requestUrl": "http://test.api.ft.com/content/22c0d426-1466-11e7-b0c1-37e417ee6c76",
  "brands": [
    "http://api.ft.com/things/dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54"
  ],
  "mainImage": {
    "id": "http://test.api.ft.com/content/639cd952-149f-11e7-2ea7-a07ecd9ac73f"
  },
  "alternativeImages": {},
  "comments": {
    "enabled": true
  },
  "standout": {
    "editorsChoice": false,
    "exclusive": false,
    "scoop": false
  },
  "publishReference": "tid_ra4srof3qc",
  "lastModified": "2017-03-31T15:42:35.266Z",
  "canBeSyndicated": "yes",
  "accessLevel": "subscribed",
  "canBeDistributed": "yes"
}
```

Response:
```
{
  "accessLevel": "subscribed",
  "alternativeImages": {},
  "alternativeStandfirsts": {},
  "alternativeTitles": {
    "promotionalTitle": "Brexit begins as Theresa May triggers Article 50"
  },
  "bodyXML": "<body></body>",
  "brands": [
    "http://api.ft.com/things/dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54"
  ],
  "byline": "George Parker and Kate Allen in London and Arthur Beesley in Brussels",
  "canBeDistributed": "yes",
  "canBeSyndicated": "yes",
  "comments": {
    "enabled": true
  },
  "firstPublishedDate": "2017-03-29T11:07:52.000Z",
  "id": "http://www.ft.com/thing/22c0d426-1466-11e7-b0c1-37e417ee6c76",
  "identifiers": [
    {
      "authority": "http://api.ft.com/system/FTCOM-METHODE",
      "identifierValue": "22c0d426-1466-11e7-b0c1-37e417ee6c76"
    }
  ],
  "lastModified": "2017-03-31T15:42:35.266Z",
  "mainImage": {
    "alternativeImages": {},
    "alternativeStandfirsts": {},
    "alternativeTitles": {},
    "canBeDistributed": "verify",
    "description": "Donald Tusk, president of the European Union (EU), holds the letter invoking Article 50 of the Lisbon Treaty from U.K. Prime Minister Theresa May as leaves following a news conference at the European Council in Brussels, Belgium, on Wednesday, March 29, 2017. The U.K. will&nbsp;start the clock&nbsp;on two years of negotiations to withdraw from the European Union on Wednesday, when Britain's ambassador hands EU President Donald Tusk&nbsp;a hand-signed&nbsp;letter&nbsp;from Prime Minister&nbsp;Theresa May&nbsp;invoking Article 50 of the Lisbon Treaty, the legal exit mechanism. Photographer: Jasper Juinen/Bloomberg",
    "firstPublishedDate": "2017-03-29T19:39:00.000Z",
    "id": "http://www.ft.com/thing/639cd952-149f-11e7-2ea7-a07ecd9ac73f",
    "identifiers": [
      {
        "authority": "http://api.ft.com/system/FTCOM-METHODE",
        "identifierValue": "639cd952-149f-11e7-2ea7-a07ecd9ac73f"
      }
    ],
    "lastModified": "2017-03-29T19:39:31.361Z",
    "members": [
      {
        "alternativeImages": {},
        "alternativeStandfirsts": {},
        "alternativeTitles": {},
        "binaryUrl": "http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/639cd952-149f-11e7-b0c1-37e417ee6c76",
        "canBeDistributed": "verify",
        "copyright": {
          "notice": "© Bloomberg"
        },
        "description": "Donald Tusk, president of the European Union (EU), holds the letter invoking Article 50 of the Lisbon Treaty from U.K. Prime Minister Theresa May as leaves following a news conference at the European Council in Brussels, Belgium, on Wednesday, March 29, 2017. The U.K. will&nbsp;start the clock&nbsp;on two years of negotiations to withdraw from the European Union on Wednesday, when Britain's ambassador hands EU President Donald Tusk&nbsp;a hand-signed&nbsp;letter&nbsp;from Prime Minister&nbsp;Theresa May&nbsp;invoking Article 50 of the Lisbon Treaty, the legal exit mechanism. Photographer: Jasper Juinen/Bloomberg",
        "firstPublishedDate": "2017-03-29T19:39:00.000Z",
        "id": "http://www.ft.com/thing/639cd952-149f-11e7-b0c1-37e417ee6c76",
        "identifiers": [
          {
            "authority": "http://api.ft.com/system/FTCOM-METHODE",
            "identifierValue": "639cd952-149f-11e7-b0c1-37e417ee6c76"
          }
        ],
        "lastModified": "2017-03-29T19:39:31.361Z",
        "pixelHeight": 1152,
        "pixelWidth": 2048,
        "publishReference": "tid_5ypvntzcpu",
        "publishedDate": "2017-03-29T19:39:00.000Z",
        "requestUrl": "http://test.api.ft.com/content/639cd952-149f-11e7-b0c1-37e417ee6c76",
        "title": "",
        "type": "http://www.ft.com/ontology/content/MediaResource"
      }
    ],
    "publishReference": "tid_5ypvntzcpu",
    "publishedDate": "2017-03-29T19:39:00.000Z",
    "requestUrl": "http://test.api.ft.com/content/639cd952-149f-11e7-2ea7-a07ecd9ac73f",
    "title": "",
    "type": "http://www.ft.com/ontology/content/ImageSet"
  },
  "publishReference": "tid_ra4srof3qc",
  "publishedDate": "2017-03-30T06:54:02.000Z",
  "requestUrl": "http://test.api.ft.com/content/22c0d426-1466-11e7-b0c1-37e417ee6c76",
  "standfirst": "Prime minister sets out Britain’s negotiating stance in statement to MPs",
  "standout": {
    "editorsChoice": false,
    "exclusive": false,
    "scoop": false
  },
  "title": "Brexit begins as Theresa May triggers Article 50",
  "type": "http://www.ft.com/ontology/content/Article"
}
```


## Example 2 (alternative images)
POST: /content
Body:
```
{
  "id": "http://www.ft.com/thing/6e1b070e-027b-11e7-ace0-1ce02ef0def9",
  "type": "http://www.ft.com/ontology/content/Content",
  "title": "Relaxed takes on tailoring for men",
  "alternativeTitles": {
    "promotionalTitle": "Relaxed takes on tailoring for men"
  },
  "alternativeStandfirsts": {
    "promotionalStandfirst": "This season’s fluid, super-comfortable tailoring is right on the sartorial money"
  },
  "publishedDate": "2017-03-08T10:39:48.000Z",
  "webUrl": "https://howtospendit.ft.com/mens-style/200263-relaxed-fluid-men-s-tailoring",
  "identifiers": [
    {
      "authority": "http://api.ft.com/system/FTCOM-METHODE",
      "identifierValue": "6e1b070e-027b-11e7-ace0-1ce02ef0def9"
    }
  ],
  "requestUrl": "http://test.api.ft.com/content/6e1b070e-027b-11e7-ace0-1ce02ef0def9",
  "brands": [
    "http://api.ft.com/things/dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54"
  ],
  "alternativeImages": {
    "promotionalImage": "http://test.api.ft.com/content/639cd952-149f-11e7-2ea7-a07ecd9ac73f"
  },
  "publishReference": "UK-V-5th-PP-CPH-Scenario-05-Trail-01",
  "lastModified": "2017-04-05T14:58:57.016Z",
  "canBeSyndicated": "verify",
  "canBeDistributed": "verify"
}
```

Response:
```
{
  "alternativeImages": {
    "promotionalImage": {
      "alternativeImages": {},
      "alternativeStandfirsts": {},
      "alternativeTitles": {},
      "canBeDistributed": "verify",
      "description": "Donald Tusk, president of the European Union (EU), holds the letter invoking Article 50 of the Lisbon Treaty from U.K. Prime Minister Theresa May as leaves following a news conference at the European Council in Brussels, Belgium, on Wednesday, March 29, 2017. The U.K. will&nbsp;start the clock&nbsp;on two years of negotiations to withdraw from the European Union on Wednesday, when Britain's ambassador hands EU President Donald Tusk&nbsp;a hand-signed&nbsp;letter&nbsp;from Prime Minister&nbsp;Theresa May&nbsp;invoking Article 50 of the Lisbon Treaty, the legal exit mechanism. Photographer: Jasper Juinen/Bloomberg",
      "firstPublishedDate": "2017-03-29T19:39:00.000Z",
      "id": "http://www.ft.com/thing/639cd952-149f-11e7-2ea7-a07ecd9ac73f",
      "identifiers": [
        {
          "authority": "http://api.ft.com/system/FTCOM-METHODE",
          "identifierValue": "639cd952-149f-11e7-2ea7-a07ecd9ac73f"
        }
      ],
      "lastModified": "2017-03-29T19:39:31.361Z",
      "members": [
        {
          "alternativeImages": {},
          "alternativeStandfirsts": {},
          "alternativeTitles": {},
          "binaryUrl": "http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/639cd952-149f-11e7-b0c1-37e417ee6c76",
          "canBeDistributed": "verify",
          "copyright": {
            "notice": "© Bloomberg"
          },
          "description": "Donald Tusk, president of the European Union (EU), holds the letter invoking Article 50 of the Lisbon Treaty from U.K. Prime Minister Theresa May as leaves following a news conference at the European Council in Brussels, Belgium, on Wednesday, March 29, 2017. The U.K. will&nbsp;start the clock&nbsp;on two years of negotiations to withdraw from the European Union on Wednesday, when Britain's ambassador hands EU President Donald Tusk&nbsp;a hand-signed&nbsp;letter&nbsp;from Prime Minister&nbsp;Theresa May&nbsp;invoking Article 50 of the Lisbon Treaty, the legal exit mechanism. Photographer: Jasper Juinen/Bloomberg",
          "firstPublishedDate": "2017-03-29T19:39:00.000Z",
          "id": "http://www.ft.com/thing/639cd952-149f-11e7-b0c1-37e417ee6c76",
          "identifiers": [
            {
              "authority": "http://api.ft.com/system/FTCOM-METHODE",
              "identifierValue": "639cd952-149f-11e7-b0c1-37e417ee6c76"
            }
          ],
          "lastModified": "2017-03-29T19:39:31.361Z",
          "pixelHeight": 1152,
          "pixelWidth": 2048,
          "publishReference": "tid_5ypvntzcpu",
          "publishedDate": "2017-03-29T19:39:00.000Z",
          "requestUrl": "http://test.api.ft.com/content/639cd952-149f-11e7-b0c1-37e417ee6c76",
          "title": "",
          "type": "http://www.ft.com/ontology/content/MediaResource"
        }
      ],
      "publishReference": "tid_5ypvntzcpu",
      "publishedDate": "2017-03-29T19:39:00.000Z",
      "requestUrl": "http://test.api.ft.com/content/639cd952-149f-11e7-2ea7-a07ecd9ac73f",
      "title": "",
      "type": "http://www.ft.com/ontology/content/ImageSet"
    }
  },
  "alternativeStandfirsts": {
    "promotionalStandfirst": "This season’s fluid, super-comfortable tailoring is right on the sartorial money"
  },
  "alternativeTitles": {
    "promotionalTitle": "Relaxed takes on tailoring for men"
  },
  "brands": [
    "http://api.ft.com/things/dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54"
  ],
  "canBeDistributed": "verify",
  "canBeSyndicated": "verify",
  "id": "http://www.ft.com/thing/6e1b070e-027b-11e7-ace0-1ce02ef0def9",
  "identifiers": [
    {
      "authority": "http://api.ft.com/system/FTCOM-METHODE",
      "identifierValue": "6e1b070e-027b-11e7-ace0-1ce02ef0def9"
    }
  ],
  "lastModified": "2017-04-05T14:58:57.016Z",
  "publishReference": "UK-V-5th-PP-CPH-Scenario-05-Trail-01",
  "publishedDate": "2017-03-08T10:39:48.000Z",
  "requestUrl": "http://test.api.ft.com/content/6e1b070e-027b-11e7-ace0-1ce02ef0def9",
  "title": "Relaxed takes on tailoring for men",
  "type": "http://www.ft.com/ontology/content/Content",
  "webUrl": "https://howtospendit.ft.com/mens-style/200263-relaxed-fluid-men-s-tailoring"
}
```

## Example 3 (lead images)
POST: /internalcontent
Body:
```
{
  "design": null,
  "tableOfContents": null,
  "topper": null,
  "leadImages": [
    {
      "id": "89f194c8-13bc-11e7-80f4-13e067d5072c",
      "type": "square"
    },
    {
      "id": "3e96c818-13bc-11e7-b0c1-37e417ee6c76",
      "type": "standard"
    },
    {
      "id": "8d7b4e22-13bc-11e7-80f4-13e067d5072c",
      "type": "wide"
    }
  ],
  "uuid": "5010e2e4-09bd-11e7-97d1-5e720a26771b",
  "lastModified": "2017-03-31T08:23:37.061Z",
  "publishReference": "tid_8pqiiuxbvz"
}
```

Response:
```
{
  "design": null,
  "lastModified": "2017-03-31T08:23:37.061Z",
  "leadImages": [
    {
      "id": "89f194c8-13bc-11e7-80f4-13e067d5072c",
      "image": {
        "alternativeImages": {},
        "alternativeStandfirsts": {},
        "alternativeTitles": {},
        "binaryUrl": "http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/89f194c8-13bc-11e7-80f4-13e067d5072c",
        "canBeDistributed": "verify",
        "description": "",
        "firstPublishedDate": "2017-03-28T13:45:00.000Z",
        "id": "http://www.ft.com/thing/89f194c8-13bc-11e7-80f4-13e067d5072c",
        "identifiers": [
          {
            "authority": "http://api.ft.com/system/FTCOM-METHODE",
            "identifierValue": "89f194c8-13bc-11e7-80f4-13e067d5072c"
          }
        ],
        "lastModified": "2017-03-28T13:45:53.438Z",
        "pixelHeight": 2612,
        "pixelWidth": 2612,
        "publishReference": "tid_lej6rdjegj",
        "publishedDate": "2017-03-28T13:45:00.000Z",
        "requestUrl": "http://test.api.ft.com/content/89f194c8-13bc-11e7-80f4-13e067d5072c",
        "title": "",
        "type": "http://www.ft.com/ontology/content/MediaResource"
      },
      "type": "square"
    },
    {
      "id": "3e96c818-13bc-11e7-b0c1-37e417ee6c76",
      "image": {
        "alternativeImages": {},
        "alternativeStandfirsts": {},
        "alternativeTitles": {},
        "binaryUrl": "http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/3e96c818-13bc-11e7-b0c1-37e417ee6c76",
        "canBeDistributed": "verify",
        "copyright": {
          "notice": "© EPA"
        },
        "description": "",
        "firstPublishedDate": "2017-03-28T13:42:00.000Z",
        "id": "http://www.ft.com/thing/3e96c818-13bc-11e7-b0c1-37e417ee6c76",
        "identifiers": [
          {
            "authority": "http://api.ft.com/system/FTCOM-METHODE",
            "identifierValue": "3e96c818-13bc-11e7-b0c1-37e417ee6c76"
          }
        ],
        "lastModified": "2017-03-28T13:43:00.375Z",
        "pixelHeight": 1152,
        "pixelWidth": 2048,
        "publishReference": "tid_tv7jfgi6jn",
        "publishedDate": "2017-03-28T13:42:00.000Z",
        "requestUrl": "http://test.api.ft.com/content/3e96c818-13bc-11e7-b0c1-37e417ee6c76",
        "title": "Leader of the PVV party Gert Wilders reacts to the election result",
        "type": "http://www.ft.com/ontology/content/MediaResource"
      },
      "type": "standard"
    },
    {
      "id": "8d7b4e22-13bc-11e7-80f4-13e067d5072c",
      "image": {
        "alternativeImages": {},
        "alternativeStandfirsts": {},
        "alternativeTitles": {},
        "binaryUrl": "http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/8d7b4e22-13bc-11e7-80f4-13e067d5072c",
        "canBeDistributed": "verify",
        "description": "",
        "firstPublishedDate": "2017-03-28T13:45:00.000Z",
        "id": "http://www.ft.com/thing/8d7b4e22-13bc-11e7-80f4-13e067d5072c",
        "identifiers": [
          {
            "authority": "http://api.ft.com/system/FTCOM-METHODE",
            "identifierValue": "8d7b4e22-13bc-11e7-80f4-13e067d5072c"
          }
        ],
        "lastModified": "2017-03-28T13:45:53.525Z",
        "pixelHeight": 1548,
        "pixelWidth": 4645,
        "publishReference": "tid_prlsj2avbn",
        "publishedDate": "2017-03-28T13:45:00.000Z",
        "requestUrl": "http://test.api.ft.com/content/8d7b4e22-13bc-11e7-80f4-13e067d5072c",
        "title": "",
        "type": "http://www.ft.com/ontology/content/MediaResource"
      },
      "type": "wide"
    }
  ],
  "publishReference": "tid_8pqiiuxbvz",
  "tableOfContents": null,
  "topper": null,
  "uuid": "5010e2e4-09bd-11e7-97d1-5e720a26771b"
}
```