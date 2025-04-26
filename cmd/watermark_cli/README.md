# watermark_cli

Консольная утилита для встраивания и извлечения водяной метки из PDF-документов.

## Использование

```sh
# Встраивание водяной метки
watermark_cli embed --config=config.yaml

# Извлечение водяной метки
watermark_cli extract --config=config.yaml
```

## Пример config.yaml

```yaml
pdf_path: "test.pdf"
image_folder: "tmp/"
output_folder: "output/"
font_path: "/System/Library/Fonts/Supplemental/Times New Roman.ttf"
language: "rus+eng"
whitelist: "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя"
blacklist: "-.,:;"
watermark: "01110101101100"
shift: 4
marker_length: 5
``` 