data "aws_iam_policy_document" "deploy" {
  statement {
    sid = "ProjectServices"
    actions = [
      "s3:*",
      "dynamodb:*",
      "cloudfront:*",
      "lambda:*",
      "apigateway:*",
      "execute-api:*",
      "logs:*",
      "cloudformation:*",
      "acm:*",
      "ssm:GetParameter*",
      "ssm:PutParameter",
      "ssm:DeleteParameter",
      "sqs:*",
      "events:*",
      "cloudwatch:*",
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

  # Чтобы этот же конфиг можно было применять локально, а не только под SSO.
  # Пользователей и ключи заводить нельзя — только политики и их привязки.
  statement {
    sid = "ManagePolicies"
    actions = [
      "iam:GetUser",
      "iam:ListAttachedUserPolicies",
      "iam:ListPolicies",
      "iam:ListEntitiesForPolicy",
      "iam:GetPolicy",
      "iam:GetPolicyVersion",
      "iam:ListPolicyVersions",
      "iam:CreatePolicy",
      "iam:CreatePolicyVersion",
      "iam:DeletePolicyVersion",
      "iam:DeletePolicy",
      "iam:AttachUserPolicy",
      "iam:DetachUserPolicy",
    ]
    resources = ["*"]
  }

  statement {
    sid = "ReadOnlyBilling"
    actions = [
      "ce:Get*",
      "ce:Describe*",
      "ce:List*",
      "budgets:View*",
      "budgets:Describe*",
      "freetier:Get*",
      "aws-portal:View*",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "deploy" {
  name        = "vsluh-deploy"
  description = "Локальный ключ: Terraform и разведка"
  policy      = data.aws_iam_policy_document.deploy.json
}

resource "aws_iam_user_policy_attachment" "deploy" {
  user       = "vsluh-deploy"
  policy_arn = aws_iam_policy.deploy.arn
}
