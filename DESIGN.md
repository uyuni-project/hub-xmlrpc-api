# Design considerations

- this app tries to conform to the [Twelve Factor-App](https://12factor.net/) manifesto
  - [currently with a known limitation according to principle VI](https://github.com/uyuni-project/hub-xmlrpc-api/issues/56)
- this app's code is structured according to the [Clean Code Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
  - the application is structured in concentric horizontal layers, according to the level of the policy the layer implements. At the same time, the application is structured in vertical layers, each layer implementing a specific usecase
  - the "gateway" package is placed at the center, and contains the bussiness rules and the highest policies of the application, and implements all the usecases. All the other packages contain lower level policies, and are plugins to the "gateway" package
  - the goal of this design is to isolate from each other things that change for different reasons and at different rates, in order to favor modularity, testatiblity, mantainablity, and specifically to avoid potential secundary effects when adding/changing/fixing a feature
- this app fundamentally implements a [Gorilla](https://www.gorillatoolkit.org/)-based XMLRPC Server to serve its API
  - in Gorilla's jargon, this implements a Codec, the Server being implemented by Gorilla itself
- it uses [kolo/xmlrpc](https://github.com/kolo/xmlrpc) as the XMLRPC client library to consume other XMLRPC APIs
- major packages contain a README.md file with a high-level explanation of its contents