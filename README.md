go-qotd
=======

A simple server which implements the Quote of the Day protocol:
[RFC 865](http://tools.ietf.org/html/rfc865).

Please note that QOTD servers are known to be vulnerable to
["ping-pong" DOS attacks](http://www.securityspace.com/smysecure/catid.html?id=10198)
so it's probably not wise to expose this service in production.

(In case you hadn't guessed, this is purely a silly just-for-fun project. :) )


Usage
-----

    $ ./qotd --help
    Usage of ./qotd:
      -debug=false: Print debug messages
      -file="/tmp/quotes": File to get quotes from
      -port=17: Port to run QOTD on

Data format
-----------

go-qotd accepts a quotes file in [fortune(6)](https://en.wikipedia.org/wiki/Fortune_%28Unix%29)
format, which consists of a set of quotations each separated by the character '%' on its
own line.

Sample data
-----------

Please note that the sample data included in this repo is the list of science
quotes from fortune-mod version 1.99.1-17.el7 from the Fedora Project.
fortune-mod BSD licensed.
