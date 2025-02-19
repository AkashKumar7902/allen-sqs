This is a modified gin mongo app, with three new function that does the following thing

* FetchAppConfig() : periodically fetches configuration from AWS AppConfig every 10 sec.
* SendEvent(): creates random event and publishes to sqs server
* ConsumeEvents(): consumes the event coming to the sqs server and does mongo call with event data.
  