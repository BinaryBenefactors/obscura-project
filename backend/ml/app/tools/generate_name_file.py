import hashlib
import uuid

def generate_name_file(data: str, extension:str) -> str:
    # Создаем хэш от входных данных
    hash_object = hashlib.md5(data.encode())
    hash_hex = hash_object.hexdigest()
    
    # Генерируем UUID
    unique_id = str(uuid.uuid4())

    # Формируем имя файла из хэша и UUID
    file_name = f"{hash_hex}_{unique_id}.{extension}"
    return file_name

if __name__ == '__main__':
    example_data = "example data"
    print(generate_name_file(example_data))
