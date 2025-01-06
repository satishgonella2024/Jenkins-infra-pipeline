resource "aws_s3_bucket" "example" {
  bucket = "jenkins-infra-pipeline-example-london"
}

resource "aws_s3_bucket_acl" "example_acl" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}