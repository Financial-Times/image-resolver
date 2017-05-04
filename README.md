[![CircleCI](https://circleci.com/gh/Financial-Times/image-resolver/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/image-resolver/tree/master)
[![Coverage Status](https://coveralls.io/repos/github/Financial-Times/image-resolver/badge.svg)](https://coveralls.io/github/Financial-Times/image-resolver)
# image-resolver

Image resolver is an internally used API for retrieving unrolled images. It receive a content and return the content plus images unrolled.

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
 host name. e.g. coco.pre-prod.public-things-api.1 ($GRAPHITE_PREFIX)
  --logMetrics=false                          Whether to log metrics. Set to true if running locally and you want metric
s output ($LOG_METRICS)

```

## Endpoints

### Application specific endpoints:

* /content

### Admin specific endpoints:

* /ping
* /build-info
* /__ping
* /__build-info
* /__health
* /__gtg


## Example 1 (main image)
POST: /image-resolver/content
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
    "id": "http://www.ft.com/thing/639cd952-149f-11e7-2ea7-a07ecd9ac73f",
    "type": "http://www.ft.com/ontology/content/ImageSet",
    "title": "",
    "alternativeTitles": {},
    "alternativeStandfirsts": {},
    "description": "Donald Tusk, president of the European Union (EU), holds the letter invoking Article 50 of the Lisbon Treaty from U.K. Prime Minister Theresa May as leaves following a news conference at the European Council in Brussels, Belgium, on Wednesday, March 29, 2017. The U.K. will&nbsp;start the clock&nbsp;on two years of negotiations to withdraw from the European Union on Wednesday, when Britain's ambassador hands EU President Donald Tusk&nbsp;a hand-signed&nbsp;letter&nbsp;from Prime Minister&nbsp;Theresa May&nbsp;invoking Article 50 of the Lisbon Treaty, the legal exit mechanism. Photographer: Jasper Juinen/Bloomberg",
    "firstPublishedDate": "2017-03-29T19:39:00.000Z",
    "publishedDate": "2017-03-29T19:39:00.000Z",
    "identifiers": [
      {
        "authority": "http://api.ft.com/system/FTCOM-METHODE",
        "identifierValue": "639cd952-149f-11e7-2ea7-a07ecd9ac73f"
      }
    ],
    "members": [
      {
        "id": "http://www.ft.com/thing/639cd952-149f-11e7-b0c1-37e417ee6c76",
        "type": "http://www.ft.com/ontology/content/MediaResource",
        "title": "",
        "alternativeTitles": {},
        "alternativeStandfirsts": {},
        "description": "Donald Tusk, president of the European Union (EU), holds the letter invoking Article 50 of the Lisbon Treaty from U.K. Prime Minister Theresa May as leaves following a news conference at the European Council in Brussels, Belgium, on Wednesday, March 29, 2017. The U.K. will&nbsp;start the clock&nbsp;on two years of negotiations to withdraw from the European Union on Wednesday, when Britain's ambassador hands EU President Donald Tusk&nbsp;a hand-signed&nbsp;letter&nbsp;from Prime Minister&nbsp;Theresa May&nbsp;invoking Article 50 of the Lisbon Treaty, the legal exit mechanism. Photographer: Jasper Juinen/Bloomberg",
        "firstPublishedDate": "2017-03-29T19:39:00.000Z",
        "publishedDate": "2017-03-29T19:39:00.000Z",
        "identifiers": [
          {
            "authority": "http://api.ft.com/system/FTCOM-METHODE",
            "identifierValue": "639cd952-149f-11e7-b0c1-37e417ee6c76"
          }
        ],
        "requestUrl": "http://test.api.ft.com/content/639cd952-149f-11e7-b0c1-37e417ee6c76",
        "binaryUrl": "http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/639cd952-149f-11e7-b0c1-37e417ee6c76",
        "pixelWidth": 2048,
        "pixelHeight": 1152,
        "alternativeImages": {},
        "copyright": {
          "notice": "© Bloomberg"
        },
        "publishReference": "tid_5ypvntzcpu",
        "lastModified": "2017-03-29T19:39:31.361Z",
        "canBeDistributed": "verify"
      }
    ],
    "requestUrl": "http://test.api.ft.com/content/639cd952-149f-11e7-2ea7-a07ecd9ac73f",
    "alternativeImages": {},
    "publishReference": "tid_5ypvntzcpu",
    "lastModified": "2017-03-29T19:39:31.361Z",
    "canBeDistributed": "verify"
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


## Example 2 (alternative images)
POST: /image-resolver/content
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
    "promotionalImage": "http://test.api.ft.com/content/4723cb4e-027c-11e7-ace0-1ce02ef0def9"
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
    "promotionalImage": {
      "id": "http://www.ft.com/thing/4723cb4e-027c-11e7-ace0-1ce02ef0def9",
      "type": "http://www.ft.com/ontology/content/MediaResource",
      "alternativeTitles": {},
      "alternativeStandfirsts": {},
      "firstPublishedDate": "2017-03-06T14:50:00.000Z",
      "publishedDate": "2017-03-06T14:50:00.000Z",
      "identifiers": [
        {
          "authority": "http://api.ft.com/system/FTCOM-METHODE",
          "identifierValue": "4723cb4e-027c-11e7-ace0-1ce02ef0def9"
        }
      ],
      "requestUrl": "http://test.api.ft.com/content/4723cb4e-027c-11e7-ace0-1ce02ef0def9",
      "binaryUrl": "http://com.ft.imagepublish.prod-us.s3.amazonaws.com/4723cb4e-027c-11e7-ace0-1ce02ef0def9",
      "alternativeImages": {},
      "publishReference": "tid_axufgrmrhm",
      "lastModified": "2017-03-06T14:50:50.298Z"
    }
  },
  "publishReference": "UK-V-5th-PP-CPH-Scenario-05-Trail-01",
  "lastModified": "2017-04-05T14:58:57.016Z",
  "canBeSyndicated": "verify",
  "canBeDistributed": "verify"
}
```

## Example 3 (lead images)
POST: /image-resolver/content
Body:
```
{
  "id": "http://www.ft.com/thing/4da6a172-2431-11e7-a24c-6bbb6ec0bc98",
  "type": "http://www.ft.com/ontology/content/Article",
  "bodyXML": "<body><p>Test Lead Image (all 3 formats)</p>\n\n\n</body>",
  "title": "Lead Image",
  "alternativeTitles": {},
  "alternativeStandfirsts": {},
  "byline": "TL",
  "firstPublishedDate": "2017-04-18T12:04:28.000Z",
  "publishedDate": "2017-04-18T12:04:28.000Z",
  "identifiers": [
    {
      "authority": "http://api.ft.com/system/FTCOM-METHODE",
      "identifierValue": "4da6a172-2431-11e7-a24c-6bbb6ec0bc98"
    }
  ],
  "requestUrl": "http://test.api.ft.com/content/4da6a172-2431-11e7-a24c-6bbb6ec0bc98",
  "brands": [
    "http://api.ft.com/things/dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54"
  ],
  "alternativeImages": {},
  "comments": {
    "enabled": true
  },
  "standout": {
    "editorsChoice": false,
    "exclusive": false,
    "scoop": false
  },
  "publishReference": "tid_hcqcors4wf",
  "lastModified": "2017-04-18T12:23:32.744Z",
  "canBeDistributed": "yes",
  "canBeSyndicated": "verify",
  "accessLevel": "subscribed",
  "leadImages": [
    {
      "id": "588d9ba6-1557-11e7-9469-afea892e4de3",
      "type": "square"
    },
    {
      "id": "588d9ba6-1557-11e7-9469-afea892e4de3",
      "type": "standard"
    },
    {
      "id": "588d9ba6-1557-11e7-9469-afea892e4de3",
      "type": "wide"
    }
  ]
}
```

Response:
```
{
  "id": "http://www.ft.com/thing/4da6a172-2431-11e7-a24c-6bbb6ec0bc98",
  "type": "http://www.ft.com/ontology/content/Article",
  "bodyXML": "<body><p>Test Lead Image (all 3 formats)</p>\n\n\n</body>",
  "title": "Lead Image",
  "alternativeTitles": {},
  "alternativeStandfirsts": {},
  "byline": "TL",
  "firstPublishedDate": "2017-04-18T12:04:28.000Z",
  "publishedDate": "2017-04-18T12:04:28.000Z",
  "identifiers": [
    {
      "authority": "http://api.ft.com/system/FTCOM-METHODE",
      "identifierValue": "4da6a172-2431-11e7-a24c-6bbb6ec0bc98"
    }
  ],
  "requestUrl": "http://test.api.ft.com/content/4da6a172-2431-11e7-a24c-6bbb6ec0bc98",
  "brands": [
    "http://api.ft.com/things/dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54"
  ],
  "alternativeImages": {},
  "comments": {
    "enabled": true
  },
  "standout": {
    "editorsChoice": false,
    "exclusive": false,
    "scoop": false
  },
  "publishReference": "tid_hcqcors4wf",
  "lastModified": "2017-04-18T12:23:32.744Z",
  "canBeSyndicated": "verify",
  "accessLevel": "subscribed",
  "canBeDistributed": "yes",
  "leadImages": [
    {
      "image": {
        "id": "http://www.ft.com/thing/588d9ba6-1557-11e7-9469-afea892e4de3",
        "type": "http://www.ft.com/ontology/content/MediaResource",
        "title": "",
        "alternativeTitles": {},
        "alternativeStandfirsts": {},
        "description": "",
        "firstPublishedDate": "2017-04-05T08:25:00.000Z",
        "publishedDate": "2017-04-05T08:25:00.000Z",
        "identifiers": [
          {
            "authority": "http://api.ft.com/system/FTCOM-METHODE",
            "identifierValue": "588d9ba6-1557-11e7-9469-afea892e4de3"
          }
        ],
        "requestUrl": "http://test.api.ft.com/content/588d9ba6-1557-11e7-9469-afea892e4de3",
        "binaryUrl": "http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/588d9ba6-1557-11e7-9469-afea892e4de3",
        "pixelWidth": 4645,
        "pixelHeight": 2612,
        "alternativeImages": {},
        "publishReference": "tid_n6vkawrq1l",
        "lastModified": "2017-04-05T08:25:20.438Z",
        "canBeDistributed": "verify"
      },
      "type": "square"
    },
    {
      "image": {
        "id": "http://www.ft.com/thing/588d9ba6-1557-11e7-9469-afea892e4de3",
        "type": "http://www.ft.com/ontology/content/MediaResource",
        "title": "",
        "alternativeTitles": {},
        "alternativeStandfirsts": {},
        "description": "",
        "firstPublishedDate": "2017-04-05T08:25:00.000Z",
        "publishedDate": "2017-04-05T08:25:00.000Z",
        "identifiers": [
          {
            "authority": "http://api.ft.com/system/FTCOM-METHODE",
            "identifierValue": "588d9ba6-1557-11e7-9469-afea892e4de3"
          }
        ],
        "requestUrl": "http://test.api.ft.com/content/588d9ba6-1557-11e7-9469-afea892e4de3",
        "binaryUrl": "http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/588d9ba6-1557-11e7-9469-afea892e4de3",
        "pixelWidth": 4645,
        "pixelHeight": 2612,
        "alternativeImages": {},
        "publishReference": "tid_n6vkawrq1l",
        "lastModified": "2017-04-05T08:25:20.438Z",
        "canBeDistributed": "verify"
      },
      "type": "standard"
    },
    {
      "image": {
        "id": "http://www.ft.com/thing/588d9ba6-1557-11e7-9469-afea892e4de3",
        "type": "http://www.ft.com/ontology/content/MediaResource",
        "title": "",
        "alternativeTitles": {},
        "alternativeStandfirsts": {},
        "description": "",
        "firstPublishedDate": "2017-04-05T08:25:00.000Z",
        "publishedDate": "2017-04-05T08:25:00.000Z",
        "identifiers": [
          {
            "authority": "http://api.ft.com/system/FTCOM-METHODE",
            "identifierValue": "588d9ba6-1557-11e7-9469-afea892e4de3"
          }
        ],
        "requestUrl": "http://test.api.ft.com/content/588d9ba6-1557-11e7-9469-afea892e4de3",
        "binaryUrl": "http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/588d9ba6-1557-11e7-9469-afea892e4de3",
        "pixelWidth": 4645,
        "pixelHeight": 2612,
        "alternativeImages": {},
        "publishReference": "tid_n6vkawrq1l",
        "lastModified": "2017-04-05T08:25:20.438Z",
        "canBeDistributed": "verify"
      },
      "type": "wide"
    }
  ]
}
```