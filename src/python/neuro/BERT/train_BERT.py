import os
import torch
import pandas as pd
import numpy as np
import warnings

from datasets import Dataset, DatasetDict, load_dataset, concatenate_datasets
from sklearn.metrics import accuracy_score, precision_recall_fscore_support, confusion_matrix, classification_report

from transformers import AutoTokenizer, AutoModelForSequenceClassification, DataCollatorWithPadding
from transformers import TrainingArguments, Trainer
from transformers.integrations import ClearMLCallback

warnings.filterwarnings("ignore")

from dotenv import load_dotenv
from pathlib import Path

load_dotenv(Path('../.env'))

from clearml import Task, Logger

print("CUDA available" if torch.cuda.is_available() else "CUDA not available")
import kagglehub


dataset_path = kagglehub.dataset_download("mikhailklemin/kinopoisks-movies-reviews")
dataset_path = os.path.join(dataset_path, "dataset")

neg_dataset = load_dataset("text", data_dir=os.path.join(dataset_path, "neg"), split="train", num_proc=4)
pos_dataset = load_dataset("text", data_dir=os.path.join(dataset_path, "pos"), split="train", num_proc=4)
neu_dataset = load_dataset("text", data_dir=os.path.join(dataset_path, "neu"), split="train", num_proc=4)

neg_dataset = neg_dataset.map(lambda x: {"label": 0})
pos_dataset = pos_dataset.map(lambda x: {"label": 1})
neu_dataset = neu_dataset.map(lambda x: {"label": 2})

kinopoisk_data = concatenate_datasets([neg_dataset, pos_dataset, neu_dataset])
train_test_split = kinopoisk_data.train_test_split(test_size=0.2, shuffle=True, seed=42)
kinopoisk_data = DatasetDict({
    "train": train_test_split["train"],
    "test": train_test_split["test"]
})

# ради скорости обучаемся на части данных
small_ds = DatasetDict({
    "train": kinopoisk_data["train"].train_test_split(train_size=0.1, seed=42)["train"],
    "test": kinopoisk_data["test"].train_test_split(test_size=0.1, seed=42)["test"]
})

models_to_train = [
    # "DeepPavlov/rubert-base-cased",
    # "cointegrated/rubert-tiny",
    # "ai-forever/ruBert-base",
    "ai-forever/ruBert-large",
]

def compute_metrics(eval_pred):
    predictions, labels = eval_pred
    preds = np.argmax(predictions, axis=1)
    acc = accuracy_score(labels, preds)
    
    precision, recall, f1_macro, _ = precision_recall_fscore_support(labels, preds, average='macro')
    return {"accuracy": acc, "precision": precision, 
            "recall": recall, "macro_f1": f1_macro}
    
# пытаемся ускорить обучение
os.environ["CUDA_LAUNCH_BLOCKING"]="1"
os.environ["TOKENIZERS_PARALLELISM"] = "false"
os.environ["PYTORCH_CUDA_ALLOC_CONF"]="expandable_segments:True"

try:
    torch._inductor.config.max_autotune = False
    torch._inductor.config.max_autotune_gemm = False
except Exception:
    pass

for checkpoint in models_to_train:
    print(f"\n === Обучение модели {checkpoint} ===")
    
    task = Task.init(project_name="BERT_kinopoisk", task_name=f"{checkpoint}_training")
    
    tokenizer = AutoTokenizer.from_pretrained(checkpoint, cache_dir="models")
    model = AutoModelForSequenceClassification.from_pretrained(checkpoint, num_labels=3, cache_dir="models")
    tokenized_data = small_ds.map(
        lambda x: tokenizer(x["text"], truncation=True, padding="max_length", max_length=512), batched=True)
    data_collator = DataCollatorWithPadding(tokenizer=tokenizer)
    output_dir = "models" + checkpoint.split("/")[-1] + "_finetuned"
    n_steps = 5000                          # каждые n_steps шагов eval-делаем логгирование-сохранение

    training_args = TrainingArguments(
        num_train_epochs=3,                 # число эпох
        output_dir=output_dir,              # папка для сохранения модели и чекпоинтов
        per_device_train_batch_size=16,     # батч на одном устройстве для обучения
        per_device_eval_batch_size=32,      # батч для оценки
        learning_rate=2e-5,                 # начальный learning rate
        weight_decay=0.01,                  # коэффициент L2-регуляризации (weight decay)
        warmup_ratio=0.05,                  # прогреваемся, так якобы будет лучше!
        lr_scheduler_type="cosine",         # меняем lr
        fp16=torch.cuda.is_available(),     # используем fp16 при наличии GPU
        optim="adamw_torch",                # оптимизатор быстрее чем стандартный
        gradient_checkpointing=True,        # экономим VRAM
        gradient_accumulation_steps=2,      # копим на всякий    
        
        eval_strategy="steps",              # выполнять оценку на тестовом наборе после N шагов
        eval_steps=n_steps,                 # считаем метрики каждые n_steps шагов
        logging_strategy="steps",           # вкидываем логи каждые N шагов
        logging_steps=n_steps,              # каждые n_steps шагов
        save_strategy="steps",              # сохранять модель после каждой эпохи
        save_steps=n_steps,                 # сохраняем каждые n_steps шагов
        save_total_limit=2,                 # чтобы не засорялось
        
        load_best_model_at_end=True,        # загружаем лучшую модель после окончания обучения
        metric_for_best_model="accuracy",  # метрика, на основе которой выбирается лучшая модель
        greater_is_better=False,                 # у нас loss

        report_to=["clearml"]               # логировать в ClearML 
    )
    
    task.connect(training_args.to_dict())
    trainer = Trainer(
        model=model,
        args=training_args,
        train_dataset=tokenized_data["train"],
        eval_dataset=tokenized_data["test"],
        tokenizer=tokenizer,
        data_collator=data_collator,
        compute_metrics=compute_metrics,
        callbacks=[ClearMLCallback()],      
    )
    
    trainer.train()
    
    predictions = trainer.predict(tokenized_data["test"])
    preds = np.argmax(predictions.predictions, axis=1)
    y_true = tokenized_data["test"]["label"]

    acc = accuracy_score(y_true, preds)
    # print(f"Accuracy: {acc:.3f}")
    # print("Classification report:\n", classification_report(y_true, preds, digits=4))
    # print("Confusion Matrix:")
    # print(cm)
    cm = confusion_matrix(y_true, preds)

    # логируем финальные метрики и матрицу ошибок
    logger = task.get_logger()
    logger.report_scalar(title="test", series="accuracy", value=acc, iteration=0)
    labels = ["neg", "pos", "neu"]
    logger.report_confusion_matrix(
        title="Confusion Matrix (test)",
        series="ConfMatrix",
        matrix=cm.tolist(),
        xaxis=labels,
        yaxis=labels,
        iteration=0,
    )
    
    trainer.save_model(output_dir)

    task.close()