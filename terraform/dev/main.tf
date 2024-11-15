provider "aws" {
  region = "eu-central-1"
}

resource "aws_eks_node_group" "node_group" {
  cluster_name    = aws_eks_cluster.k8s_cluster.name
  node_group_name = "eks-node-group"
  node_role_arn   = aws_iam_role.eks_node.arn
  subnet_ids      = [
    "subnet-00d19fc19329b6526",
    "subnet-0bd8d7416085a6086",
    "subnet-05a0f3715378f8e0f"
  ]

  scaling_config {
    desired_size = 2
    max_size     = 3
    min_size     = 1
  }

  instance_types = ["t3.medium"]
  ami_type       = "AL2_x86_64" # Amazon Linux 2

  tags = {
    Name = "eks-node-group"
  }
}

resource "aws_eks_cluster" "k8s_cluster" {
  name     = "k8s-${var.env}"
  role_arn = aws_iam_role.eks.arn

  vpc_config {
    subnet_ids = [
      "subnet-00d19fc19329b6526",
      "subnet-0bd8d7416085a6086",
      "subnet-05a0f3715378f8e0f"
    ]
    endpoint_public_access = true
  }
}

resource "aws_db_instance" "db_instance" {
  identifier        = "db-${var.env}"
  allocated_storage = 20
  engine            = "mysql"
  instance_class    = "db.t3.micro"
  db_name           = "go_crud_api"
  username          = "root"
  password          = "SecurePass123!"
  skip_final_snapshot = true
}

resource "aws_iam_role" "eks_node" {
  name = "eks-node-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role" "eks" {
  name = "eks-cluster-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "eks.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_route_table" "public_rt" {
  vpc_id = "vpc-001c504f1ace09713"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "igw-0280fd2505a27a82e"
  }
}

resource "aws_route_table_association" "public_rt_association_a" {
  subnet_id      = "subnet-00d19fc19329b6526"
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_route_table_association" "public_rt_association_b" {
  subnet_id      = "subnet-0bd8d7416085a6086"
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_route_table_association" "public_rt_association_c" {
  subnet_id      = "subnet-05a0f3715378f8e0f"
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_iam_role_policy_attachment" "node_group_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.eks_node.name
}

resource "aws_iam_role_policy_attachment" "node_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.eks_node.name
}

resource "aws_iam_role_policy_attachment" "cni_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.eks_node.name
}

resource "aws_iam_role_policy_attachment" "eks_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.eks.name
}
resource "aws_iam_role_policy_attachment" "ecr_read_only" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.eks_node.name
}

variable "env" {}