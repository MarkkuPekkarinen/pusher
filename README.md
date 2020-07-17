[![Build Status](https://travis-ci.org/topfreegames/pusher.svg?branch=master)](https://travis-ci.org/topfreegames/pusher)
[![Coverage Status](https://coveralls.io/repos/github/topfreegames/pusher/badge.svg?branch=master)](https://coveralls.io/github/topfreegames/pusher?branch=master)
[![Docs](https://readthedocs.org/projects/pusher/badge/?version=latest)](http://pusher.readthedocs.io/en/latest/)

Pusher
======

### Dependencies
* Go 1.9
* Kafka >= 0.9.0
* [librdkafka](https://github.com/edenhill/librdkafka)

### Setup
First, set your $GOPATH ([Go Lang](https://golang.org/doc/install)) env variable and add $GOPATH/bin to your $PATH

```bash
make setup
```

### Building

```bash
make build
```

### Sending pushes

#### APNS

```bash
./bin/pusher apns -d -p
```

#### GCM

```bash
./bin/pusher gcm -d -p
```

### Automated tests

We're using [Ginkgo](https://onsi.github.io/ginkgo) and [Gomega](https://onsi.github.io/gomega) for testing our code. Since we're making extensive use of interfaces, external dependencies are mocked for all unit tests.

### Linter

```bash
make lint
```

#### Unit Tests
We'll try to keep testing coverage as high as possible. To run unit tests simply use:

```bash
make unit
```

To check the test coverage use:

```bash
make cover  # opens a html file
```

or

```bash
make test-coverage-func  # prints code coverage to the console
```

#### Integration Tests

We also have a few integration tests that actually connect to pusher external dependencies. For these you'll need to have docker installed.

To run integration tests run:

```
make integration
```

If you are running integration tests locally with the most recent librdkafka version (such as installed by brew) some of them will fail due to incompatible librdkafka version. 
Tests should work for librdkafka v0.11.5

### Benchmark

#### Create fake push data


```bash
cd bench

python create_bench.py test.txt

```

text.txt will contain several lines like this:

- GCM

```json
{
  "to": "XNAJY2WCN7RDH6B5APHXTCM793X28IO7780AB51F0F8OV3ENWXOIR40JRF3K9416AD9K029NEE3XTA229NJC0Y6DHCBO13EE6IFO6VRF8FICJ317AC5I3N1FCSJ7KIVXMKZ088BJOVS3PPJUG9CWV1J2",
  "notification": {
    "title": "Come play!",
    "body": "Helena miss you! come play!"
  },
  "dry_run": true
}
```

Note: If you want to actually send the pushes you need to set `dry_run: false` (default is true).

- APNS

```json
{
  "DeviceToken":"H9CSRZHTAUPOZP1ZDLK46DN8L1DS4JFIUKHXE33K77QHQLHZ650TG66U49ZQGFZV",
  "Payload": {
    "aps": {
      "alert": "Helena miss you! come play!"
    }
  },
  "push_expiry":0
}
```

#### Send pushes using the fake data:

To send the push using the fake data you need to start `pusher` in the correct mode (apns or gcm) and then produce to the Kafka topic and broker it will be listening to:

```bash
cat test.txt | kafka-console-producer --topic push-game_gcm --broker-list localhost:9941
```

### Available Environment variables

Pusher reads from Kafka the push notifications that should be sent. The container takes environment variables to specify this connection:

* `PUSHER_QUEUE_TOPICS` - List of Kafka topics, ex: `^push-[^-_]+_(apns|gcm)`
* `PUSHER_QUEUE_BROKERS` - List of Kafka brokers;
* `PUSHER_QUEUE_GROUP` - Kafka consumer group;
* `PUSHER_QUEUE_SESSIONTIMEOUT` - Kafka session timeout;
* `PUSHER_QUEUE_OFFSETRESETSTRATEGY` - Kafka offset reset strategy;
* `PUSHER_QUEUE_HANDLEALLMESSAGESBEFOREEXITING` - Boolean indicating if shutdown should wait for all messages to be handled;

Pusher gets the GCM or APNS keys info from environment variables:

* `PUSHER_GCM_APPS` - Comma separated APNS app names, ex: appname1,appname2;
* `PUSHER_APNS_APPS` - Comma separated GCM app names, ex: appname1,appname2;
* `PUSHER_APNS_CERTS_APPNAME` - App APNS certificate path
* `PUSHER_GCM_CERTS_APPNAME_APIKEY` - GCM App Api Key
* `PUSHER_GCM_CERTS_APPNAME_SENDERID` - GCM App SenderID

For feedbacks you must specify a list of reporters:

* `PUSHER_FEEDBACK_REPORTERS` - List of feedbacks reporters;

For each specified reporter you can set its configuration. For Kafka, it is as follows:

* `PUSHER_FEEDBACK_KAFKA_TOPICS` - List of Kafka topics;
* `PUSHER_FEEDBACK_KAFKA_BROKERS` - List of Kafka brokers;

The same logic is used for stats:

* `PUSHER_STATS_REPORTERS` - List of feedbacks reporters;

For a Statsd stats reporter, it is as follows:

* `PUSHER_STATS_STATSD_HOST` - Statsd host;
* `PUSHER_STATS_STATSD_PREFIX` - Prefix used in Statsd reported metrics;
* `PUSHER_STATS_STATSD_FLUSHINTERVALINMS` - Interval (in milliseconds) during which stats are aggregated before they are sent to the statsd server;

You can also specify invalid token handlers:

* `PUSHER_INVALIDTOKEN_HANDLERS` - List of invalid token handlers;

If Pusher needs to connect to a PostgreSQL database in order to delete invalid tokens the following environment variables must be specified:

* `PUSHER_INVALIDTOKEN_PG_USER` - User of the PostgreSQL Server to connect to;
* `PUSHER_INVALIDTOKEN_PG_PASS` - Password of the PostgreSQL Server to connect to;
* `PUSHER_INVALIDTOKEN_PG_HOST` - PostgreSQL host to connect to;
* `PUSHER_INVALIDTOKEN_PG_DATABASE` - PostgreSQL database to connect to;
* `PUSHER_INVALIDTOKEN_PG_PORT` - PostgreSQL port to connect to;
* `PUSHER_INVALIDTOKEN_PG_POOLSIZE` - PostgreSQL connection pool size;
* `PUSHER_INVALIDTOKEN_PG_MAXRETRIES` - PostgreSQL connection max retries;
* `PUSHER_INVALIDTOKEN_PG_CONNECTIONTIMEOUT` - Timeout for trying to establish connection;

Other than that, there are a couple more configurations you can pass using environment variables:

* `PUSHER_GRACEFULLSHUTDOWNTIMEOUT` - Pusher is exited gracefully but you should specify a timeout for termination in case it takes too long;


The APNS library we're using supports several concurrent workers.
* `PUSHER_APNS_CONCURRENTWORKERS` - Amount of concurrent workers;

The GCM library we're using requires that we specify a ping interval and timeout for the XMPP connection.
* `PUSHER_GCM_PINGINTERVAL` - Ping interval in seconds;
* `PUSHER_GCM_PINGTIMEOUT` - Ping timeout in seconds;

GCM supports at most 100 pending messages (see [Flow Control section](https://developers.google.com/cloud-messaging/ccs#flow) in GCM documentation).

* `PUSHER_GCM_MAXPENDINGMESSAGES` - Max pending messages;

If you wish Sentry integration simply set the following environment variable:

* `PUSHER_SENTRY_URL` - Sentry Client Key (DSN);
