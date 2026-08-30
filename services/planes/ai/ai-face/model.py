import torch
import torch.nn as nn
import torch.nn.functional as F
from torchvision.models import resnet50

class ArcFaceModel(nn.Module):
    def __init__(self, embedding_size=512, num_classes=1000):
        super(ArcFaceModel, self).__init__()
        self.backbone = resnet50(pretrained=True)
        self.backbone.fc = nn.Linear(self.backbone.fc.in_features, embedding_size)
        self.bn = nn.BatchNorm1d(embedding_size)
        
        # ArcFace Layer (during training only)
        self.num_classes = num_classes
        self.embedding_size = embedding_size

    def forward(self, x):
        embeddings = self.backbone(x)
        embeddings = self.bn(embeddings)
        return F.normalize(embeddings)

class ArcMarginProduct(nn.Module):
    def __init__(self, in_features, out_features, s=30.0, m=0.50):
        super(ArcMarginProduct, self).__init__()
        self.in_features = in_features
        self.out_features = out_features
        self.s = s
        self.m = m
        self.weight = nn.Parameter(torch.FloatTensor(out_features, in_features))
        nn.init.xavier_uniform_(self.weight)

    def forward(self, input, label):
        cosine = F.linear(F.normalize(input), F.normalize(self.weight))
        sine = torch.sqrt(1.0 - torch.pow(cosine, 2))
        phi = cosine * torch.cos(torch.tensor(self.m)) - sine * torch.sin(torch.tensor(self.m))
        
        one_hot = torch.zeros(cosine.size(), device=input.device)
        one_hot.scatter_(1, label.view(-1, 1).long(), 1)
        output = (one_hot * phi) + ((1.0 - one_hot) * cosine)
        output *= self.s
        return output
