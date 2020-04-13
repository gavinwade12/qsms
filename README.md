qsms (quick-sms) is a utility for quickly sending an SMS message to someone.

# Installation

Installing qsms currently requires you to have Go installed.

    git clone https://github.com/gavinwade12/qsms
    cd qsms
    go install

# Gateways

There are currently 3 gateways supported: email, Twilio, and Messages (OSX only). The email and Twilio gateways require configuration; however, enough fields can be configured through the cli to get you going with a gateway the first time you use it. The Messages gateway requires no configuration at all.

- [Configuring the Email Gateway](#configure-email-gateway)
- [Configuring the Twilio Gateway](#configure-twilio-gateway)

# Usage

Using qsms is simple. Just provide the recipient number and the text to be sent.

    qsms 4195551234 "Hello World!"

Provide a gateway if you don't want to use your default.

    qsms -g twilio 4195551234 "Hello World!"

Send a text to a list of friends.

    printf "adam,+1 (123) 456-7890\nbeth,987654321\ncharles,419-555-1234\n" > friends.csv
    while IFS=, read -r name number; do qsms $number "Hello, $name."; done < friends.csv

# Configuration

The config file for qsms is a json file located at `$HOME/.qsms.json`. Most of the time, the cli will prompt for any missing values automatically, but there are cases where this may need edited manually.

### Configuring the Email Gateway

The email gateway has 4 fields that need configured:

- Email - the email address used to send the SMS
- Password - the password for the email address*
- SMTP Server - the SMTP server you'd like to send the SMS with e.g. smtp.gmail.com if you're using a Gmail account
- SMTP Server Port - the port your SMTP server uses e.g. 587 if using the Gmail SMTP server

These fields can be filled out through the cli the first time you try using the email gateway.

The email gateway configuration also contains a carrier -> domain mapping. Currently, the only carrier supported is Verizon. Adding other carriers to the configuration and using them is experimental at this point.

**The password is encoded to base32 before being stored in the config, so you will either need to set it via the cli or encode it before manually changing the config file.*

***If you're having issues with your Gmail credentials, you may need to [create an app password](https://myaccount.google.com/apppasswords).*

### Configuring the Twilio Gateway

Twilio provides APIs for automating things like sending SMS messages. Each SMS message currently costs $0.0075, but it's typically much faster than the email gateway and doesn't require you to provide the recipient's carrier. If you haven't already signed up, you can use my [referral link](www.twilio.com/referral/viTdGG) to give us both $10.

The Twilio gateway has 3 fields that need configured:

- Account SID
- Auth Token
- Phone Number

The `Account SID` and `Auth Token` can be obtained from the [Settings page](https://www.twilio.com/console/project/settings) in your Twilio console. You can choose any `Phone Number` from the [Active Numbers page](https://www.twilio.com/console/phone-numbers/incoming).

### Example

Here's what a config file may look like when complete:

```json
{
  "default_gateway": "twilio",
  "gateways": {
    "email": {
      "mapping": {
        "verizon": "vtext.com"
      },
      "sender": "fake@gmail.com",
      "sender_password": "JNHX5CW#JMX25CU3NOWXYAUL3N======",
      "smtp_server": "smtp.gmail.com",
      "smtp_server_port": 587
    },
    "twilio": {
      "account_sid": "AC2527fac9f4d345e77738df3db06aa1234",
      "auth_token": "5b805c3d4277ea8c6971b2f82dd4ca1b",
      "from_number": "+14805550601"
    }
  }
}
```