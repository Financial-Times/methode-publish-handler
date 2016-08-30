# MOPH [![CircleCI](https://circleci.com/gh/Financial-Times/methode-publish-handler.svg?style=svg)](https://circleci.com/gh/Financial-Times/methode-publish-handler)

Methode Publish Handler intercepts methode publication messages, and does some extra enrichment prior to forwarding on to the rest of the UPP stack (CMS Notifier etc.).

## API

See the `api-description.apib` API Blueprint file for details. To test the API documentation is still valid, simply install DreddJS and run:

```
dredd
```
