# Design considerations

- this app tries to conform to the [Twelve Factor-App](https://12factor.net/) manifesto
  - [currently with a known limitation according to principle VI](https://github.com/uyuni-project/hub-xmlrpc-api/issues/56)
- this app's code is structured according to the [Clean Code Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- this app fundamentally implements a [Gorilla](https://www.gorillatoolkit.org/) XMLRPC Server to serve its API
- it uses [kolo/xmlrpc](https://github.com/kolo/xmlrpc) as the XMLRPC client library to consume other XMLRPC APIs
- major packages contain a README.md file with a high-level explanation of its contents