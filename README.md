# Apartment Alerter
A program scrapping Kamernet, and alerting you by text (using Twilio). If you're looking for a room in Amsterdam, this might help you.

## Run with Docker
1. Create a configuration file, like this one:

```yaml
twilio:
  enable: true
  from-number: "<from number, registered in Twilio>"
  to-number: "<to number>"
  sid: "<Twilio SID>"
  token: "<Twilio Auth Token>"
```

To see more configuration options, run `docker run ericgln/apartment-alert --help`.

2. Run:
```sh
docker run -d -v '/PATH/TO/CONFIG/FILE.yml:/root/.apartment-alert/yml' ericgln/apartment-alert
```

This will run the program in the background.
