import psycopg2
from psycopg2 import sql
import random
from datetime import datetime, timedelta

# Dados de conexão ao banco de dados
DATABASE = "my_database"
USER = "my_user"
PASSWORD = "my@Pass%1234"
HOST = "localhost"  # ou o nome do serviço Docker, se estiver usando o Docker Compose
PORT = "5432"

# Funções auxiliares para gerar dados aleatórios
def generate_name(prefix, index):
    return f"{prefix}_{index}"

def generate_value():
    return round(random.uniform(10, 100), 2)

def generate_datetime():
    return datetime.now() - timedelta(days=random.randint(0, 365))

# Conectar ao banco de dados
conn = psycopg2.connect(database=DATABASE, user=USER, password=PASSWORD, host=HOST, port=PORT)
cursor = conn.cursor()

# Inserir dados nas tabelas
try:
    # Inserir classes
    for i in range(10):
        cursor.execute(
            "INSERT INTO CLASSES (NAME) VALUES (%s) RETURNING ID;",
            (generate_name('Class', i + 1),)
        )
        class_id = cursor.fetchone()[0]
        
        # Inserir objetos para cada classe
        object_ids = []
        for j in range(random.randint(1, 10)):
            cursor.execute(
                "INSERT INTO OBJECTS (NAME, VALUE, DATETIME, CLASS_ID) VALUES (%s, %s, %s, %s) RETURNING ID;",
                (generate_name('Object', j + 1), generate_value(), generate_datetime(), class_id)
            )
            object_ids.append(cursor.fetchone()[0])
        
        # Inserir itens para cada objeto
        for object_id in object_ids:
            item_ids = []
            for k in range(random.randint(2, 12)):
                cursor.execute(
                    "INSERT INTO ITEMS (NAME, VALUE, DATETIME) VALUES (%s, %s, %s) RETURNING ID;",
                    (generate_name('Item', k + 1), generate_value(), generate_datetime())
                )
                item_ids.append(cursor.fetchone()[0])

            # Inserir links entre objetos e itens
            for item_id in item_ids:
                cursor.execute(
                    "INSERT INTO OBJECT_ITEM_LINK (OBJECT_ID, ITEM_ID) VALUES (%s, %s);",
                    (object_id, item_id)
                )

    # Confirmar as transações
    conn.commit()
except Exception as e:
    print(f"An error occurred: {e}")
