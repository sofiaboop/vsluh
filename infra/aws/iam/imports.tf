# ─────────────────────────────────────────────────────────────────────────────
# Политики и привязки заводились руками до Terraform — забираем как есть.
# План обязан показать только дополнения из ci.tf и deploy.tf, не пересоздание.
# ─────────────────────────────────────────────────────────────────────────────

import {
  to = aws_iam_policy.ci
  id = "arn:aws:iam::661609779996:policy/vsluh-ci"
}

import {
  to = aws_iam_user_policy_attachment.ci
  id = "vsluh-ci/arn:aws:iam::661609779996:policy/vsluh-ci"
}

import {
  to = aws_iam_policy.deploy
  id = "arn:aws:iam::661609779996:policy/vsluh-deploy"
}

import {
  to = aws_iam_user_policy_attachment.deploy
  id = "vsluh-deploy/arn:aws:iam::661609779996:policy/vsluh-deploy"
}
