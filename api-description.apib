FORMAT: 1A

# MOPH

Methode Publish Handler intercepts methode publication messages, and does some extra enrichment prior to forwarding on to the rest of the UPP stack (CMS Notifier etc.).

# GET /__health

Returns the results of the healthchecks, which should lightly test the components dependencies.

+ Response 200 (application/json)

        {
          "checks": [
            {
              "businessImpact": "No Articles from PortalPub will be published!",
              "checkOutput": "cms-notifier service is unreachable",
              "lastUpdated": "2016-08-24T15:46:11.33737686+01:00",
              "name": "cms-notifier Availabililty Check",
              "ok": false,
              "panicGuide": "",
              "severity": 1,
              "technicalSummary": "Checks that cms-notifier Service is reachable. Methode Publish Handler forwards published articles onto cms-notifier."
            }
          ],
          "description": "A RESTful API which accepts Methode Articles and appends a vanity url before passing on to the UPP Stack",
          "name": "methode-publish-handler",
          "schemaVersion": 1,
          "ok": false,
          "severity": 1
        }

# GET /__build-info

Displays build info for the component.

+ Response 200 (application/json; charset=UTF-8)

        {
          "version": "",
          "repository": "",
          "revision": "",
          "builder": "",
          "dateTime": ""
        }

# POST /notify

Simulates a direct call to the CMS Notifier, but instead calls the Vanity Service, and appends the webUrl to the article output. It accepts a Methode JSON document.

+ Request (application/json)

        {
          "systemAttributes":"<systemAttributes></systemAttributes>",
          "lastModified": "2016-08-30T09:40:50Z",
          "uuid": "12f01e0a-c404-11e2-aa5b-00144feab7de",
          "type": "EOM::WebContainer",
          "workflowStatus": "some workflow status",
          "usageTickets": "<usageTickets></usageTickets>",
          "linkedObjects": ["linked-object data"],
          "value": "base64 value",
          "attributes": "<attributes></attributes>"
        }

+ Response 200 (application/json)

        {
          "json": {
            "systemAttributes":"",
            "lastModified": "",
            "uuid": "",
            "type": "",
            "workflowStatus": "",
            "usageTickets": "",
            "linkedObjects": [],
            "value": "",
            "attributes": "",
            "webUrl": ""
          }
        }
