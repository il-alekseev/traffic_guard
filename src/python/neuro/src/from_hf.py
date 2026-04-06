import os
from dotenv import load_dotenv
from huggingface_hub import snapshot_download

load_dotenv()
REPO_ID = "pelevinNow/traffic_analyzer"  
OUT_DIR = "models"  
REPO_TYPE = "model"

path = snapshot_download(
    repo_id=REPO_ID,
    repo_type=REPO_TYPE,
    local_dir=OUT_DIR,
    local_dir_use_symlinks=False,
    token=os.getenv("HF_TOKEN")
)
