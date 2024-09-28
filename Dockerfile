# Usar una imagen base de Go
FROM golang:1.22

# Setear el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el código fuente al contenedor
COPY . .

# Compilar la aplicación
RUN go mod download
RUN go build -o main .

# Ejecutar la aplicación
CMD ["./main"]
