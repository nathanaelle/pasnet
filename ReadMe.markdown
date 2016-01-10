# PasNet

only the v1.0 is usable and used by other project.

the current master may be usable without TLS but not tester.

with TLS client code, the code is actually incomplete and neither verified nor tested :

 - missing the verification of the certificate chain
 - missing the looking for intermediate certificate if not provided
 - missing the OSCP verification if no staple is provided
 - ...

the TLS server code is totally missing. 

## License
2-Clause BSD
