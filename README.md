# dlptest

A self hosted test suite for data loss prevention (DLP) Systems

## HTTP
Includes 8 test files by default. When the button is clicked it downloads the file from the server into a variable and then re-uploads it to the server. 

## S3
Configure the following environmental variables to use the S3 Test

* S3REGION
* S3BUCKET
* S3KEYID
* S3SECRET


# TODOs
- cleanup the UI
- add CI
- add more logging
- add docs for minimal s3 setup
- add verification for s3 upload
- add auto-delete for s3 upload
