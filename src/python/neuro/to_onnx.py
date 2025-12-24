import torch
from transformers import AutoTokenizer, AutoModelForSequenceClassification
from pathlib import Path

DEVICE = "cuda" if torch.cuda.is_available() else "cpu"

def export_model(
    model_dir: str,
    onnx_path: str,
    max_length: int = 256,
):
    model_dir = Path(model_dir)
    onnx_path = Path(onnx_path)

    tokenizer = AutoTokenizer.from_pretrained(model_dir, use_fast=True)
    model = AutoModelForSequenceClassification.from_pretrained(model_dir)
    model.eval().to(DEVICE)

    # dummy input
    dummy = tokenizer(
        "dummy text for onnx export",
        return_tensors="pt",
        padding="max_length",
        truncation=True,
        max_length=max_length,
    )

    dummy = {k: v.to(DEVICE) for k, v in dummy.items()}

    torch.onnx.export(
        model,
        args=tuple(dummy.values()),
        f=str(onnx_path),
        input_names=list(dummy.keys()),
        output_names=["logits"],
        dynamic_axes={
            "input_ids": {0: "batch", 1: "sequence"},
            "attention_mask": {0: "batch", 1: "sequence"},
            # если есть token_type_ids — добавится автоматически
            "logits": {0: "batch"},
        },
        opset_version=17,
        do_constant_folding=True,
    )

    print(f"✅ Exported: {onnx_path}")
    print(f"   num_labels = {model.config.num_labels}")

if __name__ == "__main__":
    export_model(
        "models/content/",
        "content_classifier.onnx",
        max_length=256,
    )


    export_model(
        "models/sentiment/",
        "sentiment_classifier.onnx",
        max_length=256,
    )
