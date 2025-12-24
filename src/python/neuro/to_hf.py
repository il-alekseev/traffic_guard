from huggingface_hub import HfApi

REPO_ID = "pelevinNow/traffic_analyzer"  
LOCAL_DIR = "models"  
REPO_TYPE = "model"

api = HfApi()

# (опционально) создать репо программно:
api.create_repo(repo_id=REPO_ID, repo_type=REPO_TYPE, private=True, exist_ok=True)

api.upload_folder(
    repo_id=REPO_ID,
    repo_type=REPO_TYPE,
    folder_path=LOCAL_DIR,
    path_in_repo=".",          # класть в корень репо
    commit_message="Upload model artifacts",
)
print("Done")
