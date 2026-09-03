# Evidencias de los TP's 
## Evidencias — TP1
### 1. Push directo a main rechazado
<img width="728" height="412" alt="Screenshot 2026-08-18 123551" src="https://github.com/user-attachments/assets/ba0b51d3-400d-445a-b19f-147cb616b70f" />
GitHub rechaza el push porque main está protegida y la regla alcanza también al dueño del repo.

### 2. El PR de la rama B no se puede mergear: conflicto (Aviso)
<img width="1874" height="508" alt="Screenshot 2026-08-18 124103" src="https://github.com/user-attachments/assets/2e17d3e1-4775-4289-9647-79ce85fe1b3a" />
Ocurre un conflicto por que desde dos ramas se toco el mismo archivo y mismas lineas, desde el mismo origen main.

### 3. Resolucion del conflicto
<img width="1874" height="508" alt="Screenshot 2026-08-18 124103" src="https://github.com/user-attachments/assets/4abc5439-b9cc-4989-99b0-17308d4de373" />
Seleccionamos con que cambios nos vamos a quedar

### 4. Creacion del Tag y Release v1.0.0
<img width="1415" height="649" alt="Screenshot 2026-08-18 194711" src="https://github.com/user-attachments/assets/835b730c-0c1f-43ae-8a02-5163cb9021c5" />
Fijamos una version estable del codigo con el tag y le damos una explicacion de los cambios y agregados con el release


## Evidencias — TP2
### 1. Captura de la terminal de mi repo levantando la app con "docker compose up --build", se levanta la imagen del back, front y todo.
<img width="1916" height="1017" alt="image" src="https://github.com/user-attachments/assets/06f09499-0c37-4260-88ff-a8dad1864625" />

### 2. Captura de mi navegador en localhost:3000 con la app levantada y corriendo
<img width="1916" height="967" alt="image" src="https://github.com/user-attachments/assets/3e2405dd-b82c-48eb-957b-2369c2818738" />
<img width="1917" height="970" alt="image" src="https://github.com/user-attachments/assets/41bf2932-6ed7-42aa-9149-c4fb8e0a0cd4" />

### 3. Prueba de persistencia
Hacemos docker compose down
<img width="1912" height="962" alt="image" src="https://github.com/user-attachments/assets/a1db5a19-2538-4879-9e88-09ee8f52cb13" />
Luego hacemos docker compose up -d y podemos ver que los planes persisten
<img width="1912" height="967" alt="image" src="https://github.com/user-attachments/assets/bd9e61ad-1bd8-44a2-a4c6-277f5ab99b58" />

### 4. Comparación de tamaños de imágenes (Multi-stage)
Se puede notar las diferencias de peso entre las imagenes base SDK pesadas y mis imágenes finales de producción (backend:ci y frontend:ci)
<img width="1416" height="777" alt="image" src="https://github.com/user-attachments/assets/c4f7e46e-e763-4113-b432-0ece6fee1888" />

### 5. Imagenes publicadas en el registry por my usuario de Github (frp274)
<img width="1162" height="388" alt="image" src="https://github.com/user-attachments/assets/9a4248a6-b2f3-49d8-aed9-5b219231eb9f" />


## Evidencias - TP4
### 1. Romper el build
Al abrir el PR, el workflow CI se disparará automáticamente.
El job build-backend fallará en el paso de compilación (RUN go build -o main ./cmd/api/main.go).
El status check del PR quedará en rojo (❌ Failed), bloqueando el botón de Merge si tienes configuradas las Branch Protection Rules / Required status checks.
<img width="1112" height="387" alt="Screenshot 2026-09-02 221454" src="https://github.com/user-attachments/assets/1b1dea8c-aaca-4564-98ad-e0ce5338b5ea" />
Pero al revertir el error, github automaticamente, no detecta error y permite el merge
<img width="1115" height="407" alt="image" src="https://github.com/user-attachments/assets/07fbe4d7-e191-48ac-886b-8ea04af640e5" />





