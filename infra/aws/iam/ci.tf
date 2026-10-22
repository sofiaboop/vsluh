data "aws_iam_policy_document" "ci" {
  statement {
    sid       = "CloudFormationSam"
    actions   = ["cloudformation:*"]
    resources = ["*"]
  }

  statement {
    sid = "ServerlessStack"
    actions = [
      "lambda:*",
      "apigateway:*",
      "dynamodb:*",
      "logs:*",
      "events:*",
      "sqs:*",
      "cloudwatch:*",
      "ssm:GetParameter*",
    ]
    resources = ["*"]
  }

  statement {
    sid     = "ArtifactBuckets"
    actions = ["s3:*"]
    resources = [
      "arn:aws:s3:::aws-sam-cli-managed-*",
      "arn:aws:s3:::aws-sam-cli-managed-*/*",
      "arn:aws:s3:::vsluh-*",
      "arn:aws:s3:::vsluh-*/*",
    ]
  }

  # Бакеты статики названы по домену и под маску `vsluh-*` не попадают.
  statement {
    sid     = "SiteBuckets"
    actions = ["s3:ListBucket", "s3:GetObject", "s3:PutObject", "s3:DeleteObject"]
    resources = concat(
      [for b in var.site_buckets : "arn:aws:s3:::${b}"],
      [for b in var.site_buckets : "arn:aws:s3:::${b}/*"],
    )
  }

  # Дистрибуцию CI ищет по алиасу домена, поэтому нужен и список.
  statement {
    sid = "InvalidateCloudFront"
    actions = [
      "cloudfront:ListDistributions",
      "cloudfront:CreateInvalidation",
      "cloudfront:GetInvalidation",
    ]
    resources = ["*"]
  }

  statement {
    sid = "ListBuckets"
    actions = [
      "s3:ListAllMyBuckets",
      "s3:GetBucketLocation",
      "s3:CreateBucket",
    ]
    resources = ["*"]
  }

  statement {
    sid = "IamForServerless"
    actions = [
      "iam:CreateRole",
      "iam:DeleteRole",
      "iam:GetRole",
      "iam:PassRole",
      "iam:TagRole",
      "iam:UntagRole",
      "iam:AttachRolePolicy",
      "iam:DetachRolePolicy",
      "iam:PutRolePolicy",
      "iam:DeleteRolePolicy",
      "iam:GetRolePolicy",
      "iam:ListRolePolicies",
      "iam:ListAttachedRolePolicies",
      "iam:CreateServiceLinkedRole",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "ci" {
  name        = "vsluh-ci"
  description = "GitHub Actions: SAM-стек бэкенда и выкат статики"
  policy      = data.aws_iam_policy_document.ci.json
}

resource "aws_iam_user_policy_attachment" "ci" {
  user       = "vsluh-ci"
  policy_arn = aws_iam_policy.ci.arn
}
