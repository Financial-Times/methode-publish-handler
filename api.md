FORMAT: 1A

# Methode Publish Handler

Methode Publish Handler intercepts methode publication messages, and does some extra enrichment prior to forwarding on to the rest of the UPP stack (CMS Notifier etc.).

## Group API

### /notify

#### Proxies the configured CMS-Notifier. [POST]

Given a POST body of an article via Portal Pub; MOPH will lookup the Vanity url via an API call, append it to the article (as a new json field 'webUrl'), and forward the article on to the configured CMS Notifier for processing in the UPP stack.

+ Request (application/json)

    + Headers

            Accept: application/json

    + Body

            {
              "systemAttributes": "ipsum Duis",
              "lastModified": "3232-03-20T21:23:49.810Z",
              "uuid": "48936c08-9ee0-dbe1-8cd7-e974ed7b4be7",
              "type": "esse deserunt exercitation",
              "workflowStatus": "irure esse sit",
              "usageTickets": "nostrud",
              "linkedObjects": [
                "enim ullamco dolore"
              ],
              "value": "a",
              "attributes": "irure tempor"
            }

    + Schema

            {
              "type": "object",
              "required": [
                "systemAttributes",
                "lastModified",
                "uuid",
                "type",
                "workflowStatus",
                "usageTickets",
                "linkedObjects",
                "value",
                "attributes"
              ],
              "properties": {
                "systemAttributes": {
                  "type": "string",
                  "description": "CMS system attributes as XML."
                },
                "lastModified": {
                  "type": "string",
                  "format": "date-time",
                  "description": "Last modified date of the article."
                },
                "uuid": {
                  "type": "string",
                  "pattern": "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
                  "description": "Unique ID of the article."
                },
                "type": {
                  "type": "string",
                  "description": "CMS type of the article."
                },
                "workflowStatus": {
                  "type": "string",
                  "description": "Status of the article in the CMS workflow."
                },
                "usageTickets": {
                  "type": "string",
                  "description": "Usage tickets as XML."
                },
                "linkedObjects": {
                  "type": "array",
                  "description": "Linked articles as XML.",
                  "items": {
                    "type": "string"
                  }
                },
                "value": {
                  "type": "string",
                  "description": "Base64 representation of the full article."
                },
                "attributes": {
                  "type": "string",
                  "description": "Article attributes as XML."
                },
                "webUrl": {
                  "type": "string",
                  "format": "uri",
                  "readOnly": true,
                  "description": "The vanity url for the service."
                }
              }
            }

+ Response 200 (application/json)

    CMS Notifier returned 200 for the article.

    + Headers

            X-Request-Id: 

    + Body

            {}

    + Schema

            {
              "type": "object"
            }

## Group Meta

### /__health

#### Health of the Service [GET]

Returns the results of all health checks in JSON format.

+ Request

    + Headers

            Accept: application/json

    + Body

+ Response 200 (application/json)

    Results of the healthchecks for this service.

    + Body

            {
              "checks": [
                {
                  "lastUpdated": "sint",
                  "technicalSummary": "eiusmod",
                  "panicGuide": "eiusmod proident qui in nost",
                  "businessImpact": "sit consequat sunt Excepteur",
                  "checkOutput": "ipsum incididunt amet",
                  "ok": true,
                  "severity": 75801957,
                  "name": "de"
                },
                {
                  "technicalSummary": "exercitation eu ut",
                  "severity": 42923798,
                  "name": "ulla",
                  "ok": false,
                  "panicGuide": "irure Lorem",
                  "checkOutput": "ea deserunt laboris",
                  "lastUpdated": "sed commodo exercitation consectetur"
                },
                {
                  "lastUpdated": "non ex",
                  "businessImpact": "aliquip Ut qui do nulla"
                }
              ],
              "schemaVersion": -66094298
            }

    + Schema

            {
              "type": "object",
              "properties": {
                "checks": {
                  "type": "array",
                  "items": {
                    "type": "object",
                    "properties": {
                      "businessImpact": {
                        "type": "string",
                        "description": "A short statement of what the impact is if the check fails."
                      },
                      "checkOutput": {
                        "type": "string",
                        "description": "The technical output of the check."
                      },
                      "lastUpdated": {
                        "type": "string",
                        "description": "When was the check last run."
                      },
                      "name": {
                        "type": "string",
                        "description": "Name of the check."
                      },
                      "ok": {
                        "type": "boolean",
                        "description": "Did the check pass?"
                      },
                      "panicGuide": {
                        "type": "string",
                        "description": "Link to the panic guide for the corresponding service."
                      },
                      "severity": {
                        "type": "integer",
                        "description": "Severity of the problem. 1 is Critical, 2 is Warning."
                      },
                      "technicalSummary": {
                        "type": "string",
                        "description": "A short description of what the check does."
                      }
                    }
                  }
                },
                "description": {
                  "type": "string",
                  "description": "A brief description of what this component does."
                },
                "name": {
                  "type": "string",
                  "description": "The name of the component."
                },
                "schemaVersion": {
                  "type": "integer",
                  "description": "The schema version."
                },
                "ok": {
                  "type": "boolean",
                  "description": "Did any checks fail?"
                },
                "severity": {
                  "type": "integer",
                  "description": "Overall serverity for the service. 1 is Critical, 2 is Warning."
                }
              }
            }

### /__build-info

#### Returns detailed build information for the binary. [GET]

+ Response 200

    Displays build info in JSON.

    + Body

    + Schema

            {
              "type": "object",
              "properties": {
                "version": {
                  "type": "string",
                  "description": "Component version."
                },
                "repository": {
                  "type": "string",
                  "description": "Git Repository URL."
                },
                "revision": {
                  "type": "string",
                  "description": "Git revision."
                },
                "builder": {
                  "type": "string",
                  "description": "The version of go used to build the binary."
                },
                "dateTime": {
                  "type": "string",
                  "description": "Build time."
                }
              }
            }

### /__ping

#### Returns 'pong' if running. [GET]

+ Response 200 (text/plain; charset=utf-8)

    Returns 'pong' if running.

    + Body

            pong

