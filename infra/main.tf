provider "aws" {
  region = var.region
}

resource "aws_eks_cluster" "snisid" {
  name     = "snisid-platform"
  role_arn = aws_iam_role.eks.arn

  vpc_config {
    subnet_ids = aws_subnet.private[*].id
  }
}

resource "aws_db_instance" "postgres" {
  allocated_storage    = 100
  engine               = "postgres"
  engine_version       = "16"
  instance_class       = "db.t4g.large"
  db_name              = "snisid"
  username             = "snisid"
  password             = var.db_password
  skip_final_snapshot  = true
}

resource "aws_msk_cluster" "kafka" {
  cluster_name           = "snisid-events"
  kafka_version          = "3.6.0"
  number_of_broker_nodes = 3

  broker_node_group_info {
    instance_type = "kafka.m5.large"
    client_subnets = aws_subnet.private[*].id
    security_groups = [aws_security_group.kafka.id]
  }
}
