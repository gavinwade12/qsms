qsms (quick-sms) is a utility for quickly sending an SMS message to someone.

# Installation

Installing qsms currently requires you to have Go installed.

    git clone https://github.com/gavinwade12/qsms
    cd qsms
    go install

# Gateways

Currently, only 2 gateways are supported: email and Twilio (3 if you're a Mac user - there's support for sending directly through the Messages application). The email and Twilio gateways require configuration. All but one of the email fields can be configured through the cli the first time you try to use the gateway, and the Twilio gateway can be completely configured through the cli. The Messages gateway requires no configuration.

- [Configuring the Email Gateway](#configure-email-gateway)
- [Configuring the Twilio Gateway](#configure-twilio-gateway)

# Usage

Using qsms is simple. Just provide the recipient number and the text to be sent.

    qsms 4195551234 "Hello World!"

Provide a gateway if you don't want to use your default.

    qsms -g twilio 4195551234 "Hello World!"

# Configuration

The config file for qsms is a json file located at $HOME/.qsms.json. Most of the time, the cli will prompt for any missing values automatically, but there are cases where this may need edited manually.

### Configure Email Gateway

Most of the fields for the email gateway should be self-explanatory. However, it does require a mapping from carrier to domain for sending messages. For example, to send to a Verizon number, the text is sent to number@vztext.com. By default, Verizon is included, but any other carriers will need added to this mapping manually in the config file.

### Configure Twilio Gateway

The Twilio gateway only requires the Account SID, Auth Token, and a phone number. These can all easily be obtained from the Twilio console.