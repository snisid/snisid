import torch
import torch.distributed as dist
import torch.multiprocessing as mp
from torch.nn.parallel import DistributedDataParallel as DDP
from model import ArcFaceModel
import os

def setup(rank, world_size):
    os.environ['MASTER_ADDR'] = 'localhost'
    os.environ['MASTER_PORT'] = '12355'
    dist.init_process_group("nccl", rank=rank, world_size=world_size)

def cleanup():
    dist.destroy_process_group()

def train(rank, world_size):
    setup(rank, world_size)
    
    device = torch.device(f"cuda:{rank}")
    model = ArcFaceModel().to(device)
    ddp_model = DDP(model, device_ids=[rank])

    optimizer = torch.optim.Adam(ddp_model.parameters(), lr=0.001)
    
    # Training loop would go here
    print(f"Rank {rank} starting training...")
    
    cleanup()

if __name__ == "__main__":
    world_size = torch.cuda.device_count()
    if world_size > 1:
        mp.spawn(train, args=(world_size,), nprocs=world_size, join=True)
    else:
        print("Less than 2 GPUs found, skipping distributed training demo.")
