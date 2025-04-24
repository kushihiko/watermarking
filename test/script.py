import easyocr
import cv2
import matplotlib.pyplot as plt

# Загружаем изображение
image_path = '/Users/kushihiko/Projects/watermarking/test/3.png'  # путь к изображению
image = cv2.imread(image_path)

# EasyOCR reader (можно указать язык, например ['en', 'ru'])
reader = easyocr.Reader(['ru'], gpu=False)

# Распознавание текста
results = reader.readtext(image)

# Отрисовка bounding boxes
for (bbox, text, confidence) in results:
    # Bounding box: список из 4 точек (каждая — [x, y])
    top_left = tuple(map(int, bbox[0]))
    bottom_right = tuple(map(int, bbox[2]))

    # Рисуем прямоугольник
    cv2.rectangle(image, top_left, bottom_right, (0, 255, 0), 2)

    # Подписываем текст
    cv2.putText(image, text, (top_left[0], top_left[1] - 10),
                cv2.FONT_HERSHEY_SIMPLEX, 0.7, (0, 0, 255), 2)

# Отображение результата
plt.figure(figsize=(10, 10))
plt.imshow(cv2.cvtColor(image, cv2.COLOR_BGR2RGB))
plt.axis('off')
plt.title('Detected Text with EasyOCR')
plt.show()