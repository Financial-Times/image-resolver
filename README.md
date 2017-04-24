# image-resolver

Image resolver is an internally used API for retrieving unrolled images.

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

* /content/{uuid}

### Admin specific endpoints:

* /ping
* /build-info
* /__ping
* /__build-info
* /__health
* /__gtg


## Example
GET: https://xp-up.ft.com//__image-resolver/content/22c0d426-1466-11e7-b0c1-37e417ee6c76
```
{
   "mainImage":{
      "id":"http://www.ft.com/thing/639cd952-149f-11e7-2ea7-a07ecd9ac73f",
      "type":"http://www.ft.com/ontology/content/ImageSet",
      "description":"Donald Tusk, president of the European Union (EU), holds the letter invoking Article 50 of the Lisbon Treaty from U.K. Prime Minister Theresa May as leaves following a news conference at the European Council in Brussels, Belgium, on Wednesday, March 29, 2017. The U.K. will&nbsp;start the clock&nbsp;on two years of negotiations to withdraw from the European Union on Wednesday, when Britain's ambassador hands EU President Donald Tusk&nbsp;a hand-signed&nbsp;letter&nbsp;from Prime Minister&nbsp;Theresa May&nbsp;invoking Article 50 of the Lisbon Treaty, the legal exit mechanism. Photographer: Jasper Juinen/Bloomberg",
      "publishedDate":"2017-03-29T19:39:00Z",
      "identifiers":[
         {
            "authority":"http://api.ft.com/system/FTCOM-METHODE",
            "identifierValue":"639cd952-149f-11e7-2ea7-a07ecd9ac73f"
         }
      ],
      "members":[
         {
            "id":"http://www.ft.com/thing/639cd952-149f-11e7-b0c1-37e417ee6c76",
            "type":"http://www.ft.com/ontology/content/MediaResource",
            "description":"Donald Tusk, president of the European Union (EU), holds the letter invoking Article 50 of the Lisbon Treaty from U.K. Prime Minister Theresa May as leaves following a news conference at the European Council in Brussels, Belgium, on Wednesday, March 29, 2017. The U.K. will&nbsp;start the clock&nbsp;on two years of negotiations to withdraw from the European Union on Wednesday, when Britain's ambassador hands EU President Donald Tusk&nbsp;a hand-signed&nbsp;letter&nbsp;from Prime Minister&nbsp;Theresa May&nbsp;invoking Article 50 of the Lisbon Treaty, the legal exit mechanism. Photographer: Jasper Juinen/Bloomberg",
            "publishedDate":"2017-03-29T19:39:00Z",
            "identifiers":[
               {
                  "authority":"http://api.ft.com/system/FTCOM-METHODE",
                  "identifierValue":"639cd952-149f-11e7-b0c1-37e417ee6c76"
               }
            ],
            "requestUrl":"http://test.api.ft.com/content/639cd952-149f-11e7-b0c1-37e417ee6c76",
            "binaryUrl":"http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/639cd952-149f-11e7-b0c1-37e417ee6c76",
            "copyright":{
               "notice":"© Bloomberg"
            },
            "publishReference":"tid_5ypvntzcpu",
            "pixelWidth":2048,
            "pixelHeight":1152,
            "lastModified":"2017-03-29T19:39:31.361Z",
            "alternativeTitles":{

            },
            "alternativeStandfirsts":{

            },
            "alternativeImages":{

            },
            "firstPublishedDate":"2017-03-29T19:39:00Z",
            "canBeDistributed":"verify"
         }
      ],
      "requestUrl":"http://test.api.ft.com/content/639cd952-149f-11e7-2ea7-a07ecd9ac73f",
      "publishReference":"tid_5ypvntzcpu",
      "lastModified":"2017-03-29T19:39:31.361Z",
      "alternativeTitles":{

      },
      "alternativeStandfirsts":{

      },
      "alternativeImages":{

      },
      "firstPublishedDate":"2017-03-29T19:39:00Z",
      "canBeDistributed":"verify"
   },
   "embeds":[
      {
         "id":"http://www.ft.com/thing/639cd952-149f-11e7-2ea7-a07ecd9ac73f",
         "type":"http://www.ft.com/ontology/content/ImageSet",
         "description":"Donald Tusk, president of the European Union (EU), holds the letter invoking Article 50 of the Lisbon Treaty from U.K. Prime Minister Theresa May as leaves following a news conference at the European Council in Brussels, Belgium, on Wednesday, March 29, 2017. The U.K. will&nbsp;start the clock&nbsp;on two years of negotiations to withdraw from the European Union on Wednesday, when Britain's ambassador hands EU President Donald Tusk&nbsp;a hand-signed&nbsp;letter&nbsp;from Prime Minister&nbsp;Theresa May&nbsp;invoking Article 50 of the Lisbon Treaty, the legal exit mechanism. Photographer: Jasper Juinen/Bloomberg",
         "publishedDate":"2017-03-29T19:39:00Z",
         "identifiers":[
            {
               "authority":"http://api.ft.com/system/FTCOM-METHODE",
               "identifierValue":"639cd952-149f-11e7-2ea7-a07ecd9ac73f"
            }
         ],
         "members":[
            {
               "id":"http://www.ft.com/thing/639cd952-149f-11e7-b0c1-37e417ee6c76",
               "type":"http://www.ft.com/ontology/content/MediaResource",
               "description":"Donald Tusk, president of the European Union (EU), holds the letter invoking Article 50 of the Lisbon Treaty from U.K. Prime Minister Theresa May as leaves following a news conference at the European Council in Brussels, Belgium, on Wednesday, March 29, 2017. The U.K. will&nbsp;start the clock&nbsp;on two years of negotiations to withdraw from the European Union on Wednesday, when Britain's ambassador hands EU President Donald Tusk&nbsp;a hand-signed&nbsp;letter&nbsp;from Prime Minister&nbsp;Theresa May&nbsp;invoking Article 50 of the Lisbon Treaty, the legal exit mechanism. Photographer: Jasper Juinen/Bloomberg",
               "publishedDate":"2017-03-29T19:39:00Z",
               "identifiers":[
                  {
                     "authority":"http://api.ft.com/system/FTCOM-METHODE",
                     "identifierValue":"639cd952-149f-11e7-b0c1-37e417ee6c76"
                  }
               ],
               "requestUrl":"http://test.api.ft.com/content/639cd952-149f-11e7-b0c1-37e417ee6c76",
               "binaryUrl":"http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/639cd952-149f-11e7-b0c1-37e417ee6c76",
               "copyright":{
                  "notice":"© Bloomberg"
               },
               "publishReference":"tid_5ypvntzcpu",
               "pixelWidth":2048,
               "pixelHeight":1152,
               "lastModified":"2017-03-29T19:39:31.361Z",
               "alternativeTitles":{

               },
               "alternativeStandfirsts":{

               },
               "alternativeImages":{

               },
               "firstPublishedDate":"2017-03-29T19:39:00Z",
               "canBeDistributed":"verify"
            }
         ],
         "requestUrl":"http://test.api.ft.com/content/639cd952-149f-11e7-2ea7-a07ecd9ac73f",
         "publishReference":"tid_5ypvntzcpu",
         "lastModified":"2017-03-29T19:39:31.361Z",
         "alternativeTitles":{

         },
         "alternativeStandfirsts":{

         },
         "alternativeImages":{

         },
         "firstPublishedDate":"2017-03-29T19:39:00Z",
         "canBeDistributed":"verify"
      },
      {
         "id":"http://www.ft.com/thing/71231d3a-13c7-11e7-2ea7-a07ecd9ac73f",
         "type":"http://www.ft.com/ontology/content/ImageSet",
         "publishedDate":"2017-03-29T18:28:00Z",
         "identifiers":[
            {
               "authority":"http://api.ft.com/system/FTCOM-METHODE",
               "identifierValue":"71231d3a-13c7-11e7-2ea7-a07ecd9ac73f"
            }
         ],
         "members":[
            {
               "id":"http://www.ft.com/thing/71231d3a-13c7-11e7-b0c1-37e417ee6c76",
               "type":"http://www.ft.com/ontology/content/MediaResource",
               "publishedDate":"2017-03-29T18:28:00Z",
               "identifiers":[
                  {
                     "authority":"http://api.ft.com/system/FTCOM-METHODE",
                     "identifierValue":"71231d3a-13c7-11e7-b0c1-37e417ee6c76"
                  }
               ],
               "requestUrl":"http://test.api.ft.com/content/71231d3a-13c7-11e7-b0c1-37e417ee6c76",
               "binaryUrl":"http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/71231d3a-13c7-11e7-b0c1-37e417ee6c76",
               "copyright":{
                  "notice":"© FT montage; Getty Images"
               },
               "publishReference":"tid_fsdgcbtcih",
               "pixelWidth":2048,
               "pixelHeight":1152,
               "lastModified":"2017-03-29T18:28:38.571Z",
               "alternativeTitles":{

               },
               "alternativeStandfirsts":{

               },
               "alternativeImages":{

               },
               "firstPublishedDate":"2017-03-29T18:28:00Z",
               "canBeDistributed":"verify"
            }
         ],
         "requestUrl":"http://test.api.ft.com/content/71231d3a-13c7-11e7-2ea7-a07ecd9ac73f",
         "publishReference":"tid_fsdgcbtcih",
         "lastModified":"2017-03-29T18:28:38.571Z",
         "alternativeTitles":{

         },
         "alternativeStandfirsts":{

         },
         "alternativeImages":{

         },
         "firstPublishedDate":"2017-03-29T18:28:00Z",
         "canBeDistributed":"verify"
      },
      {
         "id":"http://www.ft.com/thing/0261ea4a-1474-11e7-1e92-847abda1ac65",
         "type":"http://www.ft.com/ontology/content/ImageSet",
         "title":"Tim Barrow, Britain's ambassador to the EU, delivers formal notice of the UK's intention to leave the bloc to European Council president Donald Tusk",
         "description":"Britain's ambassador to the EU Tim Barrow delivers British Prime Minister Theresa May's formal notice of the UK's intention to leave the bloc under Article 50 of the EU's Lisbon Treaty to European Council President Donald Tusk in Brussels on March 29, 2017. Britain formally launches the process for leaving the European Union on Wednesday, a historic step that has divided the country and thrown into question the future of the European unity project. / AFP PHOTO / POOL / Emmanuel DUNAND (Photo credit should read EMMANUEL DUNAND/AFP/Getty Images)",
         "publishedDate":"2017-03-29T18:28:00Z",
         "identifiers":[
            {
               "authority":"http://api.ft.com/system/FTCOM-METHODE",
               "identifierValue":"0261ea4a-1474-11e7-1e92-847abda1ac65"
            }
         ],
         "members":[
            {
               "id":"http://www.ft.com/thing/0261ea4a-1474-11e7-80f4-13e067d5072c",
               "type":"http://www.ft.com/ontology/content/MediaResource",
               "title":"Tim Barrow, Britain's ambassador to the EU, delivers formal notice of the UK's intention to leave the bloc to European Council president Donald Tusk",
               "description":"Britain's ambassador to the EU Tim Barrow delivers British Prime Minister Theresa May's formal notice of the UK's intention to leave the bloc under Article 50 of the EU's Lisbon Treaty to European Council President Donald Tusk in Brussels on March 29, 2017. Britain formally launches the process for leaving the European Union on Wednesday, a historic step that has divided the country and thrown into question the future of the European unity project. / AFP PHOTO / POOL / Emmanuel DUNAND (Photo credit should read EMMANUEL DUNAND/AFP/Getty Images)",
               "publishedDate":"2017-03-29T18:28:00Z",
               "identifiers":[
                  {
                     "authority":"http://api.ft.com/system/FTCOM-METHODE",
                     "identifierValue":"0261ea4a-1474-11e7-80f4-13e067d5072c"
                  }
               ],
               "requestUrl":"http://test.api.ft.com/content/0261ea4a-1474-11e7-80f4-13e067d5072c",
               "binaryUrl":"http://com.ft.coco-imagepublish.pre-prod.s3.amazonaws.com/0261ea4a-1474-11e7-80f4-13e067d5072c",
               "copyright":{
                  "notice":"© AFP"
               },
               "publishReference":"tid_4tjlxfiynp",
               "pixelWidth":2048,
               "pixelHeight":1152,
               "lastModified":"2017-03-29T18:28:38.623Z",
               "alternativeTitles":{

               },
               "alternativeStandfirsts":{

               },
               "alternativeImages":{

               },
               "firstPublishedDate":"2017-03-29T18:28:00Z",
               "canBeDistributed":"verify"
            }
         ],
         "requestUrl":"http://test.api.ft.com/content/0261ea4a-1474-11e7-1e92-847abda1ac65",
         "publishReference":"tid_4tjlxfiynp",
         "lastModified":"2017-03-29T18:28:38.623Z",
         "alternativeTitles":{

         },
         "alternativeStandfirsts":{

         },
         "alternativeImages":{

         },
         "firstPublishedDate":"2017-03-29T18:28:00Z",
         "canBeDistributed":"verify"
      }
   ]
}
```